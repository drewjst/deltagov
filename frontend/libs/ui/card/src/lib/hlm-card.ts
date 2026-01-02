import { Directive } from '@angular/core';
import { classes } from '@spartan-ng/helm/utils';

@Directive({
	selector: '[hlmCard],hlm-card',
	host: {
		role: 'group',
	},
})
export class HlmCard {
	constructor() {
		classes(() => 'p-4 rounded-lg border border-border bg-background');
	}
}
