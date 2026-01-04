import { ChangeDetectionStrategy, Component, inject, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ScrollingModule } from '@angular/cdk/scrolling';
import { LucideAngularModule, GitCompare, FileText, Calendar, User } from 'lucide-angular';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmTypographyImports } from '@spartan-ng/helm/typography';
import { DiffLine, LivingBillStore } from './living-bill.store';

@Component({
  selector: 'app-living-bill',
  standalone: true,
  imports: [
    FormsModule,
    ScrollingModule,
    LucideAngularModule,
    BrnSelectImports,
    HlmSelectImports,
    HlmTypographyImports,
  ],
  templateUrl: './living-bill.html',
  styleUrl: './living-bill.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LivingBill implements OnInit {
  protected readonly store = inject(LivingBillStore);

  protected readonly GitCompareIcon = GitCompare;
  protected readonly FileTextIcon = FileText;
  protected readonly CalendarIcon = Calendar;
  protected readonly UserIcon = User;

  ngOnInit(): void {
    // Initialize with mock data until service calls are implemented
    this.store.setVersions([
      { id: 'v1', label: 'Version 1 (Dec 1)', date: '2024-12-01', contentHash: 'abc123' },
      { id: 'v2', label: 'Version 2 (Dec 10)', date: '2024-12-10', contentHash: 'def456' },
      { id: 'v3', label: 'Version 3 (Dec 15)', date: '2024-12-15', contentHash: 'ghi789' },
      { id: 'v4', label: 'Version 4 (Dec 20)', date: '2024-12-20', contentHash: 'jkl012' },
    ]);
  }

  protected onFromVersionChange(versionId: string): void {
    this.store.selectFromVersion(versionId);
  }

  protected onToVersionChange(versionId: string): void {
    this.store.selectToVersion(versionId);
  }

  protected trackByLine(_: number, item: DiffLine): string {
    return `${item.lineNumber}-${item.type}-${item.text.substring(0, 20)}`;
  }
}
