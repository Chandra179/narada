import { useEffect, useState, type ReactElement } from 'react';
import { BarChart } from './BarChart';
import { LockIcon, MarkIcon, StarIcon } from './icons';
import { getTabIcon } from './getTabIcon';
import { ELEMENT_META, type ElementKey } from '../lib/constants';
import {
  type Profile,
} from '../lib/calculations';
import { fetchProfile } from '../lib/api';
import { useConfig, type TabConfig } from '../hooks/useConfig';

export type DashboardScreenProps = {
  birthdate: string;
  birthtime?: string;
  gender?: string;
  unlocked: boolean;
  activeTab: string;
  onTabChange: (tab: string) => void;
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
  const { config, loading: configLoading } = useConfig();
  const [profile, setProfile] = useState<Profile | null>(null);
  const [profileLoading, setProfileLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    setProfileLoading(true);
    setError(false);
    fetchProfile(birthdate, birthtime, gender)
      .then(setProfile)
      .catch(() => setError(true))
      .finally(() => setProfileLoading(false));
  }, [birthdate, birthtime, gender]);

  if (configLoading || profileLoading) {
    return (
      <section className="flex flex-col min-h-[640px] pb-6 gap-3.5">
        <div className="px-5 pt-6">
          <MarkIcon className="w-[26px] h-[26px] opacity-90" />
        </div>
        <div className="px-5">
          <p className="text-ink-dim text-sm">Loading your profile…</p>
        </div>
      </section>
    );
  }

  if (error) {
    return (
      <section className="flex flex-col min-h-[640px] pb-6 gap-3.5">
        <div className="px-5 pt-6">
          <MarkIcon className="w-[26px] h-[26px] opacity-90" />
        </div>
        <div className="px-5">
          <p className="text-ink-dim text-sm">Could not load profile. Make sure the server is running.</p>
        </div>
      </section>
    );
  }

  const tabs = config?.tabs ?? [];

  // Default locked tabs: all except the first two
  const alwaysUnlockedTabs = tabs.length >= 2 ? tabs.slice(0, 2).map(t => t.key) : [];

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

      {profile && (
        <>
          {/* Core Profile Card */}
          <CoreProfileCard profile={profile} />

          {/* Season & Strength Card */}
          {profile.season && profile.season.name && (
            <div className="bg-card border border-hairline rounded-[16px] mx-5 px-5 py-4">
              <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2">
                SEASONAL CONTEXT
              </div>
              <p className="text-ink-dim text-[13px] leading-relaxed">
                Born in <strong>{profile.season.name}</strong> ({profile.season.branch}), when <strong>{profile.season.element}</strong> energy peaks.
                Your Day Master is <strong>{profile.dmStrength}</strong>.
              </p>
              {profile.yongShen && profile.yongShen.element && (
                <p className="text-ink-dim text-[13px] leading-relaxed mt-1">
                  <strong>Yong Shen (Useful God):</strong> {profile.yongShen.reason}
                </p>
              )}
              {profile.favorableElements && profile.favorableElements.length > 0 && (
                <div className="flex gap-2 mt-2 flex-wrap">
                  {profile.favorableElements.map(el => (
                    <span key={el} className="text-[11px] font-mono px-2 py-0.5 rounded-full bg-green-50 text-green-800 border border-green-200">
                      {el} ✓
                    </span>
                  ))}
                  {profile.unfavorableElements?.map(el => (
                    <span key={el} className="text-[11px] font-mono px-2 py-0.5 rounded-full bg-red-50 text-red-800 border border-red-200">
                      {el} ✗
                    </span>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* Pillars Card */}
          {profile.pillars && (
            <PillarsCard pillars={profile.pillars} />
          )}

          {/* Day Master Profile */}
          {profile.dayMasterProfile && (
            <div className="bg-card border border-hairline rounded-[16px] mx-5 px-5 py-4">
              <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2">
                YOUR ARCHETYPE
              </div>
              <h3 className="font-display font-medium text-[19px] leading-tight mb-1">
                {profile.dayMasterProfile.archetype}
              </h3>
              {profile.dayMasterProfile.coreTraits && (
                <div className="flex gap-1.5 flex-wrap mb-2">
                  {profile.dayMasterProfile.coreTraits.map(trait => (
                    <span key={trait} className="text-[11px] font-mono px-2 py-0.5 rounded-full bg-card-2 border border-hairline text-ink-dim">
                      {trait}
                    </span>
                  ))}
                </div>
              )}
              {profile.dayMasterProfile.summary && (
                <p className="text-ink-dim text-[13px] leading-relaxed">
                  {profile.dayMasterProfile.summary}
                </p>
              )}
            </div>
          )}

          {/* Tab Navigation */}
          <div className="px-5 mb-1">
            <nav className="flex gap-1 bg-card border border-hairline rounded-full p-[5px]">
              {tabs.map((tab: TabConfig) => {
                const locked = !alwaysUnlockedTabs.includes(tab.key) && !unlocked;
                const isActive = activeTab === tab.key;
                const Icon = getTabIcon(tab.icon);
                return (
                  <button
                    key={tab.key}
                    type="button"
                    onClick={() => onTabChange(tab.key)}
                    className={`flex-1 bg-transparent border-0 font-body font-bold text-xs pt-2.5 pb-2.5 rounded-full flex flex-col items-center gap-[5px] transition-colors duration-200 relative ${
                      isActive
                        ? 'text-ink bg-card-2 shadow-[0_2px_8px_-4px_rgba(60,50,30,0.3)]'
                        : `text-ink-faint hover:text-ink-dim ${locked ? 'text-ink-faint' : ''}`
                    }`}
                  >
                    <Icon className={`w-4 h-4 transition-opacity duration-200 ${isActive ? 'opacity-100' : 'opacity-65'}`} />
                    {tab.label}
                    {locked && (
                      <span className="absolute top-1.5 right-[22%] w-1.5 h-1.5 rounded-full bg-earth animate-pulse-glow" />
                    )}
                  </button>
                );
              })}
            </nav>
          </div>

          {/* Tab Content */}
          <div className="px-5 pt-3.5">
            {tabs.map((tab: TabConfig) => (
              <TabPanel
                key={tab.key}
                tab={tab}
                isActive={activeTab === tab.key}
                profile={profile}
                locked={!alwaysUnlockedTabs.includes(tab.key) && !unlocked}
                onUnlock={onUnlock}
              />
            ))}
          </div>
        </>
      )}
    </section>
  );
}

// =============================================================================
// Sub-components
// =============================================================================

function CoreProfileCard({ profile }: { profile: Profile }) {
  const dm = profile.dayMaster;
  return (
    <div className="bg-card border border-hairline rounded-[16px] shadow-[0_1px_0_rgba(255,255,255,0.6)_inset,0_8px_20px_-14px_rgba(60,50,30,0.25)] mx-5 mb-4 px-5 py-5 pb-[18px]">
      <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2">
        CORE PROFILE
      </div>
      {dm && (
        <div className="flex items-baseline gap-2 mb-1">
          <h2 className="font-display font-medium text-[21px] leading-[1.3] tracking-tight">
            {dm.stem} · {dm.element} · <span className="italic">{dm.yinYang}</span>
          </h2>
          <span className="font-mono text-[11px] text-ink-faint">
            LP {profile.lp}
          </span>
        </div>
      )}
      <BarChart balance={profile.balance} />
    </div>
  );
}

function PillarsCard({ pillars }: { pillars: Record<string, import('../lib/calculations').Pillar> }) {
  const keys = ['year', 'month', 'day', 'hour'] as const;
  const labels: Record<string, string> = { year: 'Year', month: 'Month', day: 'Day', hour: 'Hour' };

  return (
    <div className="bg-card border border-hairline rounded-[16px] mx-5 mb-4 px-5 py-4">
      <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2.5">
        FOUR PILLARS
      </div>
      <div className="grid grid-cols-4 gap-2">
        {keys.map((k) => {
          const p = pillars[k];
          if (!p) return null;
          return (
            <div key={k} className="bg-card-2 border border-hairline rounded-[10px] p-2.5 text-center">
              <div className="font-mono text-[9px] text-ink-faint uppercase tracking-wider mb-1">
                {labels[k]}
              </div>
              <div className="font-display text-[15px] font-medium leading-tight">
                {p.stem}{p.branch}
              </div>
              <div className="font-mono text-[10px] text-ink-dim">
                {p.stemElement}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

function TabPanel({
  tab,
  isActive,
  profile,
  locked,
  onUnlock,
}: {
  tab: TabConfig;
  isActive: boolean;
  profile: Profile;
  locked: boolean;
  onUnlock: () => void;
}): ReactElement | null {
  if (!isActive) return null;

  if (locked) {
    return (
      <div className="animate-slide-in">
        <div className="bg-card border border-hairline rounded-[16px] shadow-[0_1px_0_rgba(255,255,255,0.6)_inset,0_8px_20px_-14px_rgba(60,50,30,0.25)] p-[22px] pb-6">
          <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2.5">
            {tab.label.toUpperCase()} · LOCKED
          </div>
          <div className="font-display text-[19px] font-medium leading-[1.35] mb-3.5">
            Your Hour Pillar decides this one.
          </div>
          <div className="relative rounded-[10px] overflow-hidden border border-hairline bg-card-2 p-[18px] mb-[18px]">
            <BarChart balance={profile.balance} blurred animate={false} />
            <div className="absolute inset-0 flex flex-col items-center justify-center gap-2.5 bg-gradient-to-b from-[rgba(243,239,230,0.35)] to-[rgba(243,239,230,0.82)]">
              <div className="w-[38px] h-[38px] rounded-full bg-card-2 border border-hairline-strong flex items-center justify-center text-earth">
                <LockIcon className="w-4 h-4" />
              </div>
              <p className="font-mono text-[11px] text-ink-dim tracking-wide uppercase">
                Full chart hidden
              </p>
            </div>
          </div>
          <button type="button" onClick={onUnlock} className="unlock-cta">
            <StarIcon className="w-[15px] h-[15px]" />
            Unlock deep insights with birth time & city
          </button>
        </div>
      </div>
    );
  }

  const dominant = profile.dominant as ElementKey;
  const dm = profile.dayMaster;
  const dmProfile = profile.dayMasterProfile;

  return (
    <div className="animate-slide-in">
      {/* Tab Content Card */}
      <div className="bg-card border border-hairline rounded-[16px] shadow-[0_1px_0_rgba(255,255,255,0.6)_inset,0_8px_20px_-14px_rgba(60,50,30,0.25)] p-[22px] pb-6 mb-4">
        <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2.5">
          {tab.label.toUpperCase()}
        </div>
        <div className="font-display text-[19px] font-medium leading-[1.35] mb-3.5">
          {dm ? `${dm.stem} ${dm.element} · ${dmProfile?.archetype ?? ''}` : `${ELEMENT_META[dominant].label} energy`}, read for {tab.label.toLowerCase()}
        </div>
        <div className="flex gap-2.5 mb-[18px]">
          <div className="flex-1 bg-card-2 border border-hairline rounded-[10px] px-3.5 py-3">
            <div className="font-mono text-[10px] text-ink-faint uppercase tracking-[0.06em] mb-1">
              Element
            </div>
            <div className="font-display text-base">{ELEMENT_META[dominant].label}</div>
          </div>
          {dm && (
            <div className="flex-1 bg-card-2 border border-hairline rounded-[10px] px-3.5 py-3">
              <div className="font-mono text-[10px] text-ink-faint uppercase tracking-[0.06em] mb-1">
                Day Master
              </div>
              <div className="font-display text-base">{dm.stem}</div>
            </div>
          )}
          <div className="flex-1 bg-card-2 border border-hairline rounded-[10px] px-3.5 py-3">
            <div className="font-mono text-[10px] text-ink-faint uppercase tracking-[0.06em] mb-1">
              Life Path
            </div>
            <div className="font-display text-base">{profile.lp}</div>
          </div>
        </div>

        {/* Day Master tab text */}
        {profile.tabText?.[tab.key] && (
          <p className="text-ink-dim text-[14.5px] leading-[1.7] mb-[18px]">
            {profile.tabText[tab.key]}
          </p>
        )}

        {/* Life Path text */}
        {profile.lifePathText && (
          <p className="text-ink-dim text-[14.5px] leading-[1.7]">
            As a Life Path {profile.lp}, you're {profile.lifePathText}
          </p>
        )}
      </div>

      {/* Ten Gods Card */}
      {profile.tenGods && (
        <TenGodsCard tenGods={profile.tenGods} insights={profile.tenGodInsights} />
      )}

      {/* Clash Modifiers Card */}
      {profile.textModifiers && profile.textModifiers.length > 0 && (
        <div className="bg-card border border-hairline rounded-[16px] p-[22px] pb-6 mb-4">
          <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2.5">
            CHART MODIFIERS
          </div>
          <ul className="space-y-2">
            {profile.textModifiers.map((mod, i) => (
              <li key={i} className="text-ink-dim text-[13px] leading-relaxed flex gap-2">
                <span className="text-earth shrink-0 mt-0.5">•</span>
                {mod}
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Luck Cycles Card */}
      {profile.luckCycles && profile.luckCycles.length > 0 && (
        <LuckCyclesCard cycles={profile.luckCycles} />
      )}
    </div>
  );
}

function TenGodsCard({
  tenGods,
  insights,
}: {
  tenGods: Record<string, import('../lib/calculations').TenGodEntry>;
  insights?: Record<string, import('../lib/calculations').TenGodInsight>;
}) {
  const pillars = ['yearStem', 'monthStem', 'hourStem'] as const;
  const labels: Record<string, string> = {
    yearStem: 'Year',
    monthStem: 'Month',
    hourStem: 'Hour',
  };

  return (
    <div className="bg-card border border-hairline rounded-[16px] p-[22px] pb-6 mb-4">
      <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2.5">
        TEN GODS (十神)
      </div>
      <div className="space-y-3">
        {pillars.map((pillar) => {
          const tg = tenGods[pillar];
          if (!tg || tg.god === 'Day Master') return null;
          const insight = insights?.[pillar];
          return (
            <div key={pillar} className="bg-card-2 border border-hairline rounded-[10px] p-3">
              <div className="flex items-baseline gap-2 mb-1">
                <span className="font-mono text-[10px] text-ink-faint uppercase tracking-wider">
                  {labels[pillar]} · {tg.element}
                </span>
                <span className="font-display text-[15px] font-medium">
                  {tg.god}
                </span>
              </div>
              {insight?.description && (
                <p className="text-ink-dim text-[12px] leading-relaxed">
                  {insight.description}
                </p>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}

function LuckCyclesCard({ cycles }: { cycles: import('../lib/calculations').LuckCycle[] }) {
  return (
    <div className="bg-card border border-hairline rounded-[16px] p-[22px] pb-6 mb-4">
      <div className="font-mono text-[11px] tracking-[0.14em] uppercase text-ink-faint mb-2.5">
        LUCK CYCLES (大运)
      </div>
      <div className="space-y-1.5">
        {cycles.map((cycle, i) => (
          <div key={i} className="flex items-center gap-3 bg-card-2 border border-hairline rounded-[8px] px-3 py-2">
            <span className="font-mono text-[11px] text-ink-faint whitespace-nowrap">
              {cycle.startAge}–{cycle.endAge}
            </span>
            <div className="h-4 w-px bg-hairline" />
            <span className="font-display text-[15px] font-medium">
              {cycle.stem}{cycle.branch}
            </span>
            <span className="font-mono text-[10px] text-ink-dim uppercase">
              {cycle.element}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
