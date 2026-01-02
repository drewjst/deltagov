import { ChangeDetectionStrategy, Component, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideSearch, lucideUser, lucideChevronDown, lucideSettings, lucideLogOut, lucideUserCircle } from '@ng-icons/lucide';
import { HlmDropdownMenuImports } from '@spartan-ng/helm/dropdown-menu';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { HlmMenubarImports } from '@spartan-ng/helm/menubar';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, NgIcon, HlmDropdownMenuImports, HlmIcon, HlmMenubarImports],
  providers: [provideIcons({ lucideSearch, lucideUser, lucideChevronDown, lucideSettings, lucideLogOut, lucideUserCircle })],
  templateUrl: './header.html',
  styleUrl: './header.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class Header {
  protected readonly searchQuery = signal('');

  onSearch() {
    console.log('Searching for:', this.searchQuery());
  }
}
