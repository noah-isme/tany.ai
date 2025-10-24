"use client";

import { forwardRef } from "react";
import { clsx } from "clsx";

export type SwitchProps = Omit<React.InputHTMLAttributes<HTMLInputElement>, "type">;

export const Switch = forwardRef<HTMLInputElement, SwitchProps>(function Switch(
  { className, disabled, ...props },
  ref,
) {
  return (
    <label
      className={clsx(
        "inline-flex cursor-pointer items-center gap-2",
        disabled ? "cursor-not-allowed opacity-60" : "",
        className,
      )}
    >
      <input
        ref={ref}
        type="checkbox"
        className="peer sr-only"
        disabled={disabled}
        {...props}
      />
      <span
        aria-hidden="true"
        className="pointer-events-none relative inline-flex h-5 w-9 items-center rounded-full border border-slate-300 bg-slate-300 transition peer-checked:border-indigo-500 peer-checked:bg-indigo-500 dark:border-slate-700 dark:bg-slate-700"
      >
        <span
          className="pointer-events-none ml-1 inline-block h-3 w-3 rounded-full bg-white shadow-sm transition peer-checked:translate-x-4"
        />
      </span>
    </label>
  );
});
