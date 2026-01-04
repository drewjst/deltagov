import { ChangeDetectionStrategy, Component, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideSearch } from '@ng-icons/lucide';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { Workspace } from '../../components/workspace/workspace';

@Component({
  selector: 'app-bills-page',
  standalone: true,
  imports: [NgIcon, HlmIcon, Workspace],
  providers: [provideIcons({ lucideSearch })],
  templateUrl: './bills.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class BillsPage {
  protected readonly searchQuery = signal('');

  onSearch(): void {
    // TODO: Implement search functionality
  }
}
