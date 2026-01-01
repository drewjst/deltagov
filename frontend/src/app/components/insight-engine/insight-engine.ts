import { Component } from '@angular/core';
import { LucideAngularModule, Brain, TrendingUp, DollarSign, AlertTriangle, Lightbulb, MessageSquare } from 'lucide-angular';

@Component({
  selector: 'app-insight-engine',
  standalone: true,
  imports: [LucideAngularModule],
  templateUrl: './insight-engine.html',
  styleUrl: './insight-engine.scss'
})
export class InsightEngine {
  protected readonly BrainIcon = Brain;
  protected readonly TrendingUpIcon = TrendingUp;
  protected readonly DollarSignIcon = DollarSign;
  protected readonly AlertTriangleIcon = AlertTriangle;
  protected readonly LightbulbIcon = Lightbulb;
  protected readonly MessageSquareIcon = MessageSquare;
}
