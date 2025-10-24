import type { Metadata } from "next";

import { AdminHeader } from "@/components/admin/AdminHeader";
import { AdminSidebar } from "@/components/admin/AdminSidebar";
import { AdminThemeProvider } from "@/components/admin/AdminThemeProvider";
import { ToastProvider } from "@/components/admin/ToastProvider";
import { getStoredTheme, requireAdminOrRedirect } from "@/lib/auth";

export const metadata: Metadata = {
  title: "Admin Panel · tany.ai",
  description: "Kelola knowledge base tany.ai melalui panel admin.",
};

type AdminLayoutProps = {
  children: React.ReactNode;
};

export default async function AdminLayout({ children }: AdminLayoutProps) {
  const user = await requireAdminOrRedirect();
  const storedTheme = await getStoredTheme();

  return (
    <AdminThemeProvider defaultTheme={storedTheme ?? "dark"}>
      <ToastProvider>
        <a
          href="#admin-main"
          className="sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-4 focus:z-50 focus:rounded-md focus:bg-indigo-600 focus:px-4 focus:py-2 focus:text-sm focus:text-white"
        >
          Lompat ke konten utama
        </a>
        <div className="flex min-h-screen bg-slate-100 text-slate-900 transition dark:bg-slate-950 dark:text-slate-100">
          <AdminSidebar email={user.email} />
          <div className="flex flex-1 flex-col">
            <AdminHeader email={user.email} />
            <main id="admin-main" className="flex-1 overflow-y-auto px-4 py-6 sm:px-6 lg:px-8">
              {children}
            </main>
          </div>
        </div>
      </ToastProvider>
    </AdminThemeProvider>
  );
}
