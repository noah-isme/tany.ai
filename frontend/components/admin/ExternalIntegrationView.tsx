"use client";

import { useMemo, useState, useTransition } from "react";
import { useRouter } from "next/navigation";
import { clsx } from "clsx";

import type { ActionResult } from "@/lib/action-result";
import type { ExternalItem, ExternalSource } from "@/lib/types/admin";

import { Button } from "./ui/button";
import { Switch } from "./ui/switch";

type Status = { type: "success" | "error"; message: string } | null;

type ExternalIntegrationViewProps = {
  initialSources: ExternalSource[];
  initialItems: ExternalItem[];
  syncSource: (id: string) => Promise<ActionResult<{ itemsUpserted: number; message?: string }>>;
  toggleItem: (params: { id: string; visible: boolean }) => Promise<ActionResult<ExternalItem>>;
};

const dateFormatter = new Intl.DateTimeFormat("id-ID", {
  dateStyle: "medium",
  timeStyle: "short",
});

export function ExternalIntegrationView({
  initialSources,
  initialItems,
  syncSource,
  toggleItem,
}: ExternalIntegrationViewProps) {
  const router = useRouter();
  const [sources, setSources] = useState(initialSources);
  const [items, setItems] = useState(initialItems);
  const [status, setStatus] = useState<Status>(null);
  const [pendingId, setPendingId] = useState<string | null>(null);
  const [isPending, startTransition] = useTransition();

  const visibleCounts = useMemo(() => {
    const total = items.length;
    const visible = items.filter((item) => item.visible).length;
    return { total, visible };
  }, [items]);

  const handleSync = (source: ExternalSource) => {
    setPendingId(source.id);
    startTransition(async () => {
      const result = await syncSource(source.id);
      if (result.success) {
        setStatus({
          type: "success",
          message:
            result.data.message ??
            `Sinkronisasi selesai. ${result.data.itemsUpserted} item diperbarui.`,
        });
        setSources((prev) =>
          prev.map((item) =>
            item.id === source.id
              ? {
                  ...item,
                  lastSyncedAt: new Date().toISOString(),
                  lastModified: item.lastModified,
                }
              : item,
          ),
        );
        router.refresh();
      } else {
        setStatus({ type: "error", message: result.error });
      }
      setPendingId(null);
    });
  };

  const handleToggleVisibility = (item: ExternalItem, visible: boolean) => {
    setPendingId(item.id);
    startTransition(async () => {
      const result = await toggleItem({ id: item.id, visible });
      if (result.success) {
        setItems((prev) =>
          prev.map((entry) => (entry.id === item.id ? { ...entry, visible: result.data.visible } : entry)),
        );
      } else {
        setStatus({ type: "error", message: result.error });
      }
      setPendingId(null);
    });
  };

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h1 className="text-xl font-semibold text-foreground">Integrasi Konten Eksternal</h1>
        <p className="text-sm text-muted-foreground">
          Sinkronkan proyek dan artikel dari noahis.me agar muncul di jawaban AI secara otomatis.
        </p>
      </div>

      {status ? (
        <p
          className={clsx(
            "rounded-lg border px-4 py-3 text-sm",
            status.type === "success"
              ? "border-emerald-400/70 bg-emerald-500/10 text-emerald-300"
              : "border-rose-400/70 bg-rose-500/10 text-rose-200",
          )}
        >
          {status.message}
        </p>
      ) : null}

      <section className="space-y-4">
        <header className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 className="text-base font-semibold text-foreground">Sumber Terhubung</h2>
            <p className="text-xs text-muted-foreground">
              Daftar sumber yang disinkronkan secara otomatis sesuai jadwal.
            </p>
          </div>
          <span className="text-xs text-muted-foreground">
            {visibleCounts.visible} konten aktif dari {visibleCounts.total} total item.
          </span>
        </header>

        <div className="overflow-hidden rounded-2xl border border-border bg-card/95 shadow-sm supports-[backdrop-filter]:bg-card/80 supports-[backdrop-filter]:backdrop-blur">
          <table className="min-w-full divide-y divide-border/70 text-sm">
            <thead className="bg-muted/80">
              <tr>
                <th className="px-4 py-3 text-left font-semibold text-muted-foreground">Sumber</th>
                <th className="px-4 py-3 text-left font-semibold text-muted-foreground">Terakhir Sinkron</th>
                <th className="w-40 px-4 py-3 text-right font-semibold text-muted-foreground">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/70">
              {sources.map((source) => (
                <tr key={source.id} className="bg-card/80">
                  <td className="px-4 py-3 align-top">
                    <div className="font-medium text-foreground">{source.name}</div>
                    <div className="text-xs text-muted-foreground">{source.baseUrl}</div>
                  </td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">
                    {source.lastSyncedAt ? dateFormatter.format(new Date(source.lastSyncedAt)) : "Belum pernah"}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <Button
                      size="sm"
                      onClick={() => handleSync(source)}
                      disabled={isPending && pendingId === source.id}
                    >
                      {isPending && pendingId === source.id ? "Menyinkronkan..." : "Sinkron sekarang"}
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <section className="space-y-4">
        <header>
          <h2 className="text-base font-semibold text-foreground">Konten yang Tersedia</h2>
          <p className="text-xs text-muted-foreground">
            Aktifkan atau sembunyikan konten eksternal tanpa menghapus data aslinya.
          </p>
        </header>

        <div className="overflow-hidden rounded-2xl border border-border bg-card/95 shadow-sm supports-[backdrop-filter]:bg-card/80 supports-[backdrop-filter]:backdrop-blur">
          <table className="min-w-full divide-y divide-border/70 text-sm">
            <thead className="bg-muted/80">
              <tr>
                <th className="px-4 py-3 text-left font-semibold text-muted-foreground">Judul</th>
                <th className="px-4 py-3 text-left font-semibold text-muted-foreground">Sumber</th>
                <th className="px-4 py-3 text-left font-semibold text-muted-foreground">Jenis</th>
                <th className="px-4 py-3 text-left font-semibold text-muted-foreground">Terbit</th>
                <th className="w-28 px-4 py-3 text-center font-semibold text-muted-foreground">Tayang</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/70">
              {items.map((item) => (
                <tr key={item.id} className="bg-card/80">
                  <td className="px-4 py-3">
                    <div className="font-medium text-foreground">{item.title}</div>
                    {item.summary ? (
                      <p className="mt-1 text-xs text-muted-foreground line-clamp-2">{item.summary}</p>
                    ) : null}
                    <a
                      href={item.url}
                      target="_blank"
                      rel="noreferrer"
                      className="mt-1 inline-block text-xs font-medium text-primary hover:underline"
                    >
                      {item.url}
                    </a>
                  </td>
                  <td className="px-4 py-3 text-xs text-muted-foreground">{item.sourceName}</td>
                  <td className="px-4 py-3 text-xs uppercase text-muted-foreground">{item.kind}</td>
                  <td className="px-4 py-3 text-xs text-muted-foreground">
                    {item.publishedAt ? dateFormatter.format(new Date(item.publishedAt)) : "-"}
                  </td>
                  <td className="px-4 py-3 text-center">
                    <Switch
                      checked={item.visible}
                      onChange={(event) => handleToggleVisibility(item, event.target.checked)}
                      disabled={isPending && pendingId === item.id}
                      aria-label={`Atur visibilitas ${item.title}`}
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  );
}
