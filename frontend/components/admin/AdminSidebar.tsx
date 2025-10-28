"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  BarChart3,
  BriefcaseBusiness,
  FolderKanban,
  LayoutDashboard,
  Plug,
  Settings,
  Sparkles,
  UserRound,
} from "lucide-react";
import { clsx } from "clsx";

import { LogoMark } from "@/components/LogoMark";

export const ADMIN_NAV_ITEMS = [
  { href: "/admin", label: "Dashboard", icon: LayoutDashboard },
  { href: "/admin/profile", label: "Profil", icon: UserRound },
  { href: "/admin/skills", label: "Skills", icon: Sparkles },
  { href: "/admin/services", label: "Layanan", icon: BriefcaseBusiness },
  { href: "/admin/projects", label: "Proyek", icon: FolderKanban },
  { href: "/admin/integrations", label: "Integrasi", icon: Plug },
  { href: "/admin/analytics", label: "Analytics", icon: BarChart3 },
  { href: "/admin/settings", label: "Settings", icon: Settings },
];

export function isNavItemActive(pathname: string, href: string) {
  if (href === "/admin") {
    return pathname === "/admin";
  }
  return pathname === href || pathname.startsWith(`${href}/`);
}

type AdminSidebarProps = {
  email: string;
};

export function AdminSidebar({ email }: AdminSidebarProps) {
  const pathname = usePathname();

  return (
    <aside className="hidden w-64 flex-col border-r border-border/70 bg-card/95 px-4 py-6 shadow-sm transition-colors duration-300 supports-[backdrop-filter]:bg-card/80 supports-[backdrop-filter]:backdrop-blur lg:flex">
      <Link
        href="/admin"
        className="flex items-center gap-3 rounded-lg px-2 py-2 text-left text-sm font-semibold text-foreground transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40 hover:bg-muted"
      >
        <LogoMark className="h-10 w-10" />
        <div className="flex flex-col">
          <span className="text-xs uppercase tracking-[0.35em] text-primary/70">tany.ai</span>
          <span>Admin Panel</span>
        </div>
      </Link>
      <div className="mt-6 rounded-lg border border-border bg-card/90 p-3 text-xs text-muted-foreground shadow-sm supports-[backdrop-filter]:bg-card/70 supports-[backdrop-filter]:backdrop-blur">
        <p className="font-semibold text-foreground">Masuk sebagai</p>
        <p className="mt-1 break-all">{email}</p>
      </div>
      <nav className="mt-6 flex-1">
        <ul className="space-y-1">
          {ADMIN_NAV_ITEMS.map((item) => {
            const Icon = item.icon;
            const active = isNavItemActive(pathname, item.href);
            return (
              <li key={item.href}>
                <Link
                  href={item.href}
                  className={clsx(
                    "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40",
                    active
                      ? "bg-primary text-primary-foreground shadow-sm"
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
      <p className="mt-auto text-xs text-muted-foreground">Data di panel ini menjadi knowledge base untuk chat tany.ai.</p>
    </aside>
  );
}
