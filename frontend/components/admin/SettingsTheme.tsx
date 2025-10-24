"use client";

import { useTransition } from "react";

import { setThemePreferenceAction } from "@/app/admin/settings/actions";

import { useAdminTheme } from "./AdminThemeProvider";
import { Button } from "./ui/button";

const themeOptions: { value: "light" | "dark"; label: string; description: string }[] = [
  { value: "light", label: "Mode terang", description: "Latar belakang terang untuk ruang kerja cerah." },
  { value: "dark", label: "Mode gelap", description: "Kontras tinggi yang nyaman untuk sesi malam." },
];

export function SettingsTheme() {
  const { theme, setTheme } = useAdminTheme();
  const [isPending, startTransition] = useTransition();

  const handleSelect = (value: "light" | "dark") => {
    setTheme(value);
    startTransition(async () => {
      await setThemePreferenceAction(value);
    });
  };

  return (
    <div className="grid gap-3 sm:grid-cols-2">
      {themeOptions.map((option) => {
        const active = option.value === theme;
        return (
          <Button
            key={option.value}
            type="button"
            variant={active ? "secondary" : "ghost"}
            className="flex h-full flex-col items-start gap-1 text-left"
            onClick={() => handleSelect(option.value)}
            disabled={isPending}
          >
            <span className="text-sm font-semibold">{option.label}</span>
            <span className="text-xs font-normal text-slate-500 dark:text-slate-400">{option.description}</span>
          </Button>
        );
      })}
    </div>
  );
}
