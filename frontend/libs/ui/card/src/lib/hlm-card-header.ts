import { Directive } from '@angular/core';
import { classes } from '@spartan-ng/helm/utils';

@Directive({
  selector: '[hlmCardHeader],hlm-card-header',
  host: {
    role: 'heading',
  },
})
export class HlmCardHeader {
  constructor() {
    classes(() => 'mb-2');
  }
}
