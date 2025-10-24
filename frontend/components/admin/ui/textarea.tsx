"use client";

import { forwardRef } from "react";
import { clsx } from "clsx";

export type TextareaProps = React.TextareaHTMLAttributes<HTMLTextAreaElement> & {
  invalid?: boolean;
};

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(function Textarea(
  { className, invalid = false, ...props },
  ref,
) {
  return (
    <textarea
      ref={ref}
      className={clsx(
        "min-h-[120px] w-full rounded-md border bg-white/60 px-3 py-2 text-sm text-slate-900 shadow-sm transition placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400 dark:bg-slate-900/80 dark:text-slate-100",
        invalid
          ? "border-rose-400 focus:border-rose-300 focus:ring-rose-300"
          : "border-slate-300 focus:border-indigo-400 dark:border-slate-700",
        className,
      )}
      {...props}
    />
  );
});
