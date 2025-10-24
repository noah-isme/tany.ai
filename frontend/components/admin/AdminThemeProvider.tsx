"use client";

import { createContext, useContext, useEffect, useMemo, useState } from "react";

export type ThemeMode = "light" | "dark";

type ThemeContextValue = {
  theme: ThemeMode;
  setTheme: (theme: ThemeMode) => void;
  toggleTheme: () => void;
};

const STORAGE_KEY = "tany-admin-theme";

const AdminThemeContext = createContext<ThemeContextValue | undefined>(undefined);

function applyTheme(theme: ThemeMode) {
  const root = document.documentElement;
  if (theme === "dark") {
    root.classList.add("dark");
  } else {
    root.classList.remove("dark");
  }
  root.dataset.theme = theme;
}

function persist(theme: ThemeMode) {
  if (typeof window !== "undefined") {
    window.localStorage.setItem(STORAGE_KEY, theme);
    document.cookie = `ta_theme=${theme}; path=/; max-age=${60 * 60 * 24 * 365}; samesite=lax`;
  }
}

type ProviderProps = {
  children: React.ReactNode;
  defaultTheme?: ThemeMode | null;
};

export function AdminThemeProvider({ children, defaultTheme = "dark" }: ProviderProps) {
  const [theme, setThemeState] = useState<ThemeMode>(defaultTheme ?? "dark");

  useEffect(() => {
    if (typeof document === "undefined") {
      return;
    }
    applyTheme(theme);
    persist(theme);
  }, [theme]);

  const value = useMemo<ThemeContextValue>(
    () => ({
      theme,
      setTheme: (next) => setThemeState(next),
      toggleTheme: () => setThemeState((prev) => (prev === "dark" ? "light" : "dark")),
    }),
    [theme],
  );

  return <AdminThemeContext.Provider value={value}>{children}</AdminThemeContext.Provider>;
}

export function useAdminTheme(): ThemeContextValue {
  const context = useContext(AdminThemeContext);
  if (!context) {
    throw new Error("useAdminTheme must be used within AdminThemeProvider");
  }
  return context;
}
