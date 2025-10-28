import { resolveApiUrl } from "./env";

export type KnowledgeProfile = {
  name: string;
  title?: string;
  bio?: string;
  email?: string;
  phone?: string;
  location?: string;
  avatarUrl?: string;
  updatedAt?: string;
};

export type KnowledgeSkill = {
  name: string;
};

export type KnowledgeService = {
  id: string;
  name: string;
  description?: string;
  currency?: string;
  durationLabel?: string;
  priceRange?: string[];
  order: number;
};

export type KnowledgeProject = {
  id: string;
  title: string;
  description?: string;
  techStack: string[];
  projectUrl?: string;
  category?: string;
  durationLabel?: string;
  priceLabel?: string;
  budgetLabel?: string;
  isFeatured: boolean;
  order: number;
};

export type KnowledgeBase = {
  profile: KnowledgeProfile;
  skills: KnowledgeSkill[];
  services: KnowledgeService[];
  projects: KnowledgeProject[];
};

export async function fetchKnowledgeBase(): Promise<KnowledgeBase> {
  const path = "/api/v1/knowledge-base";
  if (typeof window === "undefined") {
    const { apiFetch } = await import("./api-client");
    return apiFetch<KnowledgeBase>(path, {
      withAuth: false,
      cache: "no-store",
    });
  }

  const response = await fetch(resolveApiUrl(path), {
    headers: {
      Accept: "application/json",
    },
    cache: "no-store",
  });
  if (!response.ok) {
    throw new Error(`Failed to load knowledge base: ${response.status}`);
  }
  return (await response.json()) as KnowledgeBase;
}
