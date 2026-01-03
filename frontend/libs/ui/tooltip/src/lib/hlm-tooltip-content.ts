import { Directive, computed, input, inject } from '@angular/core';
import { BrnTooltipTrigger } from '@spartan-ng/brain/tooltip';
import { hlm } from '@spartan-ng/helm/utils';
import type { ClassValue } from 'clsx';

const tooltipContentClasses =
  'z-50 overflow-hidden rounded-md bg-primary px-3 py-1.5 text-xs text-primary-foreground animate-in fade-in-0 zoom-in-95 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2';

@Directive({
  selector: '[hlmTooltipContent]',
  standalone: true,
})
export class HlmTooltipContent {
  private readonly _brnTooltipTrigger = inject(BrnTooltipTrigger, { optional: true });

  public readonly userClass = input<ClassValue>('', { alias: 'class' });

  protected readonly _computedClass = computed(() => hlm(tooltipContentClasses, this.userClass()));

  constructor() {
    this._brnTooltipTrigger?.setTooltipContentClasses(tooltipContentClasses);
  }
}
