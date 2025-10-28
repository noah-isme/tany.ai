import { AnalyticsDashboard } from "@/components/admin/analytics/AnalyticsDashboard";
import {
  fetchAnalyticsEvents,
  fetchAnalyticsLeads,
  fetchAnalyticsSummary,
} from "@/lib/admin-api";

export const dynamic = "force-dynamic";

export default async function AnalyticsPage() {
  const now = new Date();
  const from = new Date(now.getTime() - 6 * 24 * 60 * 60 * 1000);
  const fromISO = from.toISOString();
  const toISO = now.toISOString();

  const [summary, events, leads] = await Promise.all([
    fetchAnalyticsSummary({ from: fromISO, to: toISO }),
    fetchAnalyticsEvents({ from: fromISO, to: toISO, limit: 25, page: 1 }),
    fetchAnalyticsLeads({ from: fromISO, to: toISO, limit: 25, page: 1 }),
  ]);

  return (
    <AnalyticsDashboard
      initialSummary={summary}
      initialEvents={events}
      initialLeads={leads}
      defaultRange={{ from: fromISO, to: toISO }}
    />
  );
}
