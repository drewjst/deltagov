import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'app-history-page',
  standalone: true,
  templateUrl: './history.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HistoryPage {}
