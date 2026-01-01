import { Component, signal } from '@angular/core';
import { LucideAngularModule, Search, Sparkles, User, ChevronDown } from 'lucide-angular';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [LucideAngularModule],
  templateUrl: './header.html',
  styleUrl: './header.scss'
})
export class Header {
  protected readonly searchQuery = signal('');

  // Lucide icons
  protected readonly SearchIcon = Search;
  protected readonly SparklesIcon = Sparkles;
  protected readonly UserIcon = User;
  protected readonly ChevronDownIcon = ChevronDown;
}
