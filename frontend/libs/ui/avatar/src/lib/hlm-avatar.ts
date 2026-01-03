import { Component, computed, input } from '@angular/core';
import { BrnAvatarImports } from '@spartan-ng/brain/avatar';
import { hlm } from '@spartan-ng/helm/utils';
import { type VariantProps, cva } from 'class-variance-authority';
import type { ClassValue } from 'clsx';

export const avatarVariants = cva('relative flex shrink-0 overflow-hidden rounded-full', {
  variants: {
    size: {
      default: 'h-10 w-10',
      sm: 'h-8 w-8',
      lg: 'h-12 w-12',
      xl: 'h-16 w-16',
    },
  },
  defaultVariants: {
    size: 'default',
  },
});

export type AvatarVariants = VariantProps<typeof avatarVariants>;

@Component({
  selector: 'hlm-avatar',
  standalone: true,
  imports: [BrnAvatarImports],
  host: {
    '[class]': '_computedClass()',
  },
  template: `
    <brn-avatar>
      <ng-content select="[hlmAvatarImage]" />
      <ng-content select="[hlmAvatarFallback]" />
    </brn-avatar>
  `,
})
export class HlmAvatar {
  public readonly userClass = input<ClassValue>('', { alias: 'class' });
  public readonly size = input<AvatarVariants['size']>('default');

  protected readonly _computedClass = computed(() =>
    hlm(avatarVariants({ size: this.size() }), this.userClass()),
  );
}
