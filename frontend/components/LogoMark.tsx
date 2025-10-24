import { clsx } from "clsx";

type LogoMarkProps = {
  className?: string;
};

export function LogoMark({ className }: LogoMarkProps) {
  return (
    <div
      aria-hidden="true"
      className={clsx(
        "flex h-12 w-12 items-center justify-center rounded-2xl border border-white/20 bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 shadow-lg shadow-indigo-900/40",
        className,
      )}
    >
      <span className="text-xl font-semibold text-white">ta</span>
    </div>
  );
}
