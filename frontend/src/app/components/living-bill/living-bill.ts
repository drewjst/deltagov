import { Component, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { LucideAngularModule, GitCompare, FileText, Calendar, User } from 'lucide-angular';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';

@Component({
  selector: 'app-living-bill',
  standalone: true,
  imports: [FormsModule, LucideAngularModule, BrnSelectImports, HlmSelectImports],
  templateUrl: './living-bill.html',
  styleUrl: './living-bill.scss'
})
export class LivingBill {
  protected readonly GitCompareIcon = GitCompare;
  protected readonly FileTextIcon = FileText;
  protected readonly CalendarIcon = Calendar;
  protected readonly UserIcon = User;

  // Version selection
  protected readonly fromVersion = signal('v3');
  protected readonly toVersion = signal('v4');

  protected readonly versions = [
    { value: 'v1', label: 'Version 1 (Dec 1)' },
    { value: 'v2', label: 'Version 2 (Dec 10)' },
    { value: 'v3', label: 'Version 3 (Dec 15)' },
    { value: 'v4', label: 'Version 4 (Dec 20)' },
  ];
}
