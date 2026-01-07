import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideChevronDown,
  lucideSettings,
  lucideLogOut,
  lucideUserCircle,
} from '@ng-icons/lucide';
import { HlmAvatarImports } from '@spartan-ng/helm/avatar';
import { HlmButton } from '@spartan-ng/helm/button';
import { HlmDropdownMenuImports } from '@spartan-ng/helm/dropdown-menu';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { HlmMenubarImports } from '@spartan-ng/helm/menubar';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    RouterLink,
    RouterLinkActive,
    NgIcon,
    HlmAvatarImports,
    HlmButton,
    HlmDropdownMenuImports,
    HlmIcon,
    HlmMenubarImports,
  ],
  providers: [
    provideIcons({
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
  protected readonly navLinks = [
    { path: '/lex', label: 'Lex' },
    { path: '/diffs', label: 'Diffs' },
    { path: '/history', label: 'History' },
    { path: '/influence', label: 'Influence' },
    { path: '/money-trail', label: 'Money Trail' },
  ];
}
