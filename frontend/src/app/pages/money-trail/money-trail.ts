import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'app-money-trail-page',
  standalone: true,
  templateUrl: './money-trail.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MoneyTrailPage {}
