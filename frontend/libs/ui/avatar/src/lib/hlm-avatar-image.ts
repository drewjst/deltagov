import { Directive, computed, input } from '@angular/core';
import { BrnAvatarImage } from '@spartan-ng/brain/avatar';
import { hlm } from '@spartan-ng/helm/utils';
import type { ClassValue } from 'clsx';

@Directive({
  selector: 'img[hlmAvatarImage]',
  standalone: true,
  hostDirectives: [BrnAvatarImage],
  host: {
    '[class]': '_computedClass()',
  },
})
export class HlmAvatarImage {
  public readonly userClass = input<ClassValue>('', { alias: 'class' });

  protected readonly _computedClass = computed(() =>
    hlm('aspect-square h-full w-full object-cover', this.userClass()),
  );
}
