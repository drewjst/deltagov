import { computed, inject } from '@angular/core';
import { signalStore, withState, withComputed, withMethods, patchState } from '@ngrx/signals';
import { rxMethod } from '@ngrx/signals/rxjs-interop';
import { pipe, switchMap, tap, catchError, of, EMPTY } from 'rxjs';
import { BillService, Bill as ApiBill, BillVersion as ApiVersion, DiffResponse } from '../../services/bill.service';

// Types
export interface BillVersion {
  id: number;
  label: string;
  date: string;
  contentHash: string;
  versionCode: string;
}

export interface Bill {
  id: number;
  title: string;
  sponsor: string;
  status: string;
  congress?: number;
  billNumber?: number;
  billType?: string;
  originChamber?: string;
  updateDate?: string;
}

export interface DiffSegment {
  type: 'insertion' | 'deletion' | 'unchanged';
  text: string;
}

export interface DiffLine {
  lineNumber: number;
  type: 'insertion' | 'deletion' | 'unchanged';
  text: string;
}

export interface Delta {
  fromVersion: string;
  toVersion: string;
  segments: DiffSegment[];
  lines: DiffLine[];
  insertions: number;
  deletions: number;
}

export interface LivingBillState {
  bill: Bill | null;
  versions: BillVersion[];
  selectedFromVersion: number;
  selectedToVersion: number;
  delta: Delta | null;
  isLoadingBill: boolean;
  isLoadingVersions: boolean;
  isLoadingDiff: boolean;
  error: string | null;
}

const initialState: LivingBillState = {
  bill: null,
  versions: [],
  selectedFromVersion: 0,
  selectedToVersion: 0,
  delta: null,
  isLoadingBill: false,
  isLoadingVersions: false,
  isLoadingDiff: false,
  error: null,
};

export const LivingBillStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withComputed((store) => ({
    fromVersionOptions: computed(() =>
      store.versions().filter((v) => v.id !== store.selectedToVersion()),
    ),
    toVersionOptions: computed(() =>
      store.versions().filter((v) => v.id !== store.selectedFromVersion()),
    ),
    selectedFromVersionLabel: computed(() => {
      const version = store.versions().find((v) => v.id === store.selectedFromVersion());
      return version?.label ?? '';
    }),
    selectedToVersionLabel: computed(() => {
      const version = store.versions().find((v) => v.id === store.selectedToVersion());
      return version?.label ?? '';
    }),
    hasVersions: computed(() => store.versions().length > 0),
    canCompareDiff: computed(
      () => store.selectedFromVersion() !== 0 && store.selectedToVersion() !== 0,
    ),
    isLoading: computed(
      () => store.isLoadingBill() || store.isLoadingVersions() || store.isLoadingDiff(),
    ),
    // Computed diff lines for virtual scrolling
    diffLines: computed((): DiffLine[] => {
      const delta = store.delta();
      if (delta?.lines && delta.lines.length > 0) {
        return delta.lines;
      }
      // Fallback to segments if lines not available
      if (delta?.segments && delta.segments.length > 0) {
        return delta.segments.map((seg, idx) => ({
          lineNumber: idx + 1,
          type: seg.type,
          text: seg.text,
        }));
      }
      // Empty state
      return [];
    }),
    diffStats: computed(() => {
      const delta = store.delta();
      if (delta) {
        return {
          additions: delta.insertions,
          deletions: delta.deletions,
          total: delta.lines?.length || delta.segments?.length || 0,
        };
      }
      return {
        additions: 0,
        deletions: 0,
        total: 0,
      };
    }),
  })),
  withMethods((store, billService = inject(BillService)) => ({
    // Version selection
    selectFromVersion(versionId: number): void {
      patchState(store, { selectedFromVersion: versionId });
    },
    selectToVersion(versionId: number): void {
      patchState(store, { selectedToVersion: versionId });
    },

    // Load H.R. 1 - The One Big Beautiful Bill
    loadHR1: rxMethod<void>(
      pipe(
        tap(() => patchState(store, { isLoadingBill: true, error: null })),
        switchMap(() =>
          billService.getHR1().pipe(
            tap((bill) => {
              const versions: BillVersion[] = (bill.versions || []).map((v) => ({
                id: v.id,
                label: v.label,
                date: v.date,
                contentHash: v.contentHash,
                versionCode: v.versionCode,
              }));

              const mappedBill: Bill = {
                id: bill.id,
                title: bill.title,
                sponsor: bill.sponsor || 'Unknown',
                status: bill.currentStatus || 'Unknown',
                congress: bill.congress,
                billNumber: bill.billNumber,
                billType: bill.billType,
                originChamber: bill.originChamber,
                updateDate: bill.updateDate,
              };

              patchState(store, {
                bill: mappedBill,
                versions,
                isLoadingBill: false,
              });

              // Auto-select first two versions if available
              if (versions.length >= 2) {
                patchState(store, {
                  selectedFromVersion: versions[0].id,
                  selectedToVersion: versions[versions.length - 1].id,
                });
              }
            }),
            catchError((error) => {
              patchState(store, {
                error: error.message || 'Failed to load H.R. 1',
                isLoadingBill: false,
              });
              return EMPTY;
            }),
          ),
        ),
      ),
    ),

    // Load diff between selected versions
    loadDiff: rxMethod<{ billId: number; fromVersionId: number; toVersionId: number }>(
      pipe(
        tap(() => patchState(store, { isLoadingDiff: true, error: null })),
        switchMap(({ billId, fromVersionId, toVersionId }) =>
          billService.computeDiff(billId, fromVersionId, toVersionId).pipe(
            tap((diff) => {
              const delta: Delta = {
                fromVersion: diff.fromVersion,
                toVersion: diff.toVersion,
                segments: diff.segments,
                lines: diff.lines,
                insertions: diff.insertions,
                deletions: diff.deletions,
              };
              patchState(store, { delta, isLoadingDiff: false });
            }),
            catchError((error) => {
              patchState(store, {
                error: error.message || 'Failed to compute diff',
                isLoadingDiff: false,
              });
              return EMPTY;
            }),
          ),
        ),
      ),
    ),

    // Loading state management (for use with service calls)
    setLoadingBill(isLoading: boolean): void {
      patchState(store, { isLoadingBill: isLoading });
    },
    setLoadingVersions(isLoading: boolean): void {
      patchState(store, { isLoadingVersions: isLoading });
    },
    setLoadingDiff(isLoading: boolean): void {
      patchState(store, { isLoadingDiff: isLoading });
    },

    // Data setters (will be called after service responses)
    setBill(bill: Bill): void {
      patchState(store, { bill, error: null });
    },
    setVersions(versions: BillVersion[]): void {
      patchState(store, { versions, error: null });
      // Auto-select latest two versions if available
      if (versions.length >= 2) {
        patchState(store, {
          selectedFromVersion: versions[versions.length - 2].id,
          selectedToVersion: versions[versions.length - 1].id,
        });
      }
    },
    setDelta(delta: Delta): void {
      patchState(store, { delta, error: null });
    },

    // Error handling
    setError(error: string): void {
      patchState(store, {
        error,
        isLoadingBill: false,
        isLoadingVersions: false,
        isLoadingDiff: false,
      });
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
