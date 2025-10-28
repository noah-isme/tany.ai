import { apiFetch } from "./api-client";
import type {
  AnalyticsEvent,
  AnalyticsSummary,
  ApiListParams,
  ExternalItem,
  ExternalSource,
  PaginatedResponse,
  PersonalizationSummary,
  Profile,
  Project,
  Service,
  Skill,
} from "./types/admin";

export type AnalyticsFilter = {
  from?: string;
  to?: string;
  source?: string;
  provider?: string;
  page?: number;
  limit?: number;
  type?: string;
};

export async function fetchProfile(): Promise<Profile> {
  const response = await apiFetch<{ data: Profile }>("/api/admin/profile");
  return response.data;
}

export async function updateProfile(payload: Partial<Profile>): Promise<Profile> {
  const response = await apiFetch<{ data: Profile }>("/api/admin/profile", {
    method: "PUT",
    body: payload,
  });
  return response.data;
}

export async function fetchSkills(): Promise<PaginatedResponse<Skill>> {
  return apiFetch<PaginatedResponse<Skill>>("/api/admin/skills?limit=100&sort=order&dir=asc");
}

export async function createSkill(payload: { name: string; order?: number | null }): Promise<Skill> {
  const response = await apiFetch<{ data: Skill }>("/api/admin/skills", {
    method: "POST",
    body: payload,
  });
  return response.data;
}

export async function updateSkill(id: string, payload: { name: string; order?: number | null }): Promise<Skill> {
  const response = await apiFetch<{ data: Skill }>(`/api/admin/skills/${id}`, {
    method: "PUT",
    body: payload,
  });
  return response.data;
}

export async function deleteSkill(id: string): Promise<void> {
  await apiFetch(`/api/admin/skills/${id}`, { method: "DELETE" });
}

export async function reorderSkills(items: { id: string; order: number }[]): Promise<void> {
  await apiFetch("/api/admin/skills/reorder", {
    method: "PATCH",
    body: items,
  });
}

export async function fetchServices(): Promise<PaginatedResponse<Service>> {
  return apiFetch<PaginatedResponse<Service>>("/api/admin/services?limit=100&sort=order&dir=asc");
}

export async function createService(payload: Partial<Service>): Promise<Service> {
  const response = await apiFetch<{ data: Service }>("/api/admin/services", {
    method: "POST",
    body: payload,
  });
  return response.data;
}

export async function updateService(id: string, payload: Partial<Service>): Promise<Service> {
  const response = await apiFetch<{ data: Service }>(`/api/admin/services/${id}`, {
    method: "PUT",
    body: payload,
  });
  return response.data;
}

export async function deleteService(id: string): Promise<void> {
  await apiFetch(`/api/admin/services/${id}`, { method: "DELETE" });
}

export async function reorderServices(items: { id: string; order: number }[]): Promise<void> {
  await apiFetch("/api/admin/services/reorder", {
    method: "PATCH",
    body: items,
  });
}

export async function toggleService(id: string, isActive?: boolean): Promise<Service> {
  const body = typeof isActive === "boolean" ? { is_active: isActive } : undefined;
  const response = await apiFetch<{ data: Service }>(`/api/admin/services/${id}/toggle`, {
    method: "PATCH",
    body,
  });
  return response.data;
}

export async function fetchProjects(): Promise<PaginatedResponse<Project>> {
  return apiFetch<PaginatedResponse<Project>>("/api/admin/projects?limit=100&sort=order&dir=asc");
}

export async function createProject(payload: Partial<Project>): Promise<Project> {
  const response = await apiFetch<{ data: Project }>("/api/admin/projects", {
    method: "POST",
    body: payload,
  });
  return response.data;
}

export async function updateProject(id: string, payload: Partial<Project>): Promise<Project> {
  const response = await apiFetch<{ data: Project }>(`/api/admin/projects/${id}`, {
    method: "PUT",
    body: payload,
  });
  return response.data;
}

export async function deleteProject(id: string): Promise<void> {
  await apiFetch(`/api/admin/projects/${id}`, { method: "DELETE" });
}

export async function reorderProjects(items: { id: string; order: number }[]): Promise<void> {
  await apiFetch("/api/admin/projects/reorder", {
    method: "PATCH",
    body: items,
  });
}

