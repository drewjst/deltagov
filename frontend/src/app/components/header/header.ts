import { ChangeDetectionStrategy, Component, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideSearch,
  lucideChevronDown,
  lucideSettings,
  lucideLogOut,
  lucideUserCircle,
} from '@ng-icons/lucide';
import { HlmAvatarImports } from '@spartan-ng/helm/avatar';
import { HlmDropdownMenuImports } from '@spartan-ng/helm/dropdown-menu';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { HlmMenubarImports } from '@spartan-ng/helm/menubar';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [NgIcon, HlmAvatarImports, HlmDropdownMenuImports, HlmIcon, HlmMenubarImports],
  providers: [
    provideIcons({
      lucideSearch,
      lucideChevronDown,
      lucideSettings,
      lucideLogOut,
      lucideUserCircle,
    }),
  ],
  templateUrl: './header.html',
  styleUrl: './header.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class Header {
  protected readonly searchQuery = signal('');

  onSearch(): void {
    console.log('Searching for:', this.searchQuery());
  }
}
