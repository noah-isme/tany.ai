import type { AnalyticsSummary } from "@/lib/types/admin";

type SummaryMetrics = {
  chats: number;
  avgResponse: number;
  successRate: number;
  conversions: number;
  uniqueUsers: number;
  conversionRate: number;
};

type DailyChartPoint = {
  date: string;
  totalChats: number;
  successRate: number;
  conversions: number;
};

type ProviderChartPoint = {
  provider: string;
  totalChats: number;
  successRate: number;
};

export function calculateSummaryMetrics(summary?: AnalyticsSummary | null): SummaryMetrics {
  const chats = summary?.totalChats ?? 0;
  const conversions = summary?.conversions ?? 0;
  return {
    chats,
    avgResponse: summary?.avgResponseTime ?? 0,
    successRate: summary?.successRate ?? 0,
    conversions,
    uniqueUsers: summary?.uniqueUsers ?? 0,
    conversionRate: chats === 0 ? 0 : conversions / chats,
  };
}

export function buildDailyChart(summary?: AnalyticsSummary | null): DailyChartPoint[] {
  return (summary?.daily ?? []).map((item) => ({
    date: item.date,
    totalChats: item.totalChats,
    successRate: item.successRate * 100,
    conversions: item.conversions,
  }));
}

export function buildProviderChart(summary?: AnalyticsSummary | null): ProviderChartPoint[] {
  const breakdown = summary?.providerBreakdown ?? {};
  return Object.entries(breakdown).map(([provider, stats]) => ({
    provider,
    totalChats: stats.totalChats,
    successRate: stats.successRate * 100,
  }));
}
