import { ChangeDetectionStrategy, Component, effect, inject, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ScrollingModule } from '@angular/cdk/scrolling';
import { LucideAngularModule, GitCompare, FileText, Calendar, User, Loader2 } from 'lucide-angular';
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
  protected readonly LoaderIcon = Loader2;

  constructor() {
    // Effect to load diff when versions change
    effect(() => {
      const fromVersion = this.store.selectedFromVersion();
      const toVersion = this.store.selectedToVersion();
      const bill = this.store.bill();

      if (bill && fromVersion && toVersion && fromVersion !== toVersion) {
        this.store.loadDiff({
          billId: bill.id,
          fromVersionId: fromVersion,
          toVersionId: toVersion,
        });
      }
    });
  }

  ngOnInit(): void {
    // Load H.R. 1 - The One Big Beautiful Bill on component init
    this.store.loadHR1();
  }

  protected onFromVersionChange(versionId: number): void {
    this.store.selectFromVersion(versionId);
  }

  protected onToVersionChange(versionId: number): void {
    this.store.selectToVersion(versionId);
  }

  protected trackByLine(_: number, item: DiffLine): string {
    return `${item.lineNumber}-${item.type}-${item.text.substring(0, 20)}`;
  }
}
