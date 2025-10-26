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
        "w-full rounded-md border border-border bg-card/80 px-3 py-2 text-sm text-foreground shadow-sm transition-colors duration-200 placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring supports-[backdrop-filter]:bg-card/60 supports-[backdrop-filter]:backdrop-blur",
        invalid
          ? "border-destructive/70 focus:border-destructive focus:ring-destructive/40"
          : "focus:border-primary/40",
        className,
      )}
      {...props}
    />
  );
});
