import { Directive } from '@angular/core';
import { classes } from '@spartan-ng/helm/utils';

@Directive({
	selector: '[hlmCardContent],hlm-card-content'
})
export class HlmCardContent {
	constructor() {
		classes(() => 'text-sm');
	}
}
