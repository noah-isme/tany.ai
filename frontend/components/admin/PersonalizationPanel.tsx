"use client";

import { useMemo, useState, useTransition } from "react";

import { useToast } from "@/components/admin/ToastProvider";
import type { ActionResult } from "@/lib/action-result";
import type { PersonalizationSummary } from "@/lib/types/admin";
import { Button } from "./ui/button";

export type PersonalizationPanelProps = {
  summary: PersonalizationSummary;
  updateWeight: (weight: number) => Promise<ActionResult<number>>;
  reindex: () => Promise<ActionResult<{ indexed: number }>>;
  reset: () => Promise<ActionResult<null>>;
};

export function PersonalizationPanel({ summary, updateWeight, reindex, reset }: PersonalizationPanelProps) {
  const toast = useToast();
  const [weight, setWeight] = useState(summary.weight ?? 0);
  const [draftWeight, setDraftWeight] = useState(summary.weight ?? 0);
  const [isPending, startTransition] = useTransition();
  const [isActionPending, startActionTransition] = useTransition();

  const percent = useMemo(() => Math.round((draftWeight ?? 0) * 100), [draftWeight]);
  const lastReindex = summary.lastReindexedAt ? new Date(summary.lastReindexedAt) : undefined;
  const lastReset = summary.lastResetAt ? new Date(summary.lastResetAt) : undefined;
  const disabled = !summary.enabled;

  const handleSaveWeight = () => {
    startTransition(async () => {
      const result = await updateWeight(Number(draftWeight.toFixed(2)));
      if (result.success) {
        setWeight(result.data);
        setDraftWeight(result.data);
        toast({ type: "success", message: result.message ?? "Bobot personalisasi tersimpan." });
      } else {
        setDraftWeight(weight);
        toast({ type: "error", message: result.error ?? "Tidak dapat memperbarui bobot personalisasi." });
      }
    });
  };

  const handleReindex = () => {
    startActionTransition(async () => {
      const result = await reindex();
      if (result.success) {
        toast({
          type: "success",
          message: result.message ?? `Menjadwalkan ${result.data?.indexed ?? 0} embedding untuk diperbarui.`,
        });
      } else {
        toast({ type: "error", message: result.error ?? "Tidak dapat memproses reindex." });
      }
    });
  };

  const handleReset = () => {
    startActionTransition(async () => {
      const result = await reset();
      if (result.success) {
        toast({ type: "success", message: result.message ?? "Semua embedding dihapus." });
      } else {
        toast({ type: "error", message: result.error ?? "Tidak dapat mereset embedding." });
      }
    });
  };

  return (
    <section className="space-y-6 rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
      <div className="grid gap-6 md:grid-cols-2">
        <div className="space-y-4">
          <StatusBadge enabled={summary.enabled} provider={summary.provider} />
          <Stat label="Dimensi embedding" value={summary.dimension.toLocaleString("id-ID")} />
          <Stat label="Total vektor" value={summary.count.toLocaleString("id-ID")}
            helper="Jumlah embedding aktif di Postgres." />
          <Timeline label="Terakhir reindex" timestamp={lastReindex} fallback="Belum pernah" />
          <Timeline label="Terakhir reset" timestamp={lastReset} fallback="Belum pernah" />
        </div>
        <div className="flex flex-col gap-6">
          <div className="rounded-xl border border-slate-200/70 bg-slate-50/80 p-5 dark:border-slate-800 dark:bg-slate-900/60">
            <h3 className="text-sm font-semibold text-slate-800 dark:text-slate-200">Bobot personalisasi</h3>
            <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">
              Atur proporsi instruksi persona yang digabungkan ke prompt AI. Nilai 0 mematikan tone khusus, 1 memaksa gaya penuh.
            </p>
            <div className="mt-4 flex items-center gap-3">
              <input
                type="range"
                min={0}
                max={1}
                step={0.01}
                value={draftWeight}
                onChange={(event) => setDraftWeight(Number(event.target.value))}
                className="flex-1"
                aria-label="Personalization weight"
                disabled={disabled}
              />
              <span className="w-12 text-right text-sm font-semibold text-slate-700 dark:text-slate-200">
                {percent}%
              </span>
            </div>
            <div className="mt-3 flex items-center justify-between text-xs text-slate-500 dark:text-slate-400">
              <span>Nilai tersimpan: {(weight * 100).toFixed(0)}%</span>
              <Button size="sm" onClick={handleSaveWeight} disabled={disabled || isPending || draftWeight === weight}>
                {isPending ? "Menyimpan..." : "Simpan"}
              </Button>
            </div>
            {disabled ? (
              <p className="mt-3 text-xs italic text-slate-500 dark:text-slate-400">
                Personalization dimatikan melalui konfigurasi server. Aktifkan variabel <code>ENABLE_PERSONALIZATION</code> untuk
                mengubah bobot.
              </p>
            ) : null}
          </div>
          <div className="flex flex-wrap gap-3">
            <Button onClick={handleReindex} disabled={isActionPending}>
              {isActionPending ? "Memproses..." : "Reindex embedding"}
            </Button>
            <Button onClick={handleReset} variant="secondary" disabled={isActionPending}>
              Reset embedding
            </Button>
          </div>
        </div>
      </div>
    </section>
  );
}

type StatusBadgeProps = {
  enabled: boolean;
  provider: string;
};

function StatusBadge({ enabled, provider }: StatusBadgeProps) {
  const label = enabled ? "Aktif" : "Nonaktif";
  const color = enabled ? "bg-emerald-500/15 text-emerald-600 dark:bg-emerald-500/20 dark:text-emerald-300" : "bg-slate-200/60 text-slate-600 dark:bg-slate-800/60 dark:text-slate-300";

  return (
    <div className="flex items-center gap-3">
      <span className={`rounded-full px-3 py-1 text-xs font-medium ${color}`}>{label}</span>
      <span className="text-xs text-slate-500 dark:text-slate-400">Provider: {provider || "-"}</span>
    </div>
  );
}

type StatProps = {
  label: string;
  value: string;
  helper?: string;
};

function Stat({ label, value, helper }: StatProps) {
  return (
    <div className="space-y-1">
      <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">{label}</p>
      <p className="text-lg font-semibold text-slate-800 dark:text-slate-100">{value}</p>
      {helper ? <p className="text-xs text-slate-500 dark:text-slate-400">{helper}</p> : null}
    </div>
  );
}

type TimelineProps = {
  label: string;
  timestamp?: Date;
  fallback: string;
};

function Timeline({ label, timestamp, fallback }: TimelineProps) {
  return (
    <div>
      <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">{label}</p>
      <p className="text-sm text-slate-700 dark:text-slate-300">
        {timestamp ? timestamp.toLocaleString("id-ID") : fallback}
      </p>
    </div>
  );
}
