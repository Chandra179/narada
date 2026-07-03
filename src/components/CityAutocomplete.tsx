import { useEffect, useRef, useState } from 'react';
import { CITIES, type City } from '../lib/constants';

export type CityAutocompleteProps = {
  value: string;
  onChange: (value: string) => void;
  id?: string;
  placeholder?: string;
};

export function CityAutocomplete({
  value,
  onChange,
  id,
  placeholder = 'Start typing a city…',
}: CityAutocompleteProps) {
  const [open, setOpen] = useState(false);
  const [matches, setMatches] = useState<City[]>([]);
  const wrapperRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener('click', handleClickOutside);
    return () => document.removeEventListener('click', handleClickOutside);
  }, []);

  function handleInput(raw: string) {
    onChange(raw);
    const q = raw.trim().toLowerCase();
    if (q.length < 2) {
      setOpen(false);
      setMatches([]);
      return;
    }
    const next = CITIES.filter((c) => c.name.toLowerCase().includes(q)).slice(0, 6);
    setMatches(next);
    setOpen(true);
  }

  function selectCity(name: string) {
    onChange(name);
    setOpen(false);
  }

  const showEmpty = open && matches.length === 0 && value.trim().length >= 2;

  return (
    <div ref={wrapperRef} className="relative">
      <input
        id={id}
        type="text"
        value={value}
        onChange={(e) => handleInput(e.target.value)}
        placeholder={placeholder}
        autoComplete="off"
      />
      {open && matches.length > 0 && (
        <div className="absolute top-[calc(100%+6px)] left-0 right-0 bg-card-2 border border-hairline-strong rounded-[10px] overflow-hidden z-20 max-h-[220px] overflow-y-auto shadow-[0_12px_30px_-10px_rgba(60,50,30,0.25)]">
          {matches.map((city) => (
            <button
              key={city.name}
              type="button"
              onClick={() => selectCity(city.name)}
              className="w-full text-left px-4 py-3 text-sm border-b border-hairline last:border-b-0 flex justify-between items-center hover:bg-bg-alabaster focus:bg-bg-alabaster focus:outline-none"
            >
              <span>{city.name}</span>
              <span className="font-mono text-[11px] text-ink-faint">{city.country}</span>
            </button>
          ))}
        </div>
      )}
      {showEmpty && (
        <div className="mt-[6px] px-4 py-3.5 text-xs text-ink-dim leading-relaxed bg-card-2 border border-hairline-strong rounded-[10px]">
          Can't find your town? Try typing the nearest major city or regional capital instead. This keeps your timezone calculation accurate.
        </div>
      )}
    </div>
  );
}
