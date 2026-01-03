export * from './lib/hlm-tooltip';
export * from './lib/hlm-tooltip-trigger';
export * from './lib/hlm-tooltip-content';

import { HlmTooltip } from './lib/hlm-tooltip';
import { HlmTooltipTrigger } from './lib/hlm-tooltip-trigger';
import { HlmTooltipContent } from './lib/hlm-tooltip-content';

export const HlmTooltipImports = [HlmTooltip, HlmTooltipTrigger, HlmTooltipContent] as const;
