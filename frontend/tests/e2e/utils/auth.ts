import { expect, Page } from "@playwright/test";

const DEFAULT_API_PORT = process.env.MOCK_API_PORT ?? "4000";
export const API_BASE_URL =
  process.env.PLAYWRIGHT_API_BASE_URL ?? `http://127.0.0.1:${DEFAULT_API_PORT}`;

type ApiRequestOptions = {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  data?: Record<string, unknown> | unknown[];
};

type SkillResponse = { id: string; name: string; order: number };
type ServiceResponse = {
  id: string;
  name: string;
  description?: string | null;
  is_active: boolean;
  order: number;
  price_min?: number | null;
  price_max?: number | null;
  currency?: string | null;
  duration_label?: string | null;
};
type ProjectResponse = {
  id: string;
  title: string;
  description?: string | null;
  tech_stack?: string[];
  image_url?: string | null;
  project_url?: string | null;
  category?: string | null;
  order: number;
  is_featured: boolean;
};

function authHeaders(token: string): Record<string, string> {
  return {
    Authorization: `Bearer ${token}`,
    Accept: "application/json",
  };
}

async function apiRequest<T>(
  page: Page,
  token: string,
  path: string,
  { method = "GET", data }: ApiRequestOptions = {},
): Promise<T> {
  const response = await page.request.fetch(`${API_BASE_URL}${path}`, {
    method,
    data,
    headers: {
      ...authHeaders(token),
      "Content-Type": "application/json",
    },
  });

  expect(response.ok()).toBeTruthy();
  const contentType = response.headers()["content-type"] ?? "";
  if (!contentType.includes("application/json")) {
    return {} as T;
  }
  return (await response.json()) as T;
}

export async function authenticateAdmin(page: Page, email: string, password: string): Promise<string> {
  const response = await page.request.post(`${API_BASE_URL}/api/auth/login`, {
    data: { email, password },
  });

  expect(response.ok()).toBeTruthy();
  const payload = (await response.json()) as { accessToken: string };
  const currentUrl = page.url() || `http://localhost:${process.env.PORT ?? "3000"}/login`;
  const { origin } = new URL(currentUrl);

  await page.context().addCookies([
    {
      name: "ta_access",
      value: payload.accessToken,
      url: origin,
      httpOnly: true,
      sameSite: "Lax",
      secure: false,
    },
  ]);

  return payload.accessToken;
}

export async function fetchSkills(page: Page, token: string): Promise<SkillResponse[]> {
  const response = await apiRequest<{ items: SkillResponse[] }>(
    page,
    token,
    "/api/admin/skills?limit=100&sort=order&dir=asc",
  );
  return response.items;
}

export async function createSkill(page: Page, token: string, name: string): Promise<SkillResponse> {
  const response = await apiRequest<{ data: SkillResponse }>(page, token, "/api/admin/skills", {
    method: "POST",
    data: { name },
  });
  return response.data;
}

export async function reorderSkills(
  page: Page,
  token: string,
  items: { id: string; order: number }[],
): Promise<void> {
  await apiRequest(page, token, "/api/admin/skills/reorder", { method: "PATCH", data: items });
}

export async function deleteSkill(page: Page, token: string, id: string): Promise<void> {
  await apiRequest(page, token, `/api/admin/skills/${id}`, { method: "DELETE" });
}

export async function fetchServices(page: Page, token: string): Promise<ServiceResponse[]> {
  const response = await apiRequest<{ items: ServiceResponse[] }>(
    page,
    token,
    "/api/admin/services?limit=100&sort=order&dir=asc",
  );
  return response.items;
}

export async function createService(
  page: Page,
  token: string,
  payload: Partial<ServiceResponse> & { name: string },
): Promise<ServiceResponse> {
  const response = await apiRequest<{ data: ServiceResponse }>(page, token, "/api/admin/services", {
    method: "POST",
    data: payload,
  });
  return response.data;
}

export async function toggleService(
  page: Page,
  token: string,
  id: string,
  isActive: boolean,
): Promise<ServiceResponse> {
  const response = await apiRequest<{ data: ServiceResponse }>(
    page,
    token,
    `/api/admin/services/${id}/toggle`,
    {
      method: "PATCH",
      data: { is_active: isActive },
    },
  );
  return response.data;
}

export async function fetchProjects(page: Page, token: string): Promise<ProjectResponse[]> {
  const response = await apiRequest<{ items: ProjectResponse[] }>(
    page,
    token,
    "/api/admin/projects?limit=100&sort=order&dir=asc",
  );
  return response.items;
}

export async function createProject(
  page: Page,
  token: string,
  payload: Partial<ProjectResponse> & { title: string },
): Promise<ProjectResponse> {
  const response = await apiRequest<{ data: ProjectResponse }>(page, token, "/api/admin/projects", {
    method: "POST",
    data: payload,
  });
  return response.data;
}

export async function deleteProject(page: Page, token: string, id: string): Promise<void> {
  await apiRequest(page, token, `/api/admin/projects/${id}`, { method: "DELETE" });
}

export async function deleteService(page: Page, token: string, id: string): Promise<void> {
  await apiRequest(page, token, `/api/admin/services/${id}`, { method: "DELETE" });
}

export async function updateProfile(
  page: Page,
  token: string,
  payload: Record<string, unknown>,
): Promise<{ title?: string | null; location?: string | null }> {
  const response = await apiRequest<{ data: { title?: string | null; location?: string | null } }>(
    page,
    token,
    "/api/admin/profile",
    { method: "PUT", data: payload },
  );
  return response.data;
}
