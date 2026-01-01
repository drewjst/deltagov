import { Component } from '@angular/core';
import { LucideAngularModule, GitCompare, FileText, Calendar, User } from 'lucide-angular';

@Component({
  selector: 'app-living-bill',
  standalone: true,
  imports: [LucideAngularModule],
  templateUrl: './living-bill.html',
  styleUrl: './living-bill.scss'
})
export class LivingBill {
  protected readonly GitCompareIcon = GitCompare;
  protected readonly FileTextIcon = FileText;
  protected readonly CalendarIcon = Calendar;
  protected readonly UserIcon = User;
}
