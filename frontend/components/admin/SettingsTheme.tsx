"use client";

import { useTransition } from "react";

import { setThemePreferenceAction } from "@/app/admin/settings/actions";

import { useAdminTheme } from "./AdminThemeProvider";
import { Button } from "./ui/button";

const themeOptions: { value: "light" | "dark" | "system"; label: string; description: string }[] = [
  { value: "light", label: "Mode terang", description: "Latar belakang terang untuk ruang kerja cerah." },
  { value: "dark", label: "Mode gelap", description: "Kontras tinggi yang nyaman untuk sesi malam." },
  {
    value: "system",
    label: "Ikuti sistem",
    description: "Aktifkan penyesuaian otomatis mengikuti preferensi OS.",
  },
];

export function SettingsTheme() {
  const { theme, resolvedTheme, setTheme } = useAdminTheme();
  const [isPending, startTransition] = useTransition();

  const handleSelect = (value: "light" | "dark" | "system") => {
    setTheme(value);
    startTransition(async () => {
      await setThemePreferenceAction(value);
    });
  };

  return (
    <div className="grid gap-3 sm:grid-cols-2">
      {themeOptions.map((option) => {
        const active = option.value === theme || (option.value === resolvedTheme && theme === "system");
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
            <span className="text-xs font-normal text-muted-foreground">{option.description}</span>
          </Button>
        );
      })}
    </div>
  );
}
