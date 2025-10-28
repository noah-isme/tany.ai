"use client";

import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Line,
  LineChart,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { useCallback, useEffect, useMemo, useState } from "react";

import {
  AnalyticsFilter,
  fetchAnalyticsEvents,
  fetchAnalyticsLeads,
  fetchAnalyticsSummary,
} from "@/lib/admin-api";
import type {
  AnalyticsEvent,
  AnalyticsSummary,
  PaginatedResponse,
} from "@/lib/types/admin";
import {
  buildDailyChart,
  buildProviderChart,
  calculateSummaryMetrics,
} from "./utils";

const AUTO_REFRESH_INTERVAL = 5 * 60 * 1000; // 5 minutes
const PROVIDER_COLORS = ["#6366f1", "#22c55e", "#f97316", "#a855f7", "#0ea5e9", "#ef4444"];

type AnalyticsDashboardProps = {
  initialSummary: AnalyticsSummary;
  initialEvents: PaginatedResponse<AnalyticsEvent>;
  initialLeads: PaginatedResponse<AnalyticsEvent>;
  defaultRange: { from: string; to: string };
};

type FilterState = {
  from: string;
  to: string;
  provider: string;
  source: string;
};

type FetchState = {
  summary: AnalyticsSummary;
  events: PaginatedResponse<AnalyticsEvent>;
  leads: PaginatedResponse<AnalyticsEvent>;
};

function toISODate(value: string): string {
  if (!value) {
    return "";
  }
  if (value.includes("T")) {
    return value;
  }
  return new Date(`${value}T00:00:00Z`).toISOString();
}

function formatPercent(value: number): string {
  return `${(value * 100).toFixed(1)}%`;
}

function formatLatency(value: number): string {
  if (value >= 1000) {
    return `${(value / 1000).toFixed(2)}s`;
  }
  return `${value}ms`;
}

function formatDateLabel(value: string): string {
  const date = new Date(value);
  return date.toLocaleDateString("id-ID", { month: "short", day: "numeric" });
}

