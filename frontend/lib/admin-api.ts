import { apiFetch } from "./api-client";
import type { Profile, Skill, Service, Project, PaginatedResponse } from "./types/admin";

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
