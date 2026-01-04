import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideArrowRight, lucideGitCommit, lucideGitBranch, lucideFileText } from '@ng-icons/lucide';
import { HlmButton } from '@spartan-ng/helm/button';
import { HlmIcon } from '@spartan-ng/helm/icon';

@Component({
  selector: 'app-landing-page',
  standalone: true,
  imports: [
    RouterLink,
    NgIcon,
    HlmButton,
    HlmIcon
  ],
  providers: [
    provideIcons({
      lucideArrowRight,
      lucideGitCommit,
      lucideGitBranch,
      lucideFileText
    }),
  ],
  templateUrl: './landing.html',
  styleUrl: './landing.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LandingPage {}
