import { type ElementKey } from './constants';

export type ElementalBalance = Record<ElementKey, number>;

export type Pillar = {
  stem: string;
  stemElement: string;
  stemYinYang: string;
  branch: string;
  branchElement: string;
  hiddenStems: string[];
};

export type DayMaster = {
  stem: string;
  element: string;
  yinYang: string;
};

export type SeasonInfo = {
  name: string;
  branch: string;
  element: string;
};

export type YongShenInfo = {
  element: string;
  reason: string;
};

export type TenGodEntry = {
  element: string;
  god: string;
};

export type LuckCycle = {
  startAge: number;
  endAge: number;
  stem: string;
  branch: string;
  element: string;
};

export type DayMasterProfileInfo = {
  archetype: string;
  coreTraits: string[];
  summary: string;
};

export type TenGodInsight = {
  god: string;
  category?: string;
  description: string;
};

export type Profile = {
  dominant: ElementKey;
  lp: number;
  balance: ElementalBalance;
  pillars?: Record<string, Pillar>;
  dayMaster?: DayMaster;
  dmStrength?: string;
  season?: SeasonInfo;
  yongShen?: YongShenInfo;
  favorableElements?: string[];
  unfavorableElements?: string[];
  tenGods?: Record<string, TenGodEntry>;
  luckCycles?: LuckCycle[];
  tabText?: Record<string, string>;
  dayMasterProfile?: DayMasterProfileInfo;
  tenGodInsights?: Record<string, TenGodInsight>;
  textModifiers?: string[];
  lifePathText?: string;
};

export function reduceToBase(n: number): number {
  if (n === 11) return 2;
  if (n === 22) return 4;
  if (n === 33) return 6;
  return n;
}
