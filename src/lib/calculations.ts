import { type ElementKey } from './constants';

export type ElementalBalance = Record<ElementKey, number>;

export type Profile = {
  dominant: ElementKey;
  lp: number;
  balance: ElementalBalance;
};

export function reduceToBase(n: number): number {
  if (n === 11) return 2;
  if (n === 22) return 4;
  if (n === 33) return 6;
  return n;
}
