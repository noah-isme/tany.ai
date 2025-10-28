import { describe, expect, it } from "vitest";

import type { AnalyticsSummary } from "@/lib/types/admin";
import { buildDailyChart, buildProviderChart, calculateSummaryMetrics } from "../analytics/utils";

describe("analytics utils", () => {
  const summary: AnalyticsSummary = {
    rangeStart: "2025-01-01T00:00:00Z",
    rangeEnd: "2025-01-07T23:59:59Z",
    totalChats: 20,
    avgResponseTime: 150,
    successRate: 0.85,
    uniqueUsers: 8,
    conversions: 4,
    providerBreakdown: {
      gemini: {
        totalChats: 12,
        avgResponseTime: 140,
        successRate: 0.9,
      },
      mock: {
        totalChats: 8,
        avgResponseTime: 170,
        successRate: 0.75,
      },
    },
    daily: [
      {
        date: "2025-01-01T00:00:00Z",
        totalChats: 5,
        avgResponseTime: 160,
        successRate: 0.8,
        conversions: 1,
      },
      {
        date: "2025-01-02T00:00:00Z",
        totalChats: 7,
        avgResponseTime: 140,
        successRate: 0.9,
        conversions: 2,
      },
    ],
  };

  it("calculates summary metrics with fallback values", () => {
    const metrics = calculateSummaryMetrics(summary);
    expect(metrics.chats).toBe(20);
    expect(metrics.avgResponse).toBe(150);
    expect(metrics.successRate).toBeCloseTo(0.85);
    expect(metrics.conversionRate).toBeCloseTo(0.2);
    expect(metrics.uniqueUsers).toBe(8);
  });

  it("builds daily chart data with percentage conversion", () => {
    const result = buildDailyChart(summary);
    expect(result).toHaveLength(2);
    expect(result[0]).toEqual({
      date: "2025-01-01T00:00:00Z",
      totalChats: 5,
      successRate: 80,
      conversions: 1,
    });
  });

  it("builds provider chart data with percentage conversion", () => {
    const result = buildProviderChart(summary);
    expect(result).toEqual([
      { provider: "gemini", totalChats: 12, successRate: 90 },
      { provider: "mock", totalChats: 8, successRate: 75 },
    ]);
  });
});
