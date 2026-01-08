import { ChangeDetectionStrategy, Component } from '@angular/core';
import { InfluenceComponent } from '../../components/influence/influence.component';

@Component({
  selector: 'app-influence-page',
  standalone: true,
  templateUrl: './influence.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [InfluenceComponent],
})
export class InfluencePage {}
