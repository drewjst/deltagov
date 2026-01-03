import { Directive } from '@angular/core';
import { BrnTooltip } from '@spartan-ng/brain/tooltip';

@Directive({
  selector: '[hlmTooltip]',
  standalone: true,
  hostDirectives: [BrnTooltip],
})
export class HlmTooltip {}
