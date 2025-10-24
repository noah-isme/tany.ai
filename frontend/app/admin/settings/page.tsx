import { ApiKeysForm, type StoredApiKeySnapshot } from "@/components/admin/ApiKeysForm";
import { SettingsTheme } from "@/components/admin/SettingsTheme";

import { getStoredApiKeys } from "./actions";

export const dynamic = "force-dynamic";

export default async function SettingsPage() {
  const stored = (await getStoredApiKeys()) as StoredApiKeySnapshot | undefined;

  return (
    <div className="space-y-8">
      <section className="space-y-3">
        <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-100">Pengaturan</h1>
        <p className="text-sm text-slate-600 dark:text-slate-400">
          Sesuaikan tampilan panel dan simpan placeholder kredensial integrasi tanpa menampilkan nilai asli.
        </p>
      </section>

      <section className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
        <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Preferensi tema</h2>
        <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">
          Pilih mode gelap atau terang. Perubahan akan diterapkan ke seluruh panel admin.
        </p>
        <div className="mt-4">
          <SettingsTheme />
        </div>
      </section>

      <section className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
        <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">API Keys (placeholder)</h2>
        <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">
          Nilai akan di-hash dan tidak disertakan di bundle klien. Gunakan untuk pengujian alur sebelum integrasi vault.
        </p>
        <div className="mt-4">
          <ApiKeysForm initial={stored} />
        </div>
      </section>
    </div>
  );
}
