import { Component, ChangeDetectionStrategy, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { UpperCasePipe } from '@angular/common';
import { LexStore } from './lex.store';
import { HlmButton } from '@spartan-ng/helm/button';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideSearch,
  lucideLoader2,
  lucideChevronLeft,
  lucideChevronRight,
  lucideX,
  lucideFilter,
  lucideFileText,
} from '@ng-icons/lucide';

@Component({
  selector: 'app-lex-page',
  standalone: true,
  imports: [FormsModule, UpperCasePipe, HlmButton, NgIcon],
  providers: [
    provideIcons({
      lucideSearch,
      lucideLoader2,
      lucideChevronLeft,
      lucideChevronRight,
      lucideX,
      lucideFilter,
      lucideFileText,
    }),
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="h-full flex flex-col bg-background">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-border">
        <h1 class="text-2xl font-bold text-foreground">Lex</h1>
        <p class="text-sm text-muted-foreground mt-1">Search and explore legislative bills</p>
      </div>

      <!-- Search and Filters -->
      <div class="px-6 py-4 border-b border-border bg-card">
        <div class="flex flex-col gap-4">
          <!-- Search Row -->
          <div class="flex gap-3">
            <!-- Sponsor Search -->
            <div class="flex-1 relative">
              <ng-icon
                name="lucideSearch"
                class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
                size="16"
              />
              <input
                type="text"
                placeholder="Search by sponsor name..."
                [ngModel]="sponsorInput()"
                (ngModelChange)="sponsorInput.set($event)"
                (keyup.enter)="onSearch()"
                class="w-full h-9 pl-9 pr-3 rounded-md border border-input bg-background text-sm
                       placeholder:text-muted-foreground focus:outline-none focus:ring-2
                       focus:ring-ring focus:border-transparent transition-all"
              />
            </div>

            <!-- Query Search -->
            <div class="flex-1 relative">
              <ng-icon
                name="lucideFileText"
                class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
                size="16"
              />
              <input
                type="text"
                placeholder="Search bill titles..."
                [ngModel]="queryInput()"
                (ngModelChange)="queryInput.set($event)"
                (keyup.enter)="onSearch()"
                class="w-full h-9 pl-9 pr-3 rounded-md border border-input bg-background text-sm
                       placeholder:text-muted-foreground focus:outline-none focus:ring-2
                       focus:ring-ring focus:border-transparent transition-all"
              />
            </div>

            <!-- Search Button -->
            <button
              hlmBtn
              variant="default"
              (click)="onSearch()"
              [disabled]="store.isLoading()"
              class="px-4"
            >
              @if (store.isLoading()) {
                <ng-icon name="lucideLoader2" class="animate-spin" size="16" />
              } @else {
                <ng-icon name="lucideSearch" size="16" />
              }
              <span>Search</span>
            </button>
          </div>

          <!-- Filter Row -->
          <div class="flex items-center gap-3 flex-wrap">
            <!-- Congress Filter -->
            <select
              [ngModel]="congressInput()"
              (ngModelChange)="congressInput.set($event)"
              class="h-9 px-3 rounded-md border border-input bg-background text-sm
                     focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
            >
              <option value="">All Congresses</option>
              <option value="119">119th Congress</option>
              <option value="118">118th Congress</option>
              <option value="117">117th Congress</option>
              <option value="116">116th Congress</option>
            </select>

            <!-- Bill Type Filter -->
            <select
              [ngModel]="typeInput()"
              (ngModelChange)="typeInput.set($event)"
              class="h-9 px-3 rounded-md border border-input bg-background text-sm
                     focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
            >
              <option value="">All Types</option>
              <option value="hr">House Bills (H.R.)</option>
              <option value="s">Senate Bills (S.)</option>
              <option value="hjres">House Joint Resolutions</option>
              <option value="sjres">Senate Joint Resolutions</option>
              <option value="hconres">House Concurrent Resolutions</option>
              <option value="sconres">Senate Concurrent Resolutions</option>
              <option value="hres">House Simple Resolutions</option>
              <option value="sres">Senate Simple Resolutions</option>
            </select>

            <!-- Spending Only Toggle -->
            <label class="flex items-center gap-2 text-sm text-foreground cursor-pointer">
              <input
                type="checkbox"
                [ngModel]="spendingInput()"
                (ngModelChange)="spendingInput.set($event)"
                class="h-4 w-4 rounded border-input text-primary focus:ring-ring"
              />
              <span>Spending Bills Only</span>
            </label>

            <!-- Clear Filters -->
            @if (store.hasActiveFilters()) {
              <button
                hlmBtn
                variant="ghost"
                size="sm"
                (click)="onClearFilters()"
                class="text-muted-foreground"
              >
                <ng-icon name="lucideX" size="14" />
                <span>Clear Filters</span>
              </button>
            }
          </div>
        </div>
      </div>

      <!-- Results Info -->
      <div class="px-6 py-2 border-b border-border bg-muted/30 text-sm text-muted-foreground">
        @if (store.isLoading()) {
          <span>Searching...</span>
        } @else if (store.isEmpty()) {
          <span>No bills found. Try adjusting your search criteria.</span>
        } @else {
          <span>
            Showing {{ store.showingFrom() }}-{{ store.showingTo() }} of {{ store.total() }} bills
          </span>
        }
      </div>

      <!-- Error Display -->
      @if (store.error()) {
        <div class="px-6 py-3 bg-destructive/10 border-b border-destructive/20 text-destructive text-sm">
          {{ store.error() }}
        </div>
      }

      <!-- Bills Table -->
      <div class="flex-1 overflow-auto">
        @if (store.isLoading() && store.bills().length === 0) {
          <div class="flex items-center justify-center h-64">
            <ng-icon name="lucideLoader2" class="animate-spin text-primary" size="32" />
          </div>
        } @else {
          <table class="w-full">
            <thead class="sticky top-0 bg-muted/50 backdrop-blur-sm border-b border-border">
              <tr>
                <th class="text-left px-6 py-3 text-sm font-semibold text-foreground">Title</th>
                <th class="text-left px-6 py-3 text-sm font-semibold text-foreground w-48">Sponsor</th>
                <th class="text-left px-6 py-3 text-sm font-semibold text-foreground w-48">Status</th>
                <th class="text-left px-6 py-3 text-sm font-semibold text-foreground w-32">Congress</th>
              </tr>
            </thead>
            <tbody>
              @for (bill of store.bills(); track bill.id) {
                <tr
                  class="border-b border-border hover:bg-muted/50 cursor-pointer transition-colors"
                  (click)="onBillClick(bill.id)"
                >
                  <td class="px-6 py-4">
                    <div class="font-medium text-foreground">
                      {{ bill.billType | uppercase }}.{{ bill.billNumber }}
                    </div>
                    <div class="text-sm text-muted-foreground line-clamp-2">
                      {{ bill.title }}
                    </div>
                  </td>
                  <td class="px-6 py-4 text-sm text-muted-foreground">
                    {{ bill.sponsor || 'Unknown' }}
                  </td>
                  <td class="px-6 py-4">
                    <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium
                                 bg-primary/10 text-primary">
                      {{ bill.currentStatus || 'Pending' }}
                    </span>
                  </td>
                  <td class="px-6 py-4 text-sm text-muted-foreground">
                    {{ bill.congress }}th
                  </td>
                </tr>
              } @empty {
                @if (!store.isLoading()) {
                  <tr>
                    <td colspan="4" class="px-6 py-12 text-center text-muted-foreground">
                      <ng-icon name="lucideFileText" class="mx-auto mb-2 opacity-50" size="32" />
                      <p>No bills found</p>
                      <p class="text-sm mt-1">Try adjusting your search filters</p>
                    </td>
                  </tr>
                }
              }
            </tbody>
          </table>
        }
      </div>

      <!-- Pagination -->
      @if (store.totalPages() > 1) {
        <div class="px-6 py-3 border-t border-border bg-card flex items-center justify-between">
          <div class="text-sm text-muted-foreground">
            Page {{ store.currentPage() }} of {{ store.totalPages() }}
          </div>
          <div class="flex gap-2">
            <button
              hlmBtn
              variant="outline"
              size="sm"
              [disabled]="!store.hasPrevPage() || store.isLoading()"
              (click)="onPrevPage()"
            >
              <ng-icon name="lucideChevronLeft" size="16" />
              <span>Previous</span>
            </button>
            <button
              hlmBtn
              variant="outline"
              size="sm"
              [disabled]="!store.hasNextPage() || store.isLoading()"
              (click)="onNextPage()"
            >
              <span>Next</span>
              <ng-icon name="lucideChevronRight" size="16" />
            </button>
          </div>
        </div>
      }
    </div>
  `,
  styles: `
    .line-clamp-2 {
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
  `,
})
export class LexPage implements OnInit {
  readonly store = inject(LexStore);
  private readonly router = inject(Router);

  // Local input signals for form binding
  readonly sponsorInput = signal('');
  readonly queryInput = signal('');
  readonly congressInput = signal('');
  readonly typeInput = signal('');
  readonly spendingInput = signal(false);

  ngOnInit(): void {
    // Load initial bills
    this.store.loadAll();
  }

  onSearch(): void {
    // Update store with current filter values
    this.store.setSponsorFilter(this.sponsorInput());
    this.store.setSearchQuery(this.queryInput());
    this.store.setCongressFilter(this.congressInput() ? parseInt(this.congressInput(), 10) : null);
    this.store.setTypeFilter(this.typeInput());
    this.store.setSpendingOnly(this.spendingInput());

    // Trigger search
    this.store.loadAll();
  }

  onClearFilters(): void {
    this.sponsorInput.set('');
    this.queryInput.set('');
    this.congressInput.set('');
    this.typeInput.set('');
    this.spendingInput.set(false);
    this.store.clearFilters();
    this.store.loadAll();
  }

  onBillClick(billId: number): void {
    this.router.navigate(['/lex', billId]);
  }

  onPrevPage(): void {
    this.store.prevPage();
    this.store.loadAll();
  }

  onNextPage(): void {
    this.store.nextPage();
    this.store.loadAll();
  }
}
