import { useEffect, useState } from 'react';
import { BarChart } from './BarChart';
import {
  CareerIcon,
  HealthIcon,
  LockIcon,
  MarkIcon,
  RomanceIcon,
  StarIcon,
  WealthIcon,
} from './icons';
import {
  ELEMENT_META,
  ELEMENT_TEXT,
  type ElementKey,
  type TabKey,
} from '../lib/constants';
import {
  reduceToBase,
  type ElementalBalance,
  type Profile,
} from '../lib/calculations';
import { LIFE_PATH_TEXT } from '../lib/constants';
import { fetchProfile } from '../lib/api';

const TAB_CONFIG: { key: TabKey; label: string; icon: React.ElementType }[] = [
  { key: 'romance', label: 'Romance', icon: RomanceIcon },
  { key: 'health', label: 'Health', icon: HealthIcon },
  { key: 'career', label: 'Career', icon: CareerIcon },
  { key: 'wealth', label: 'Wealth', icon: WealthIcon },
];

export type DashboardScreenProps = {
  birthdate: string;
  birthtime?: string;
  gender?: string;
  unlocked: boolean;
  activeTab: TabKey;
  onTabChange: (tab: TabKey) => void;
  onUnlock: () => void;
  onReset: () => void;
};

export function DashboardScreen({
  birthdate,
  birthtime,
  gender,
  unlocked,
  activeTab,
  onTabChange,
  onUnlock,
  onReset,
}: DashboardScreenProps) {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    setLoading(true);
    setError(false);
    console.log('[dashboard] fetchProfile called with:', { birthdate, birthtime, gender });
    fetchProfile(birthdate, birthtime, gender)
      .then(setProfile)
      .catch(() => setError(true))
      .finally(() => setLoading(false));
  }, [birthdate, birthtime, gender]);

  return (
    <section className="flex flex-col min-h-[640px] pb-6 gap-3.5">
      <div className="px-5 pt-6">
        <div className="flex justify-between items-start mb-5">
          <MarkIcon className="w-[26px] h-[26px] opacity-90" />
          <button
            type="button"
            onClick={onReset}
            className="bg-transparent border-0 text-ink-faint font-mono text-[11px] cursor-pointer tracking-wide hover:text-ink-dim"
          >
            START OVER
          </button>
        </div>
      </div>

      <div className="bg-card border border-hairline rounded-[16px] shadow-[0_1px_0_rgba(255,255,255,0.6)_inset,0_8px_20px_-14px_rgba(60,50,30,0.25)] mx-5 mb-4 px-5 py-5 pb-[18px]">
        <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2">
          CORE PROFILE
        </div>
        {loading && (
          <p className="text-ink-dim text-sm">Loading your profile…</p>
        )}
        {error && (
          <p className="text-ink-dim text-sm">Could not load profile. Make sure the server is running.</p>
        )}
        {profile && (
          <>
            <h2 className="font-display font-medium text-[21px] leading-[1.3] tracking-tight mb-[18px]">
              Balanced <span className="italic">{ELEMENT_META[profile.dominant].label}</span> Element &nbsp;/&nbsp; Life Path {profile.lp}
            </h2>
            <BarChart balance={profile.balance} />
          </>
        )}
      </div>

      <div className="px-5 mb-1">
        <nav className="flex gap-1 bg-card border border-hairline rounded-full p-[5px]">
          {TAB_CONFIG.map(({ key, label, icon: Icon }) => {
            const locked = (key === 'career' || key === 'wealth') && !unlocked;
            const isActive = activeTab === key;
            return (
              <button
                key={key}
                type="button"
                onClick={() => onTabChange(key)}
                className={`flex-1 bg-transparent border-0 font-body font-bold text-xs pt-2.5 pb-2.5 rounded-full flex flex-col items-center gap-[5px] transition-colors duration-200 relative ${
                  isActive
                    ? 'text-ink bg-card-2 shadow-[0_2px_8px_-4px_rgba(60,50,30,0.3)]'
                    : `text-ink-faint hover:text-ink-dim ${locked ? 'text-ink-faint' : ''}`
                }`}
              >
                <Icon
                  className={`w-4 h-4 transition-opacity duration-200 ${
                    isActive ? 'opacity-100' : 'opacity-65'
                  }`}
                />
                {label}
                {locked && (
                  <span className="absolute top-1.5 right-[22%] w-1.5 h-1.5 rounded-full bg-earth animate-pulse-glow" />
                )}
              </button>
            );
          })}
        </nav>
      </div>

      {profile && (
        <div className="px-5 pt-3.5">
          {TAB_CONFIG.map(({ key }) => (
            <TabPanel
              key={key}
              tab={key}
              isActive={activeTab === key}
              dominant={profile.dominant}
              lp={profile.lp}
              balance={profile.balance}
              locked={(key === 'career' || key === 'wealth') && !unlocked}
              onUnlock={onUnlock}
            />
          ))}
        </div>
      )}
    </section>
  );
}

