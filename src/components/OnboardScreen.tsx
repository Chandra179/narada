import { useEffect, useState } from 'react';
import { CityAutocomplete } from './CityAutocomplete';
import { MarkIcon } from './icons';

export type OnboardScreenProps = {
  initialBirthdate?: string;
  initialBirthtime?: string;
  initialCity?: string;
  pendingAddTime?: boolean;
  onReveal: (data: { birthdate: string; birthtime: string; city: string; unlocked: boolean }) => void;
  onAddTimeAcknowledged?: () => void;
};

export function OnboardScreen({
  initialBirthdate = '',
  initialBirthtime = '',
  initialCity = '',
  pendingAddTime = false,
  onReveal,
  onAddTimeAcknowledged,
}: OnboardScreenProps) {
  const [birthdate, setBirthdate] = useState(initialBirthdate);
  const [birthtime, setBirthtime] = useState(initialBirthtime);
  const [city, setCity] = useState(initialCity);
  const [optionalOpen, setOptionalOpen] = useState(Boolean(initialBirthtime || initialCity));

  useEffect(() => {
    if (pendingAddTime) {
      setOptionalOpen(true);
      onAddTimeAcknowledged?.();
      // Defer focus until the optional fields have expanded.
      const timer = setTimeout(() => {
        document.getElementById('in-birthtime')?.focus();
      }, 360);
      return () => clearTimeout(timer);
    }
  }, [pendingAddTime, onAddTimeAcknowledged]);

  function handleReveal() {
    if (!birthdate) {
      document.getElementById('in-birthdate')?.focus();
      return;
    }

    const unlocked = Boolean(birthtime && city);
    onReveal({ birthdate, birthtime, city, unlocked });
  }

  return (
    <section className="flex flex-col min-h-[640px] px-[26px] pt-11 pb-[30px] justify-between">
      <div>
        <MarkIcon className="w-[46px] h-[46px] mx-auto mb-[22px]" />
        <div className="text-center mb-[34px]">
          <h1 className="font-display font-medium text-[28px] leading-[1.15] mb-2.5 tracking-tight">
            What are you made of?
          </h1>
          <p className="text-ink-dim text-[14.5px] leading-[1.55] max-w-[300px] mx-auto">
            One date reveals your elemental core. Add your birth time later to unlock how it moves through career and wealth.
          </p>
        </div>

        <div className="mb-[18px]">
          <label className="field-label" htmlFor="in-birthdate">
            Date of birth
          </label>
          <input
            id="in-birthdate"
            type="date"
            value={birthdate}
            onChange={(e) => setBirthdate(e.target.value)}
          />
        </div>

        <div className="flex items-center gap-2 mx-0.5 mb-5.5">
          <button
            type="button"
            onClick={() => setOptionalOpen((o) => !o)}
            className="btn-link"
          >
            {optionalOpen ? '– Hide extra details' : '+ Add birth time & city for deeper insight'}
          </button>
        </div>

        <div
          className={`overflow-hidden transition-[max-height] duration-[350ms] ease-in-out ${
            optionalOpen ? 'max-h-[400px]' : 'max-h-0'
          }`}
        >
          <div className="pt-4 flex flex-col gap-3.5">
            <div>
              <label className="field-label" htmlFor="in-birthtime">
                Birth time
              </label>
              <input
                id="in-birthtime"
                type="time"
                value={birthtime}
                onChange={(e) => setBirthtime(e.target.value)}
              />
            </div>
            <div>
              <label className="field-label" htmlFor="in-city">
                Birth city
              </label>
              <CityAutocomplete
                id="in-city"
                value={city}
                onChange={setCity}
              />
            </div>
          </div>
        </div>
      </div>

      <div className="mt-auto pt-6">
        <button type="button" onClick={handleReveal} className="btn btn-primary">
          Reveal my profile
        </button>
        <p className="text-center text-xs text-ink-faint mt-3.5 leading-relaxed">
          Your data stays on this device. Nothing is shared.
        </p>
      </div>
    </section>
  );
}
