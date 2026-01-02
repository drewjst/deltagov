import { Component, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideSearch, lucideSparkles, lucideUser, lucideChevronDown, lucideFileText, lucideGitCompare, lucideHistory, lucideTrendingUp } from '@ng-icons/lucide';
import { BrnCommandImports } from '@spartan-ng/brain/command';
import { HlmCommandImports } from '@spartan-ng/helm/command';
import { HlmIcon } from '@spartan-ng/helm/icon';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [NgIcon, BrnCommandImports, HlmCommandImports, HlmIcon],
  providers: [provideIcons({ lucideSearch, lucideSparkles, lucideUser, lucideChevronDown, lucideFileText, lucideGitCompare, lucideHistory, lucideTrendingUp })],
  templateUrl: './header.html',
  styleUrl: './header.scss'
})
export class Header {
  protected readonly searchQuery = signal('');
}
