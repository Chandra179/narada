import { ClockIcon } from './icons';

export type ModalProps = {
  open: boolean;
  onProceedBasic: () => void;
  onAddTime: () => void;
};

export function Modal({ open, onProceedBasic, onAddTime }: ModalProps) {
  if (!open) return null;

  return (
    <div className="absolute inset-0 bg-[rgba(34,31,25,0.45)] backdrop-blur-[3px] flex items-end justify-center z-40 p-4 animate-fade-in">
      <div className="bg-card-2 border border-hairline-strong rounded-[24px] p-7 pb-6 w-full max-w-[392px] shadow-[0_20px_50px_-20px_rgba(60,50,30,0.35)] animate-rise-in">
        <ClockIcon className="w-[34px] h-[34px] mb-4 text-earth" />
        <h3 className="font-display font-medium text-xl leading-tight mb-2.5">Missing your birth time</h3>
        <p className="text-ink-dim text-sm leading-relaxed mb-6">
          Without it, we can't calculate your <strong className="text-ink font-semibold">Hour Pillar</strong> — the piece that governs career and late-life fortune. You'll still get your core personality reading, but career and wealth stay locked.
        </p>
        <div className="flex flex-col gap-2.5">
          <button
            type="button"
            onClick={onProceedBasic}
            className="btn btn-ghost"
          >
            Proceed with basic profile
          </button>
          <button
            type="button"
            onClick={onAddTime}
            className="btn btn-primary"
          >
            Add time now
          </button>
        </div>
      </div>
    </div>
  );
}
