import { ChangeDetectionStrategy, Component } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideGithub, lucideGamepad2, lucideChevronDown } from '@ng-icons/lucide';
import { HlmButton } from '@spartan-ng/helm/button';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [HlmButton, NgIcon],
  providers: [
    provideIcons({
      lucideGithub,
      lucideGamepad2,
      lucideChevronDown,
    }),
  ],
  templateUrl: './header.html',
  styleUrl: './header.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class Header {}
