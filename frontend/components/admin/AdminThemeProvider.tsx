"use client";

import { createContext, useContext, useEffect, useMemo } from "react";
import { useTheme } from "next-themes";

export type ThemeMode = "light" | "dark" | "system";

type ThemeContextValue = {
  theme: ThemeMode;
  resolvedTheme: "light" | "dark";
  setTheme: (theme: ThemeMode) => void;
  toggleTheme: () => void;
};

const AdminThemeContext = createContext<ThemeContextValue | undefined>(undefined);

type ProviderProps = {
  children: React.ReactNode;
};

export function AdminThemeProvider({ children }: ProviderProps) {
  const { theme = "system", resolvedTheme = "light", setTheme } = useTheme();

  const value = useMemo<ThemeContextValue>(() => {
    const activeTheme = (resolvedTheme ?? "light") as "light" | "dark";
    const currentTheme = (theme ?? "system") as ThemeMode;
    return {
      theme: currentTheme,
      resolvedTheme: activeTheme,
      setTheme: (next) => setTheme(next),
      toggleTheme: () => setTheme(activeTheme === "dark" ? "light" : "dark"),
    };
  }, [resolvedTheme, setTheme, theme]);

  useEffect(() => {
    const selected = value.theme;
    const cookieValue = selected;
    const maxAge = 60 * 60 * 24 * 365;
    const secure = typeof window !== "undefined" && window.location.protocol === "https:";
    document.cookie = `ta_theme=${cookieValue}; path=/; max-age=${maxAge}; samesite=lax${secure ? "; secure" : ""}`;
  }, [value.theme]);

  return <AdminThemeContext.Provider value={value}>{children}</AdminThemeContext.Provider>;
}

export function useAdminTheme(): ThemeContextValue {
  const context = useContext(AdminThemeContext);
  if (!context) {
    throw new Error("useAdminTheme must be used within AdminThemeProvider");
  }
  return context;
}
