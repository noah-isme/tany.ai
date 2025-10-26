"use client";

import { Moon, Sun } from "lucide-react";

import { Button } from "./ui/button";
import { useAdminTheme } from "./AdminThemeProvider";

export function ThemeToggle() {
  const { resolvedTheme, toggleTheme } = useAdminTheme();
  const isDark = resolvedTheme === "dark";
  return (
    <Button
      type="button"
      variant="ghost"
      size="sm"
      aria-label={isDark ? "Aktifkan mode terang" : "Aktifkan mode gelap"}
      onClick={toggleTheme}
    >
      {isDark ? <Moon className="h-4 w-4" /> : <Sun className="h-4 w-4" />}
      <span className="sr-only">Toggle theme</span>
    </Button>
  );
}
