export * from './lib/hlm-avatar';
export * from './lib/hlm-avatar-image';
export * from './lib/hlm-avatar-fallback';

import { HlmAvatar } from './lib/hlm-avatar';
import { HlmAvatarImage } from './lib/hlm-avatar-image';
import { HlmAvatarFallback } from './lib/hlm-avatar-fallback';

export const HlmAvatarImports = [HlmAvatar, HlmAvatarImage, HlmAvatarFallback] as const;
