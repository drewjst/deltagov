import { computed } from '@angular/core';
import { signalStore, withState, withComputed, withMethods, patchState } from '@ngrx/signals';

// Types
export interface BillVersion {
  id: string;
  label: string;
  date: string;
  contentHash: string;
}

export interface Bill {
  id: string;
  title: string;
  sponsor: string;
  status: string;
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
}

export interface LivingBillState {
  bill: Bill | null;
  versions: BillVersion[];
  selectedFromVersion: string;
  selectedToVersion: string;
  delta: Delta | null;
  isLoadingBill: boolean;
  isLoadingVersions: boolean;
  isLoadingDiff: boolean;
  error: string | null;
}

const initialState: LivingBillState = {
  bill: null,
  versions: [],
  selectedFromVersion: '',
  selectedToVersion: '',
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
      () => store.selectedFromVersion() !== '' && store.selectedToVersion() !== '',
    ),
    isLoading: computed(
      () => store.isLoadingBill() || store.isLoadingVersions() || store.isLoadingDiff(),
    ),
    // Computed diff lines for virtual scrolling - will be derived from delta when available
    diffLines: computed((): DiffLine[] => {
      // Mock data until real diff engine is connected
      const mockLines: DiffLine[] = [
        { lineNumber: 1, type: 'unchanged', text: 'SECTION 1. SHORT TITLE.' },
        { lineNumber: 2, type: 'unchanged', text: 'This Act may be cited as the "Federal Budget Act of 2025".' },
        { lineNumber: 3, type: 'unchanged', text: '' },
        { lineNumber: 4, type: 'unchanged', text: 'SECTION 2. APPROPRIATIONS.' },
        { lineNumber: 5, type: 'deletion', text: '(a) There is appropriated $500,000,000 for infrastructure.' },
        { lineNumber: 5, type: 'insertion', text: '(a) There is appropriated $750,000,000 for infrastructure.' },
        { lineNumber: 6, type: 'unchanged', text: '' },
        { lineNumber: 7, type: 'deletion', text: '(b) Funds shall be distributed over a period of 3 years.' },
        { lineNumber: 7, type: 'insertion', text: '(b) Funds shall be distributed over a period of 5 years.' },
        { lineNumber: 8, type: 'unchanged', text: '' },
        { lineNumber: 9, type: 'insertion', text: '(c) Priority shall be given to rural communities.' },
        { lineNumber: 10, type: 'unchanged', text: '' },
        { lineNumber: 11, type: 'unchanged', text: 'SECTION 3. OVERSIGHT.' },
        { lineNumber: 12, type: 'unchanged', text: 'The Government Accountability Office shall conduct annual audits.' },
      ];
      // Return mock data or parse delta.segments when available
      return store.delta()?.segments
        ? store.delta()!.segments.map((seg, idx) => ({
            lineNumber: idx + 1,
            type: seg.type,
            text: seg.text,
          }))
        : mockLines;
    }),
    diffStats: computed(() => {
      const lines = store.delta()?.segments ?? [];
      return {
        additions: lines.filter(l => l.type === 'insertion').length || 3,
        deletions: lines.filter(l => l.type === 'deletion').length || 2,
        total: lines.length || 14,
      };
    }),
  })),
  withMethods((store) => ({
    // Version selection
    selectFromVersion(versionId: string): void {
      patchState(store, { selectedFromVersion: versionId });
    },
    selectToVersion(versionId: string): void {
      patchState(store, { selectedToVersion: versionId });
    },

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
