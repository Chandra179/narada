import { useEffect, useState } from 'react';
import { ELEMENTS, ELEMENT_META } from '../lib/constants';
import type { ElementalBalance } from '../lib/calculations';

export type BarChartProps = {
  balance: ElementalBalance;
  blurred?: boolean;
  animate?: boolean;
};

export function BarChart({ balance, blurred = false, animate = true }: BarChartProps) {
  const [animated, setAnimated] = useState(!animate);

  useEffect(() => {
    if (!animate) return;
    const id = requestAnimationFrame(() => setAnimated(true));
    return () => cancelAnimationFrame(id);
  }, [balance, animate]);

  const sorted = [...ELEMENTS].sort((a, b) => balance[b] - balance[a]);

  return (
    <div className={`flex flex-col gap-3 ${blurred ? 'blur-[5px] opacity-60 pointer-events-none' : ''}`}>
      {sorted.map((el) => {
        const meta = ELEMENT_META[el];
        return (
          <div
            key={el}
            className="group grid grid-cols-[58px_1fr_38px] items-center gap-2.5 cursor-default"
          >
            <div className="font-mono text-[11px] tracking-wide text-ink-dim transition-colors duration-150 group-hover:text-ink group-hover:font-bold">
              {meta.short}
            </div>
            <div className="relative h-3.5 bg-hairline rounded-full overflow-hidden">
              <div
                className="absolute left-0 top-0 bottom-0 rounded-full transition-all duration-[900ms] ease-[cubic-bezier(0.2,0.8,0.2,1)] group-hover:brightness-[1.12] group-hover:shadow-[0_0_0_2px_rgba(34,31,25,0.08)]"
                style={{
                  backgroundColor: meta.color,
                  width: animated ? `${balance[el]}%` : '0%',
                }}
              />
            </div>
            <div className="font-mono text-[11.5px] text-ink-faint text-right transition-colors duration-150 group-hover:text-ink">
              {balance[el]}%
            </div>
          </div>
        );
      })}
    </div>
  );
}
