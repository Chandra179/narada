import { ELEMENTS, type ElementKey } from './constants';

export function lifePathNumber(dateStr: string): number {
  const digits = dateStr.replace(/-/g, '').split('').map(Number);
  const sum = digits.reduce((a, b) => a + b, 0);

  const reduce = (n: number): number => {
    while (n > 9 && n !== 11 && n !== 22 && n !== 33) {
      n = String(n)
        .split('')
        .reduce((acc, digit) => acc + Number(digit), 0);
    }
    return n;
  };

  return reduce(sum);
}

export function elementFromYear(year: number): ElementKey {
  const last = year % 10;
  if (last === 0 || last === 1) return 'metal';
  if (last === 2 || last === 3) return 'water';
  if (last === 4 || last === 5) return 'wood';
  if (last === 6 || last === 7) return 'fire';
  return 'earth';
}

export function seedFromString(str: string): number {
  let h = 0;
  for (let i = 0; i < str.length; i++) {
    h = (h * 31 + str.charCodeAt(i)) >>> 0;
  }
  return h;
}

export type ElementalBalance = Record<ElementKey, number>;

export function elementalBalance(
  dateStr: string,
  dominant: ElementKey
): ElementalBalance {
  const seed = seedFromString(dateStr);
  const rand = (i: number) => ((seed >> (i * 3)) % 17) / 17;

  const vals = {} as Record<ElementKey, number>;
  ELEMENTS.forEach((el, i) => {
    vals[el] = 8 + Math.floor(rand(i) * 18);
  });
  vals[dominant] += 34;

  const total = Object.values(vals).reduce((a, b) => a + b, 0);
  ELEMENTS.forEach((el) => {
    vals[el] = Math.round((vals[el] / total) * 100);
  });

  const diff = 100 - Object.values(vals).reduce((a, b) => a + b, 0);
  vals[dominant] += diff;

  return vals;
}

export type Profile = {
  dominant: ElementKey;
  lp: number;
  balance: ElementalBalance;
};

export function computeProfile(birthdate: string): Profile {
  const d = new Date(`${birthdate}T00:00:00`);
  const year = d.getFullYear();
  const dominant = elementFromYear(year);
  const lp = lifePathNumber(birthdate);
  const balance = elementalBalance(birthdate, dominant);
  return { dominant, lp, balance };
}

export function reduceToBase(n: number): number {
  if (n === 11) return 2;
  if (n === 22) return 4;
  if (n === 33) return 6;
  return n;
}
