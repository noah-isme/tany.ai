import type { Metadata } from "next";
import Link from "next/link";

import { LogoMark } from "@/components/LogoMark";

import { LoginForm } from "./LoginForm";

export const metadata: Metadata = {
  title: "Masuk Admin · tany.ai",
  description: "Panel login admin tany.ai untuk mengelola knowledge base.",
};

export default function LoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-950 px-6 py-16 text-slate-100">
      <div className="w-full max-w-md space-y-8">
        <div className="flex flex-col items-center gap-4 text-center">
          <LogoMark className="h-10 w-10 text-indigo-300" />
          <div>
            <h1 className="text-xl font-semibold text-white">Masuk ke Admin Panel</h1>
            <p className="mt-2 text-sm text-slate-300">
              Kelola profil, layanan, proyek, dan knowledge base tany.ai dari satu tempat.
            </p>
          </div>
        </div>

        <div className="rounded-3xl border border-white/10 bg-slate-900/60 p-8 shadow-xl shadow-black/40">
          <LoginForm />
        </div>

        <p className="text-center text-xs text-slate-400">
          Kembali ke {" "}
          <Link href="/" className="font-medium text-indigo-200 underline">
            halaman utama tany.ai
          </Link>
        </p>
      </div>
    </div>
  );
}
