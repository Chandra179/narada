import { useEffect, useState } from 'react';

const STORAGE_KEY = 'wu-xing-state-v1';

export type AppState = {
  birthdate: string;
  birthtime: string;
  city: string;
  unlocked: boolean;
  screen: 'onboard' | 'deep' | 'dashboard';
  activeTab: 'romance' | 'health' | 'career' | 'wealth';
  pendingAddTime?: boolean;
};

const DEFAULT_STATE: AppState = {
  birthdate: '',
  birthtime: '',
  city: '',
  unlocked: false,
  screen: 'onboard',
  activeTab: 'romance',
  pendingAddTime: false,
};

function loadState(): AppState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return DEFAULT_STATE;
    const parsed = JSON.parse(raw) as Partial<AppState>;
    return { ...DEFAULT_STATE, ...parsed };
  } catch {
    return DEFAULT_STATE;
  }
}

function saveState(state: AppState) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  } catch {
    // ignore storage errors
  }
}

export function usePersistentState() {
  const [state, setState] = useState<AppState>(loadState);

  useEffect(() => {
    saveState(state);
  }, [state]);

  function updateState(partial: Partial<AppState>) {
    setState((prev) => ({ ...prev, ...partial }));
  }

  function resetState() {
    setState(DEFAULT_STATE);
  }

  function consumePendingAddTime() {
    setState((prev) => ({ ...prev, pendingAddTime: false }));
  }

  return { state, updateState, resetState, consumePendingAddTime };
}