export function AnalyticsDashboard({
  initialSummary,
  initialEvents,
  initialLeads,
  defaultRange,
}: AnalyticsDashboardProps) {
  const [filters, setFilters] = useState<FilterState>({
    from: defaultRange.from.split("T")[0],
    to: defaultRange.to.split("T")[0],
    provider: "",
    source: "",
  });
  const [state, setState] = useState<FetchState>({
    summary: initialSummary,
    events: initialEvents,
    leads: initialLeads,
  });
  const [loading, setLoading] = useState(false);
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date());
  const [initialLoad, setInitialLoad] = useState(true);

  const providerOptions = useMemo(
    () => Object.keys(state.summary?.providerBreakdown ?? {}),
    [state.summary?.providerBreakdown],
  );

  const summaryMetrics = useMemo(
    () => calculateSummaryMetrics(state.summary),
    [state.summary],
  );

  const dailyChartData = useMemo(
    () => buildDailyChart(state.summary),
    [state.summary],
  );

  const providerChartData = useMemo(
    () => buildProviderChart(state.summary),
    [state.summary],
  );

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      const request: AnalyticsFilter = {
        from: toISODate(filters.from),
        to: toISODate(filters.to),
        provider: filters.provider || undefined,
        source: filters.source || undefined,
      };
      const [summary, events, leads] = await Promise.all([
        fetchAnalyticsSummary(request),
        fetchAnalyticsEvents({ ...request, limit: 25, page: 1 }),
        fetchAnalyticsLeads({ ...request, limit: 25, page: 1 }),
      ]);
      setState({ summary, events, leads });
      setLastUpdated(new Date());
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    if (initialLoad) {
      setInitialLoad(false);
      return;
    }
    void refresh();
  }, [refresh, initialLoad]);

  useEffect(() => {
    const id = setInterval(() => {
      void refresh();
    }, AUTO_REFRESH_INTERVAL);
    return () => {
      clearInterval(id);
    };
  }, [refresh]);

  const handleFilterChange = (key: keyof FilterState, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-100">Analytics</h1>
          <p className="text-sm text-slate-600 dark:text-slate-400">
            Pantau performa chat tany.ai, engagement klien, dan insight konversi secara real-time.
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-3 rounded-2xl border border-slate-200 bg-white/80 p-4 text-xs shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
          <span className="font-medium text-slate-500 dark:text-slate-300">Auto refresh</span>
          <span className="rounded-full bg-emerald-100 px-2 py-1 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-200">
            setiap 5 menit
          </span>
          <span className="text-slate-400 dark:text-slate-500">
            Terakhir diperbarui: {lastUpdated.toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit" })}
          </span>
          {loading && <span className="text-indigo-500">· Memuat data…</span>}
        </div>
      </div>

      <section className="rounded-2xl border border-slate-200 bg-white/80 p-4 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <div className="rounded-xl border border-slate-200/60 bg-white/90 p-4 dark:border-slate-700/60 dark:bg-slate-900/70">
            <p className="text-xs uppercase tracking-[0.35em] text-slate-400">Total Chat</p>
            <p className="mt-3 text-2xl font-semibold text-slate-900 dark:text-slate-100">{summaryMetrics.chats}</p>
            <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">Rentang waktu terpilih</p>
          </div>
          <div className="rounded-xl border border-slate-200/60 bg-white/90 p-4 dark:border-slate-700/60 dark:bg-slate-900/70">
            <p className="text-xs uppercase tracking-[0.35em] text-slate-400">Rata-rata Latensi</p>
            <p className="mt-3 text-2xl font-semibold text-slate-900 dark:text-slate-100">
              {summaryMetrics.avgResponse.toFixed(0)} ms
            </p>
            <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">Hitung dari semua interaksi</p>
          </div>
          <div className="rounded-xl border border-slate-200/60 bg-white/90 p-4 dark:border-slate-700/60 dark:bg-slate-900/70">
            <p className="text-xs uppercase tracking-[0.35em] text-slate-400">Success Rate</p>
            <p className="mt-3 text-2xl font-semibold text-emerald-500">
              {formatPercent(summaryMetrics.successRate)}
            </p>
            <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">Respons AI tanpa fallback</p>
          </div>
          <div className="rounded-xl border border-slate-200/60 bg-white/90 p-4 dark:border-slate-700/60 dark:bg-slate-900/70">
            <p className="text-xs uppercase tracking-[0.35em] text-slate-400">Konversi</p>
            <p className="mt-3 text-2xl font-semibold text-slate-900 dark:text-slate-100">{summaryMetrics.conversions}</p>
            <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">
              Conversion rate {formatPercent(summaryMetrics.conversionRate)} · {summaryMetrics.uniqueUsers} pengunjung unik
            </p>
          </div>
        </div>
      </section>

      <section className="rounded-2xl border border-slate-200 bg-white/80 p-4 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
        <div className="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
          <div>
            <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Filter &amp; Rentang Waktu</h2>
            <p className="text-xs text-slate-500 dark:text-slate-400">
              Sesuaikan periode analisis, sumber, dan provider untuk mempersempit insight.
            </p>
          </div>
          <div className="flex flex-wrap items-center gap-3 text-xs">
            <label className="flex flex-col gap-1">
              <span className="text-slate-500">Dari</span>
              <input
                type="date"
                value={filters.from}
                onChange={(event) => handleFilterChange("from", event.target.value)}
                className="rounded-md border border-slate-300 bg-white px-2 py-1 text-slate-700 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/40 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-200"
              />
            </label>
            <label className="flex flex-col gap-1">
              <span className="text-slate-500">Sampai</span>
              <input
                type="date"
                value={filters.to}
                onChange={(event) => handleFilterChange("to", event.target.value)}
                className="rounded-md border border-slate-300 bg-white px-2 py-1 text-slate-700 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/40 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-200"
              />
            </label>
            <label className="flex flex-col gap-1">
              <span className="text-slate-500">Provider</span>
              <select
                value={filters.provider}
                onChange={(event) => handleFilterChange("provider", event.target.value)}
                className="min-w-[140px] rounded-md border border-slate-300 bg-white px-2 py-1 text-slate-700 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/40 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-200"
              >
                <option value="">Semua provider</option>
                {providerOptions.map((provider) => (
                  <option key={provider} value={provider}>
                    {provider}
                  </option>
                ))}
              </select>
            </label>
            <label className="flex flex-col gap-1">
              <span className="text-slate-500">Sumber</span>
              <input
                value={filters.source}
                onChange={(event) => handleFilterChange("source", event.target.value)}
                placeholder="mis. landing-page"
                className="rounded-md border border-slate-300 bg-white px-2 py-1 text-slate-700 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/40 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-200"
              />
            </label>
          </div>
        </div>
      </section>

      <section className="grid gap-4 xl:grid-cols-[2fr_1fr]">
        <div className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Tren Chat Harian</h2>
              <p className="text-xs text-slate-500 dark:text-slate-400">
                Volume chat vs success rate. Sumbu kiri menunjukkan jumlah chat, sumbu kanan success rate.
              </p>
            </div>
          </div>
          <div className="mt-4 h-64">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={dailyChartData} margin={{ top: 10, right: 40, bottom: 0, left: 0 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
                <XAxis dataKey="date" tickFormatter={formatDateLabel} stroke="#94a3b8" />
                <YAxis yAxisId="left" stroke="#94a3b8" allowDecimals={false} />
                <YAxis yAxisId="right" orientation="right" stroke="#94a3b8" domain={[0, 100]} />
                <Tooltip
                  contentStyle={{ borderRadius: 12 }}
                  labelFormatter={(value) => new Date(value).toLocaleDateString("id-ID", { dateStyle: "medium" })}
                  formatter={(value, name) => {
                    if (name === "successRate") {
                      return [`${(value as number).toFixed(1)}%`, "Success rate"];
                    }
                    return [value, name === "totalChats" ? "Total chat" : name];
                  }}
                />
                <Line
                  type="monotone"
                  dataKey="totalChats"
                  name="Total chat"
                  stroke="#6366f1"
                  strokeWidth={2}
                  dot={false}
                  yAxisId="left"
                />
                <Line
                  type="monotone"
                  dataKey="successRate"
                  name="Success rate"
                  stroke="#22c55e"
                  strokeWidth={2}
                  dot={false}
                  yAxisId="right"
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
          <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Performa Provider</h2>
          <p className="text-xs text-slate-500 dark:text-slate-400">
            Bandingkan distribusi chat dan success rate antar provider model.
          </p>
          <div className="mt-4 h-64">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={providerChartData}
                  dataKey="totalChats"
                  nameKey="provider"
                  innerRadius={50}
                  outerRadius={80}
                  paddingAngle={3}
                >
                  {providerChartData.map((entry, index) => (
                    <Cell key={entry.provider} fill={PROVIDER_COLORS[index % PROVIDER_COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip
                  formatter={(value, name, payload) => {
                    const item = payload?.payload as { successRate: number };
                    return [
                      `${value} chat · success ${(item.successRate ?? 0).toFixed(1)}%`,
                      String(name),
                    ];
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
          </div>
          <ul className="mt-4 space-y-2 text-sm">
            {providerChartData.length === 0 && (
              <li className="text-slate-500 dark:text-slate-400">Belum ada data provider untuk rentang ini.</li>
            )}
            {providerChartData.map((item, index) => (
              <li key={item.provider} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span
                    className="inline-block h-2 w-2 rounded-full"
                    style={{ backgroundColor: PROVIDER_COLORS[index % PROVIDER_COLORS.length] }}
                  />
                  <span className="text-slate-600 dark:text-slate-300">{item.provider}</span>
                </div>
                <div className="text-xs text-slate-500 dark:text-slate-400">
                  {item.totalChats} chat · success {item.successRate.toFixed(1)}%
                </div>
              </li>
            ))}
          </ul>
        </div>
      </section>

      <section className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
        <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Lead &amp; Konversi</h2>
        <p className="text-xs text-slate-500 dark:text-slate-400">
          Monitor sumber lead terbaru dan konversi yang berasal dari chat.
        </p>
        <div className="mt-4 h-64">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={dailyChartData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
              <XAxis dataKey="date" tickFormatter={formatDateLabel} stroke="#94a3b8" />
              <YAxis stroke="#94a3b8" allowDecimals={false} />
              <Tooltip
                labelFormatter={(value) => new Date(value).toLocaleDateString("id-ID", { dateStyle: "medium" })}
                formatter={(value) => [`${value} lead`, "Konversi"]}
              />
              <Bar dataKey="conversions" name="Konversi" fill="#6366f1" radius={[6, 6, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
        <div className="mt-6 overflow-x-auto">
          <table className="min-w-full divide-y divide-slate-200 text-sm dark:divide-slate-800">
            <thead className="text-left text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
              <tr>
                <th className="px-4 py-2">Timestamp</th>
                <th className="px-4 py-2">Sumber</th>
                <th className="px-4 py-2">Provider</th>
                <th className="px-4 py-2">Metadata</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
              {state.leads.items.length === 0 && (
                <tr>
                  <td colSpan={4} className="px-4 py-6 text-center text-slate-500 dark:text-slate-400">
                    Belum ada lead tercatat pada rentang ini.
                  </td>
                </tr>
              )}
              {state.leads.items.map((lead) => (
                <tr key={lead.id} className="text-slate-600 dark:text-slate-300">
                  <td className="px-4 py-2">
                    {new Date(lead.timestamp).toLocaleString("id-ID", {
                      dateStyle: "medium",
                      timeStyle: "short",
                    })}
                  </td>
                  <td className="px-4 py-2">{lead.source}</td>
                  <td className="px-4 py-2">{lead.provider}</td>
                  <td className="px-4 py-2 text-xs">
                    {lead.metadata && Object.keys(lead.metadata).length > 0
                      ? JSON.stringify(lead.metadata)
                      : "-"}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <section className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
        <div className="flex flex-col gap-2">
          <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Log Interaksi Chat</h2>
          <p className="text-xs text-slate-500 dark:text-slate-400">
            Riwayat percakapan terbaru lengkap dengan latensi dan status sukses provider.
          </p>
        </div>
        <div className="mt-4 overflow-x-auto">
          <table className="min-w-full divide-y divide-slate-200 text-sm dark:divide-slate-800">
            <thead className="text-left text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
              <tr>
                <th className="px-4 py-2">Timestamp</th>
                <th className="px-4 py-2">Sumber</th>
                <th className="px-4 py-2">Provider</th>
                <th className="px-4 py-2">Durasi</th>
                <th className="px-4 py-2">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
              {state.events.items.length === 0 && (
                <tr>
                  <td colSpan={5} className="px-4 py-6 text-center text-slate-500 dark:text-slate-400">
                    Belum ada interaksi pada rentang ini.
                  </td>
                </tr>
              )}
              {state.events.items.map((event) => (
                <tr key={event.id} className="text-slate-600 dark:text-slate-300">
                  <td className="px-4 py-2">
                    {new Date(event.timestamp).toLocaleString("id-ID", {
                      dateStyle: "medium",
                      timeStyle: "short",
                    })}
                  </td>
                  <td className="px-4 py-2">{event.source}</td>
                  <td className="px-4 py-2">{event.provider}</td>
                  <td className="px-4 py-2">{formatLatency(event.durationMs)}</td>
                  <td className="px-4 py-2">
                    <span
                      className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium ${
                        event.success
                          ? "bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-200"
                          : "bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-200"
                      }`}
                    >
                      <span className="h-1.5 w-1.5 rounded-full bg-current" aria-hidden="true" />
                      {event.success ? "Berhasil" : "Fallback"}
                    </span>
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
