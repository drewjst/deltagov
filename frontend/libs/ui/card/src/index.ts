export { HlmCard } from './lib/hlm-card';
export { HlmCardHeader } from './lib/hlm-card-header';
export { HlmCardTitle } from './lib/hlm-card-title';
export { HlmCardContent } from './lib/hlm-card-content';

import { HlmCard } from './lib/hlm-card';
import { HlmCardHeader } from './lib/hlm-card-header';
import { HlmCardTitle } from './lib/hlm-card-title';
import { HlmCardContent } from './lib/hlm-card-content';

export const HlmCardImports = [HlmCard, HlmCardHeader, HlmCardTitle, HlmCardContent] as const;
