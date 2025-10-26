"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { clsx } from "clsx";
import { Menu, X } from "lucide-react";

import { LogoutButton } from "./LogoutButton";
import { ThemeToggle } from "./ThemeToggle";
import { Button } from "./ui/button";
import { ADMIN_NAV_ITEMS, isNavItemActive } from "./AdminSidebar";

const descriptions: Record<string, string> = {
  Dashboard: "Status singkat data knowledge base Anda.",
  Profil: "Perbarui informasi profil dan kontak utama.",
  Skills: "Atur skill yang akan menjadi konteks AI.",
  Layanan: "Kelola layanan dan rentang harga.",
  Proyek: "Dokumentasikan portofolio yang relevan.",
  Statistik: "Pantau performa chat dan lead (dummy).",
  Settings: "Preferensi tampilan dan kredensial integrasi.",
};

type AdminHeaderProps = {
  email: string;
};

export function AdminHeader({ email }: AdminHeaderProps) {
  const pathname = usePathname();
  const activeItem = ADMIN_NAV_ITEMS.find((item) => isNavItemActive(pathname, item.href));
  const [showMobileMenu, setShowMobileMenu] = useState(false);

  return (
    <header className="relative flex items-center justify-between border-b border-border bg-card/95 px-4 py-4 shadow-sm transition-colors duration-300 supports-[backdrop-filter]:bg-card/70 supports-[backdrop-filter]:backdrop-blur">
      <div className="flex items-center gap-3">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="lg:hidden"
          aria-label="Buka navigasi"
          onClick={() => setShowMobileMenu(true)}
        >
          <Menu className="h-5 w-5" />
        </Button>
        <div>
          <h1 className="text-lg font-semibold text-foreground">
            {activeItem?.label ?? "Dashboard"}
          </h1>
          <p className="text-xs text-muted-foreground">
            {descriptions[activeItem?.label ?? "Dashboard"]}
          </p>
        </div>
      </div>

      <div className="flex items-center gap-2">
        <ThemeToggle />
        <div className="hidden items-center gap-2 rounded-full border border-border bg-card/90 px-3 py-1 text-xs text-muted-foreground shadow-sm supports-[backdrop-filter]:bg-card/60 supports-[backdrop-filter]:backdrop-blur sm:flex">
          <span className="font-medium">{email}</span>
        </div>
        <LogoutButton />
      </div>

      {showMobileMenu ? (
        <div className="fixed inset-0 z-50 bg-background/70 backdrop-blur">
          <div className="absolute inset-y-0 left-0 w-72 max-w-[85vw] border-r border-border bg-card p-6 text-card-foreground shadow-xl">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs uppercase tracking-[0.35em] text-primary/70">tany.ai</p>
                <p className="font-semibold">Admin Panel</p>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                aria-label="Tutup navigasi"
                onClick={() => setShowMobileMenu(false)}
              >
                <X className="h-5 w-5" />
              </Button>
            </div>
            <p className="mt-4 text-xs text-muted-foreground">{email}</p>
            <nav className="mt-6">
              <ul className="space-y-2">
                {ADMIN_NAV_ITEMS.map((item) => {
                  const active = isNavItemActive(pathname, item.href);
                  const Icon = item.icon;
                  return (
                    <li key={item.href}>
                      <Link
                        href={item.href}
                        onClick={() => setShowMobileMenu(false)}
                        className={clsx(
                          "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                          active
                            ? "bg-primary text-primary-foreground"
                            : "text-muted-foreground hover:bg-muted hover:text-foreground",
                        )}
                      >
                        <Icon className="h-4 w-4" />
                        <span>{item.label}</span>
                      </Link>
                    </li>
                  );
                })}
              </ul>
            </nav>
          </div>
        </div>
      ) : null}
    </header>
  );
}
