import { useEffect, useState } from 'react';

export type TabConfig = {
  key: string;
  label: string;
  icon: string;
};

type ElementMeta = {
  label: string;
  short: string;
  color: string;
};

type City = {
  name: string;
  country: string;
  timezone: number;
};

export type AppConfig = {
  tabs: TabConfig[];
  elements: Record<string, ElementMeta>;
  cities: City[];
  lifePathNumbers: Record<number, string>;
};

const CACHE_KEY = 'wu-xing-config-v1';

function isAppConfig(obj: unknown): obj is AppConfig {
  if (!obj || typeof obj !== 'object') return false;
  const c = obj as Record<string, unknown>;
  return (
    Array.isArray(c.tabs) &&
    typeof c.elements === 'object' &&
    c.elements !== null &&
    Array.isArray(c.cities) &&
    typeof c.lifePathNumbers === 'object'
  );
}

function loadCached(): AppConfig | null {
  try {
    const raw = localStorage.getItem(CACHE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    return isAppConfig(parsed) ? parsed : null;
  } catch {
    return null;
  }
}

function cacheConfig(config: AppConfig) {
  try {
    localStorage.setItem(CACHE_KEY, JSON.stringify(config));
  } catch {
    // ignore storage errors
  }
}

export function useConfig() {
  const [config, setConfig] = useState<AppConfig | null>(loadCached);
  const [loading, setLoading] = useState(!config);
  const [error, setError] = useState(false);

  useEffect(() => {
    let cancelled = false;
    fetch('/api/config')
      .then((res) => {
        if (!res.ok) throw new Error(`config error: ${res.status}`);
        return res.json();
      })
      .then((data: AppConfig) => {
        if (cancelled) return;
        if (isAppConfig(data)) {
          cacheConfig(data);
          setConfig(data);
        }
      })
      .catch(() => {
        if (!cancelled) setError(true);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, []);

  return { config, loading, error };
}
