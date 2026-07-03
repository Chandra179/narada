import { useState } from 'react';
import { CityAutocomplete } from './CityAutocomplete';

export type DeepScreenProps = {
  onBack: () => void;
  onSubmit: (data: { birthtime: string; city: string }) => void;
};

export function DeepScreen({ onBack, onSubmit }: DeepScreenProps) {
  const [birthtime, setBirthtime] = useState('');
  const [city, setCity] = useState('');

  function handleSubmit() {
    if (!birthtime || !city) {
      const target = !birthtime ? 'in-birthtime-2' : 'in-city-2';
      document.getElementById(target)?.focus();
      return;
    }
    onSubmit({ birthtime, city });
  }

  return (
    <section className="flex flex-col min-h-[640px]">
      <div className="px-6 pt-[26px] pb-1.5">
        <button
          type="button"
          onClick={onBack}
          className="bg-transparent border-0 text-ink-faint font-mono text-[11px] cursor-pointer flex items-center gap-1.5 mb-[18px]"
        >
          &larr; Back to dashboard
        </button>
        <h2 className="font-display font-medium text-[23px] mb-2">Pinpoint the hour.</h2>
        <p className="text-ink-dim text-[13.5px] leading-relaxed mb-6">
          Your Hour Pillar governs career trajectory and late-life fortune. It needs an exact time and place to calculate.
        </p>
      </div>

      <div className="px-6 pb-[26px] flex flex-col gap-[18px]">
        <div>
          <label className="field-label" htmlFor="in-birthtime-2">
            Birth time
          </label>
          <input
            id="in-birthtime-2"
            type="time"
            value={birthtime}
            onChange={(e) => setBirthtime(e.target.value)}
          />
        </div>
        <div>
          <label className="field-label" htmlFor="in-city-2">
            Birth city
          </label>
          <CityAutocomplete
            id="in-city-2"
            value={city}
            onChange={setCity}
          />
        </div>
        <button
          type="button"
          onClick={handleSubmit}
          className="btn btn-primary mt-2"
        >
          Unlock deep insights
        </button>
      </div>
    </section>
  );
}
