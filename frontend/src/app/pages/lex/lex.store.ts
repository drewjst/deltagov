import { computed, inject } from '@angular/core';
import { signalStore, withState, withComputed, withMethods, patchState } from '@ngrx/signals';
import { rxMethod } from '@ngrx/signals/rxjs-interop';
import { pipe, switchMap, tap, catchError, EMPTY, debounceTime } from 'rxjs';
import { BillService, Bill, LexSearchParams } from '../../services/bill.service';

// Types
export interface LexSearchState {
  bills: Bill[];
  total: number;
  limit: number;
  offset: number;
  isLoading: boolean;
  error: string | null;
  // Search filters
  searchQuery: string;
  sponsorFilter: string;
  congressFilter: number | null;
  typeFilter: string;
  spendingOnly: boolean;
}

const initialState: LexSearchState = {
  bills: [],
  total: 0,
  limit: 20,
  offset: 0,
  isLoading: false,
  error: null,
  searchQuery: '',
  sponsorFilter: '',
  congressFilter: null,
  typeFilter: '',
  spendingOnly: false,
};

export const LexStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withComputed((store) => ({
    // Check if there are any active filters
    hasActiveFilters: computed(() =>
      store.searchQuery().length > 0 ||
      store.sponsorFilter().length > 0 ||
      store.congressFilter() !== null ||
      store.typeFilter().length > 0 ||
      store.spendingOnly()
    ),
    // Pagination info
    currentPage: computed(() => Math.floor(store.offset() / store.limit()) + 1),
    totalPages: computed(() => Math.ceil(store.total() / store.limit())),
    hasNextPage: computed(() => store.offset() + store.limit() < store.total()),
    hasPrevPage: computed(() => store.offset() > 0),
    // Display info
    showingFrom: computed(() => store.total() > 0 ? store.offset() + 1 : 0),
    showingTo: computed(() => Math.min(store.offset() + store.limit(), store.total())),
    // Empty state
    isEmpty: computed(() => !store.isLoading() && store.bills().length === 0),
  })),
  withMethods((store, billService = inject(BillService)) => ({
    // Update search filters
    setSearchQuery(query: string): void {
      patchState(store, { searchQuery: query, offset: 0 });
    },
    setSponsorFilter(sponsor: string): void {
      patchState(store, { sponsorFilter: sponsor, offset: 0 });
    },
    setCongressFilter(congress: number | null): void {
      patchState(store, { congressFilter: congress, offset: 0 });
    },
    setTypeFilter(type: string): void {
      patchState(store, { typeFilter: type, offset: 0 });
    },
    setSpendingOnly(spending: boolean): void {
      patchState(store, { spendingOnly: spending, offset: 0 });
    },

    // Clear all filters
    clearFilters(): void {
      patchState(store, {
        searchQuery: '',
        sponsorFilter: '',
        congressFilter: null,
        typeFilter: '',
        spendingOnly: false,
        offset: 0,
      });
    },

    // Pagination
    nextPage(): void {
      if (store.offset() + store.limit() < store.total()) {
        patchState(store, { offset: store.offset() + store.limit() });
      }
    },
    prevPage(): void {
      if (store.offset() > 0) {
        patchState(store, { offset: Math.max(0, store.offset() - store.limit()) });
      }
    },
    goToPage(page: number): void {
      const newOffset = (page - 1) * store.limit();
      if (newOffset >= 0 && newOffset < store.total()) {
        patchState(store, { offset: newOffset });
      }
    },

    // Main search method using rxMethod
    loadBills: rxMethod<LexSearchParams | void>(
      pipe(
        debounceTime(300), // Debounce rapid searches
        tap(() => patchState(store, { isLoading: true, error: null })),
        switchMap((params) => {
          // Build search params from current state if not provided
          const searchParams: LexSearchParams = params || {
            query: store.searchQuery() || undefined,
            sponsor: store.sponsorFilter() || undefined,
            congress: store.congressFilter() ?? undefined,
            type: store.typeFilter() || undefined,
            spending: store.spendingOnly() || undefined,
            limit: store.limit(),
            offset: store.offset(),
          };

          return billService.searchBills(searchParams).pipe(
            tap((response) => {
              patchState(store, {
                bills: response.bills,
                total: response.total,
                limit: response.limit,
                offset: response.offset,
                isLoading: false,
              });
            }),
            catchError((error) => {
              patchState(store, {
                error: error.message || 'Failed to search bills',
                isLoading: false,
                bills: [],
              });
              return EMPTY;
            }),
          );
        }),
      ),
    ),

    // Load all bills (convenience method)
    loadAll: rxMethod<void>(
      pipe(
        tap(() => patchState(store, { isLoading: true, error: null })),
        switchMap(() => {
          const searchParams: LexSearchParams = {
            query: store.searchQuery() || undefined,
            sponsor: store.sponsorFilter() || undefined,
            congress: store.congressFilter() ?? undefined,
            type: store.typeFilter() || undefined,
            spending: store.spendingOnly() || undefined,
            limit: store.limit(),
            offset: store.offset(),
          };

          return billService.searchBills(searchParams).pipe(
            tap((response) => {
              patchState(store, {
                bills: response.bills,
                total: response.total,
                limit: response.limit,
                offset: response.offset,
                isLoading: false,
              });
            }),
            catchError((error) => {
              patchState(store, {
                error: error.message || 'Failed to load bills',
                isLoading: false,
                bills: [],
              });
              return EMPTY;
            }),
          );
        }),
      ),
    ),

    // Error handling
    setError(error: string): void {
      patchState(store, { error, isLoading: false });
    },
    clearError(): void {
      patchState(store, { error: null });
    },

    // Reset store
    reset(): void {
      patchState(store, initialState);
    },
  })),
);
