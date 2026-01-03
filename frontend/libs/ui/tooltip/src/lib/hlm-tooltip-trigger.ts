import { Directive, input } from '@angular/core';
import { BrnTooltipTrigger } from '@spartan-ng/brain/tooltip';

@Directive({
  selector: '[hlmTooltipTrigger]',
  standalone: true,
  hostDirectives: [
    {
      directive: BrnTooltipTrigger,
      inputs: [
        'brnTooltipTrigger: hlmTooltipTrigger',
        'position',
        'showDelay',
        'hideDelay',
        'brnTooltipDisabled',
        'exitAnimationDuration',
      ],
    },
  ],
})
export class HlmTooltipTrigger {
  public readonly hlmTooltipTrigger = input<string | null>(null);
}
