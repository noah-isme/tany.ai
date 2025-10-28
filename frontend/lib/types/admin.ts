export type Profile = {
  id: string;
  name: string;
  title: string;
  bio: string;
  email: string;
  phone: string;
  location: string;
  avatar_url: string;
  updated_at?: string;
};

export type Skill = {
  id: string;
  name: string;
  order: number;
};

export type Service = {
  id: string;
  name: string;
  description: string;
  price_min: number | null;
  price_max: number | null;
  currency: string;
  duration_label: string;
  is_active: boolean;
  order: number;
};

export type Project = {
  id: string;
  title: string;
  description: string;
  tech_stack: string[];
  image_url: string;
  project_url: string;
  category: string;
  duration_label: string;
  price_label: string;
  budget_label: string;
  order: number;
  is_featured: boolean;
};

export type ExternalSource = {
  id: string;
  name: string;
  baseUrl: string;
  sourceType: string;
  enabled: boolean;
  lastSyncedAt?: string;
  lastModified?: string;
};

export type ExternalItem = {
  id: string;
  sourceName: string;
  kind: string;
  title: string;
  summary?: string;
  url: string;
  visible: boolean;
  publishedAt?: string;
  metadata: Record<string, unknown>;
};

export type PaginatedResponse<T> = {
  items: T[];
  page: number;
  limit: number;
  total: number;
};

export type ApiListParams = {
  page?: number;
  limit?: number;
  sort?: string;
  dir?: "asc" | "desc";
};

export type AnalyticsProviderSnapshot = {
  totalChats: number;
  avgResponseTime: number;
  successRate: number;
};

export type AnalyticsDailyPoint = {
  date: string;
  totalChats: number;
  avgResponseTime: number;
  successRate: number;
  conversions: number;
};

export type AnalyticsSummary = {
  rangeStart: string;
  rangeEnd: string;
  totalChats: number;
  avgResponseTime: number;
  successRate: number;
  uniqueUsers: number;
  conversions: number;
  providerBreakdown: Record<string, AnalyticsProviderSnapshot>;
  daily: AnalyticsDailyPoint[];
};

export type AnalyticsEvent = {
  id: string;
  timestamp: string;
  eventType: string;
  source: string;
  provider: string;
  durationMs: number;
  success: boolean;
  userAgent?: string;
  metadata: Record<string, unknown>;
};

export type PersonalizationSummary = {
  enabled: boolean;
  provider: string;
  dimension: number;
  count: number;
  weight: number;
  lastReindexedAt?: string;
  lastResetAt?: string;
};