export async function featureProject(id: string, isFeatured: boolean): Promise<Project> {
  const response = await apiFetch<{ data: Project }>(`/api/admin/projects/${id}/feature`, {
    method: "PATCH",
    body: { is_featured: isFeatured },
  });
  return response.data;
}

export async function fetchExternalSources(): Promise<PaginatedResponse<ExternalSource>> {
  return apiFetch<PaginatedResponse<ExternalSource>>(
    "/api/admin/external/sources?limit=50&sort=name&dir=asc",
  );
}

export async function syncExternalSource(
  id: string,
): Promise<{ itemsUpserted: number; message?: string; etag?: string; lastModified?: string }> {
  const response = await apiFetch<{ data: { itemsUpserted: number; message?: string; etag?: string; lastModified?: string } }>(
    `/api/admin/external/sources/${id}/sync`,
    { method: "POST" },
  );
  return response.data;
}

export async function fetchExternalItems(
  params: Pick<ApiListParams, "sort" | "dir"> & { kind?: string; visible?: boolean } = {},
): Promise<PaginatedResponse<ExternalItem>> {
  const search = new URLSearchParams({ limit: "100" });
  if (params.sort) {
    search.set("sort", params.sort);
  } else {
    search.set("sort", "published_at");
  }
  if (params.dir) {
    search.set("dir", params.dir);
  } else {
    search.set("dir", "desc");
  }
  if (params.kind) {
    search.set("kind", params.kind);
  }
  if (typeof params.visible === "boolean") {
    search.set("visible", String(params.visible));
  }
  return apiFetch<PaginatedResponse<ExternalItem>>(
    `/api/admin/external/items?${search.toString()}`,
  );
}

export async function setExternalItemVisibility(id: string, visible: boolean): Promise<ExternalItem> {
  const response = await apiFetch<{ data: ExternalItem }>(
    `/api/admin/external/items/${id}/visibility`,
    {
      method: "PATCH",
      body: { visible },
    },
  );
  return response.data;
}

export async function fetchPersonalizationSummary(): Promise<PersonalizationSummary> {
  const response = await apiFetch<{ data: PersonalizationSummary }>("/api/admin/personalization");
  return response.data;
}

export async function updatePersonalizationWeight(weight: number): Promise<number> {
  const response = await apiFetch<{ data: { weight: number } }>("/api/admin/personalization/weight", {
    method: "PATCH",
    body: { weight },
  });
  return response.data.weight;
}

export async function reindexPersonalization(): Promise<{ indexed: number }> {
  const response = await apiFetch<{ data: { indexed: number } }>("/api/admin/personalization/reindex", {
    method: "POST",
  });
  return response.data;
}

export async function resetPersonalization(): Promise<void> {
  await apiFetch("/api/admin/personalization/reset", { method: "POST" });
}

function buildAnalyticsQuery(params: AnalyticsFilter = {}): string {
  const search = new URLSearchParams();
  if (params.from) {
    search.set("from", params.from);
  }
  if (params.to) {
    search.set("to", params.to);
  }
  if (params.source) {
    search.set("source", params.source);
  }
  if (params.provider) {
    search.set("provider", params.provider);
  }
  if (params.page) {
    search.set("page", String(params.page));
  }
  if (params.limit) {
    search.set("limit", String(params.limit));
  }
  if (params.type) {
    search.set("type", params.type);
  }
  const query = search.toString();
  return query ? `?${query}` : "";
}

export async function fetchAnalyticsSummary(params: AnalyticsFilter = {}): Promise<AnalyticsSummary> {
  const query = buildAnalyticsQuery({
    from: params.from,
    to: params.to,
    source: params.source,
    provider: params.provider,
  });
  const response = await apiFetch<{ data: AnalyticsSummary }>(`/api/admin/analytics/summary${query}`);
  return response.data;
}

export async function fetchAnalyticsEvents(params: AnalyticsFilter = {}): Promise<PaginatedResponse<AnalyticsEvent>> {
  const query = buildAnalyticsQuery(params);
  return apiFetch<PaginatedResponse<AnalyticsEvent>>(`/api/admin/analytics/events${query}`);
}

export async function fetchAnalyticsLeads(params: AnalyticsFilter = {}): Promise<PaginatedResponse<AnalyticsEvent>> {
  const query = buildAnalyticsQuery(params);
  return apiFetch<PaginatedResponse<AnalyticsEvent>>(`/api/admin/analytics/leads${query}`);
}
