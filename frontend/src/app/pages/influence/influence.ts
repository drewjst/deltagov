import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'app-influence-page',
  standalone: true,
  templateUrl: './influence.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class InfluencePage {}
