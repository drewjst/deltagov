import { Directive, computed, input } from '@angular/core';
import { BrnAvatarFallback } from '@spartan-ng/brain/avatar';
import { hlm } from '@spartan-ng/helm/utils';
import type { ClassValue } from 'clsx';

@Directive({
  selector: '[hlmAvatarFallback]',
  standalone: true,
  hostDirectives: [BrnAvatarFallback],
  host: {
    '[class]': '_computedClass()',
  },
})
export class HlmAvatarFallback {
  public readonly userClass = input<ClassValue>('', { alias: 'class' });

  protected readonly _computedClass = computed(() =>
    hlm(
      'flex h-full w-full items-center justify-center rounded-full bg-muted text-muted-foreground',
      this.userClass(),
    ),
  );
}
