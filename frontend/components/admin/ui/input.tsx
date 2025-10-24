"use client";

import { forwardRef } from "react";
import { clsx } from "clsx";

export type InputProps = React.InputHTMLAttributes<HTMLInputElement> & {
  invalid?: boolean;
};

export const Input = forwardRef<HTMLInputElement, InputProps>(function Input(
  { className, invalid = false, type = "text", ...props },
  ref,
) {
  return (
    <input
      ref={ref}
      type={type}
      className={clsx(
        "w-full rounded-md border bg-white/60 px-3 py-2 text-sm text-slate-900 shadow-sm transition placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400 dark:bg-slate-900/80 dark:text-slate-100 dark:placeholder:text-slate-500",
        invalid
          ? "border-rose-400 focus:border-rose-300 focus:ring-rose-300"
          : "border-slate-300 focus:border-indigo-400 dark:border-slate-700",
        className,
      )}
      {...props}
    />
  );
});
