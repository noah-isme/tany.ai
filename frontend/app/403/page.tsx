import type { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
  title: "Akses ditolak · tany.ai",
  description: "Anda tidak memiliki hak akses untuk membuka halaman ini.",
};

export default function ForbiddenPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-slate-950 px-6 py-16 text-slate-100">
      <div className="w-full max-w-md space-y-6 text-center">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-rose-300">403</p>
        <h1 className="text-2xl font-semibold text-white">Akses ditolak</h1>
        <p className="text-sm text-slate-300">
          Akun Anda tidak memiliki hak akses admin. Silakan hubungi pengelola sistem untuk mendapatkan izin atau gunakan akun lain.
        </p>
        <div className="flex flex-col gap-3 sm:flex-row sm:justify-center">
          <Link
            href="/"
            className="inline-flex items-center justify-center rounded-md border border-white/10 px-4 py-2 text-sm font-medium text-white transition hover:border-white/40"
          >
            Kembali ke beranda
          </Link>
          <Link
            href="/login"
            className="inline-flex items-center justify-center rounded-md bg-indigo-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-indigo-400"
          >
            Masuk sebagai admin
          </Link>
        </div>
      </div>
    </div>
  );
}
