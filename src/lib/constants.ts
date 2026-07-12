export type ElementKey = 'wood' | 'fire' | 'earth' | 'metal' | 'water';
export type Screen = 'onboard' | 'deep' | 'dashboard';

export const ELEMENTS: ElementKey[] = ['wood', 'fire', 'earth', 'metal', 'water'];

export type ElementMeta = {
  label: string;
  short: string;
  color: string;
};

export const ELEMENT_META: Record<ElementKey, ElementMeta> = {
  wood: { label: 'Wood', short: 'WOOD', color: '#5E8A4C' },
  fire: { label: 'Fire', short: 'FIRE', color: '#BD5B3B' },
  earth: { label: 'Earth', short: 'EARTH', color: '#A97C34' },
  metal: { label: 'Metal', short: 'METAL', color: '#6E7B8C' },
  water: { label: 'Water', short: 'WATER', color: '#35708E' },
};
