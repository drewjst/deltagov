import { computed } from '@angular/core';
import {
  signalStore,
  withState,
  withComputed,
  withMethods,
  patchState,
} from '@ngrx/signals';

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
      store.versions().filter((v) => v.id !== store.selectedToVersion())
    ),
    toVersionOptions: computed(() =>
      store.versions().filter((v) => v.id !== store.selectedFromVersion())
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
      () => store.selectedFromVersion() !== '' && store.selectedToVersion() !== ''
    ),
    isLoading: computed(
      () => store.isLoadingBill() || store.isLoadingVersions() || store.isLoadingDiff()
    ),
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
  }))
);
