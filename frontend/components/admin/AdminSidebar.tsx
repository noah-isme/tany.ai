"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  BarChart3,
  BriefcaseBusiness,
  FolderKanban,
  LayoutDashboard,
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
  { href: "/admin/stats", label: "Statistik", icon: BarChart3 },
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
    <aside className="hidden w-64 flex-col border-r border-slate-200 bg-white/70 px-4 py-6 dark:border-slate-800 dark:bg-slate-950/70 lg:flex">
      <Link href="/admin" className="flex items-center gap-3 rounded-lg px-2 py-2 text-left text-sm font-semibold text-slate-900 transition hover:bg-slate-100 dark:text-slate-100 dark:hover:bg-slate-800">
        <LogoMark className="h-10 w-10" />
        <div className="flex flex-col">
          <span className="text-xs uppercase tracking-[0.35em] text-indigo-400">tany.ai</span>
          <span>Admin Panel</span>
        </div>
      </Link>
      <div className="mt-6 rounded-lg border border-slate-200 bg-white/70 p-3 text-xs text-slate-500 dark:border-slate-800 dark:bg-slate-900/70 dark:text-slate-400">
        <p className="font-semibold text-slate-700 dark:text-slate-100">Masuk sebagai</p>
        <p className="mt-1 break-all text-slate-600 dark:text-slate-300">{email}</p>
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
                    "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-400",
                    active
                      ? "bg-indigo-500/90 text-white shadow"
                      : "text-slate-600 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-300 dark:hover:bg-slate-800",
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
      <p className="mt-auto text-xs text-slate-400">Data di panel ini menjadi knowledge base untuk chat tany.ai.</p>
    </aside>
  );
}
