import { useState } from 'react';
import { DeepScreen } from './components/DeepScreen';
import { DashboardScreen } from './components/DashboardScreen';
import { Modal } from './components/Modal';
import { OnboardScreen } from './components/OnboardScreen';
import { usePersistentState } from './hooks/usePersistentState';
import type { TabKey } from './lib/constants';

export default function App() {
  const { state, updateState, resetState, consumePendingAddTime } = usePersistentState();
  const [showModal, setShowModal] = useState(false);

  function handleReveal(data: {
    birthdate: string;
    birthtime: string;
    city: string;
    gender: string;
    unlocked: boolean;
  }) {
    console.log('[app] handleReveal:', data);
    if (!data.birthtime || !data.city) {
      updateState({
        birthdate: data.birthdate,
        birthtime: data.birthtime,
        city: data.city,
        gender: data.gender,
        unlocked: false,
        pendingAddTime: false,
      });
      setShowModal(true);
      return;
    }

    updateState({
      birthdate: data.birthdate,
      birthtime: data.birthtime,
      city: data.city,
      gender: data.gender,
      unlocked: true,
      screen: 'dashboard',
      activeTab: 'romance',
      pendingAddTime: false,
    });
  }

  function handleProceedBasic() {
    setShowModal(false);
    updateState({
      unlocked: false,
      screen: 'dashboard',
      activeTab: 'romance',
    });
  }

  function handleAddTime() {
    setShowModal(false);
    updateState({ screen: 'onboard', pendingAddTime: true });
  }

  function handleDeepSubmit(data: { birthtime: string; city: string }) {
    updateState({
      birthtime: data.birthtime,
      city: data.city,
      unlocked: true,
      screen: 'dashboard',
    });
  }

  function handleTabChange(tab: TabKey) {
    updateState({ activeTab: tab });
  }

  function handleUnlock() {
    updateState({ screen: 'deep' });
  }

  function handleReset() {
    resetState();
  }

  return (
    <div className="min-h-screen flex justify-center items-start px-3.5 py-7">
      <div className="w-full max-w-[440px] bg-bg-alabaster border border-hairline-strong rounded-[28px] shadow-[0_30px_60px_-30px_rgba(60,50,30,0.28),0_0_0_1px_rgba(255,255,255,0.4)_inset] overflow-hidden relative min-h-[640px]">
        {state.screen === 'onboard' && (
          <OnboardScreen
            initialBirthdate={state.birthdate}
            initialBirthtime={state.birthtime}
            initialCity={state.city}
            initialGender={state.gender}
            pendingAddTime={state.pendingAddTime}
            onReveal={handleReveal}
            onAddTimeAcknowledged={consumePendingAddTime}
          />
        )}
        {state.screen === 'deep' && (
          <DeepScreen
            onBack={() => updateState({ screen: 'dashboard' })}
            onSubmit={handleDeepSubmit}
          />
        )}
        {state.screen === 'dashboard' && state.birthdate && (
          <DashboardScreen
            birthdate={state.birthdate}
            birthtime={state.birthtime}
            gender={state.gender}
            unlocked={state.unlocked}
            activeTab={state.activeTab}
            onTabChange={handleTabChange}
            onUnlock={handleUnlock}
            onReset={handleReset}
          />
        )}
        <Modal
          open={showModal}
          onProceedBasic={handleProceedBasic}
          onAddTime={handleAddTime}
        />
      </div>
    </div>
  );
}
