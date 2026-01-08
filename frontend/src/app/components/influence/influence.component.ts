import { Component, ChangeDetectionStrategy, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { toSignal } from '@angular/core/rxjs-interop';

import { HlmAvatar, HlmAvatarImage, HlmAvatarFallback } from '@spartan-ng/helm/avatar';
import { HlmCard, HlmCardHeader, HlmCardTitle, HlmCardContent } from '@spartan-ng/helm/card';
import { HlmButton } from '@spartan-ng/helm/button';

import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideSearch, lucideFilter, lucideInfo, lucideFileText } from '@ng-icons/lucide';

import { InfluenceService } from '../../services/influence.service';

@Component({
  selector: 'app-influence',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    HlmAvatar,
    HlmAvatarImage,
    HlmAvatarFallback,
    HlmCard,
    HlmCardHeader,
    HlmCardTitle,
    HlmCardContent,
    HlmButton,
    NgIcon
  ],
  providers: [
    provideIcons({
      lucideSearch,
      lucideFilter,
      lucideInfo,
      lucideFileText
    })
  ],
  templateUrl: './influence.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class InfluenceComponent {
  private influenceService = inject(InfluenceService);

  // In a real scenario, billId would come from route params or input
  data = toSignal(this.influenceService.getInfluenceData());

  searchTerm = signal('');

  filteredStakeholders = computed(() => {
    const currentData = this.data();
    if (!currentData) return [];

    const search = this.searchTerm().toLowerCase().trim();
    if (!search) return currentData.stakeholders;

    return currentData.stakeholders.filter(s =>
      s.organization.toLowerCase().includes(search)
    );
  });
}
