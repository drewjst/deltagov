import { ChangeDetectionStrategy, Component } from '@angular/core';
import { LivingBill } from '../living-bill/living-bill';
import { InsightEngine } from '../insight-engine/insight-engine';

@Component({
  selector: 'app-workspace',
  standalone: true,
  imports: [LivingBill, InsightEngine],
  templateUrl: './workspace.html',
  styleUrl: './workspace.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class Workspace {}