type TabPanelProps = {
  tab: TabKey;
  isActive: boolean;
  dominant: ElementKey;
  lp: number;
  balance: ElementalBalance;
  locked: boolean;
  onUnlock: () => void;
};

function TabPanel({ tab, isActive, dominant, lp, balance, locked, onUnlock }: TabPanelProps) {
  const domainLabel = {
    romance: 'Romance',
    health: 'Health',
    career: 'Career',
    wealth: 'Wealth',
  }[tab];

  const lpText = LIFE_PATH_TEXT[lp] ?? LIFE_PATH_TEXT[reduceToBase(lp)];

  if (!isActive) return null;

  if (locked) {
    return (
      <div className="animate-slide-in">
        <div className="bg-card border border-hairline rounded-[16px] shadow-[0_1px_0_rgba(255,255,255,0.6)_inset,0_8px_20px_-14px_rgba(60,50,30,0.25)] p-[22px] pb-6">
          <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2.5">
            {domainLabel.toUpperCase()} · LOCKED
          </div>
          <div className="font-display text-[19px] font-medium leading-[1.35] mb-3.5">
            Your Hour Pillar decides this one.
          </div>
          <div className="relative rounded-[10px] overflow-hidden border border-hairline bg-card-2 p-[18px] mb-[18px]">
            <BarChart balance={balance} blurred animate={false} />
            <div className="absolute inset-0 flex flex-col items-center justify-center gap-2.5 bg-gradient-to-b from-[rgba(243,239,230,0.35)] to-[rgba(243,239,230,0.82)]">
              <div className="w-[38px] h-[38px] rounded-full bg-card-2 border border-hairline-strong flex items-center justify-center text-earth">
                <LockIcon className="w-4 h-4" />
              </div>
              <p className="font-mono text-[11px] text-ink-dim tracking-wide uppercase">
                Full chart hidden
              </p>
            </div>
          </div>
          <button
            type="button"
            onClick={onUnlock}
            className="unlock-cta"
          >
            <StarIcon className="w-[15px] h-[15px]" />
            Unlock deep insights with birth time & city
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="animate-slide-in">
      <div className="bg-card border border-hairline rounded-[16px] shadow-[0_1px_0_rgba(255,255,255,0.6)_inset,0_8px_20px_-14px_rgba(60,50,30,0.25)] p-[22px] pb-6">
        <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2.5">
          {domainLabel.toUpperCase()}
        </div>
        <div className="font-display text-[19px] font-medium leading-[1.35] mb-3.5">
          {ELEMENT_META[dominant].label} energy, read for {domainLabel.toLowerCase()}
        </div>
        <div className="flex gap-2.5 mb-[18px]">
          <div className="flex-1 bg-card-2 border border-hairline rounded-[10px] px-3.5 py-3">
            <div className="font-mono text-[10px] text-ink-faint uppercase tracking-[0.06em] mb-1">
              Element
            </div>
            <div className="font-display text-base">{ELEMENT_META[dominant].label}</div>
          </div>
          <div className="flex-1 bg-card-2 border border-hairline rounded-[10px] px-3.5 py-3">
            <div className="font-mono text-[10px] text-ink-faint uppercase tracking-[0.06em] mb-1">
              Life Path
            </div>
            <div className="font-display text-base">{lp}</div>
          </div>
        </div>
        <p className="text-ink-dim text-[14.5px] leading-[1.7] mb-[18px]">
          {ELEMENT_TEXT[dominant][tab]}
        </p>
        <p className="text-ink-dim text-[14.5px] leading-[1.7]">
          As a Life Path {lp}, you're {lpText}
        </p>
      </div>
    </div>
  );
}
