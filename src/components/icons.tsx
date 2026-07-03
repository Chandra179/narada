import type { SVGProps } from 'react';

export function MarkIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 46 46"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <circle cx="23" cy="23" r="21" stroke="#5E8A4C" strokeWidth="1.1" opacity="0.75" />
      <circle cx="23" cy="23" r="15" stroke="#BD5B3B" strokeWidth="1.1" opacity="0.75" />
      <circle cx="23" cy="23" r="9" stroke="#35708E" strokeWidth="1.1" opacity="0.75" />
      <circle cx="23" cy="23" r="2.4" fill="#A97C34" />
    </svg>
  );
}

export function RomanceIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" {...props}>
      <path d="M12 20s-7.5-4.6-10-9.3C.5 7 2.3 3.5 6 3c2.3-.3 4.2 1 6 3.2C13.8 4 15.7 2.7 18 3c3.7.5 5.5 4 4 7.7C19.5 15.4 12 20 12 20z" />
    </svg>
  );
}

export function HealthIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" {...props}>
      <path d="M12 3c-2 3-3 5.5-3 8a3 3 0 0 0 6 0c0-2.5-1-5-3-8z" />
      <path d="M6.5 21c0-4 2.5-6.5 5.5-6.5s5.5 2.5 5.5 6.5" />
    </svg>
  );
}

export function CareerIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" {...props}>
      <rect x="3" y="7" width="18" height="13" rx="1.5" />
      <path d="M8 7V5a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
    </svg>
  );
}

export function WealthIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" {...props}>
      <circle cx="9" cy="9" r="5.5" />
      <circle cx="15" cy="15" r="5.5" />
    </svg>
  );
}

export function ClockIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.4" {...props}>
      <circle cx="12" cy="12" r="9" />
      <path d="M12 7.5v5l3.2 1.9" />
    </svg>
  );
}

export function LockIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" {...props}>
      <rect x="5" y="10" width="14" height="10" rx="2" />
      <path d="M8 10V7a4 4 0 0 1 8 0v3" />
    </svg>
  );
}

export function StarIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" {...props}>
      <path d="M12 3l2.2 5 5.8.5-4.4 3.9L17 18l-5-3-5 3 1.4-5.6L4 8.5l5.8-.5z" />
    </svg>
  );
}
