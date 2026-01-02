import { Directive } from '@angular/core';
import { classes } from '@spartan-ng/helm/utils';

@Directive({
	selector: '[hlmCardTitle],hlm-card-title'
})
export class HlmCardTitle {
	constructor() {
		classes(() => 'text-sm font-medium');
	}
}
