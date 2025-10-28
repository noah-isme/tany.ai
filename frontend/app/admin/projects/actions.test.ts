import { beforeEach, describe, expect, it, vi } from "vitest";

import { ApiError } from "@/lib/api-client";

import {
  createProjectAction,
  featureProjectAction,
} from "./actions";

const createProjectMock = vi.fn();
const featureProjectMock = vi.fn();

vi.mock("next/cache", () => ({
  revalidatePath: vi.fn(),
}));

vi.mock("@/lib/admin-api", () => ({
  createProject: (...args: unknown[]) => createProjectMock(...args),
  updateProject: vi.fn(),
  deleteProject: vi.fn(),
  reorderProjects: vi.fn(),
  featureProject: (...args: unknown[]) => featureProjectMock(...args),
}));

describe("projects actions", () => {
  beforeEach(() => {
    createProjectMock.mockReset();
    featureProjectMock.mockReset();
  });

  it("validates project title", async () => {
    const result = await createProjectAction({
      title: "",
      description: "",
      tech_stack: [],
      image_url: "",
      project_url: "",
      category: "",
      duration_label: "",
      price_label: "",
      budget_label: "",
      is_featured: false,
    });

    expect(result.success).toBe(false);
    expect(result.fieldErrors?.title).toBeDefined();
    expect(createProjectMock).not.toHaveBeenCalled();
  });

  it("creates project with trimmed tech stack", async () => {
    const project = {
      id: "prj-1",
      title: "AI Assistant",
      description: "",
      tech_stack: ["Next.js", "LangChain"],
      image_url: "https://example.com/image.png",
      project_url: "https://example.com",
      category: "AI",
      duration_label: "8 minggu",
      price_label: "Growth",
      budget_label: "IDR 120Jt",
      order: 0,
      is_featured: false,
    };
    createProjectMock.mockResolvedValue(project);

    const result = await createProjectAction({
      title: "AI Assistant",
      description: "",
      tech_stack: ["Next.js", "  LangChain  ", ""],
      image_url: "https://example.com/image.png",
      project_url: "https://example.com",
      category: "AI",
      duration_label: "8 minggu",
      price_label: "Growth",
      budget_label: "IDR 120Jt",
      is_featured: false,
    });

    expect(result.success).toBe(true);
    expect(createProjectMock).toHaveBeenCalledWith({
      title: "AI Assistant",
      description: "",
      tech_stack: ["Next.js", "LangChain"],
      image_url: "https://example.com/image.png",
      project_url: "https://example.com",
      category: "AI",
      duration_label: "8 minggu",
      price_label: "Growth",
      budget_label: "IDR 120Jt",
      is_featured: false,
    });
  });

  it("handles feature api error", async () => {
    featureProjectMock.mockRejectedValue(new ApiError("Gagal", 500));

    const result = await featureProjectAction({ id: "prj-1", is_featured: true });

    expect(result.success).toBe(false);
    expect(result.error).toBe("Gagal");
  });

  it("toggles project featured flag", async () => {
    const featured = {
      id: "prj-1",
      title: "AI Assistant",
      description: "",
      tech_stack: [],
      image_url: "",
      project_url: "",
      category: "",
      duration_label: "",
      price_label: "",
      budget_label: "",
      order: 0,
      is_featured: true,
    };
    featureProjectMock.mockResolvedValue(featured);

    const result = await featureProjectAction({ id: "prj-1", is_featured: true });

    expect(result.success).toBe(true);
    expect(result.data).toEqual(featured);
    expect(featureProjectMock).toHaveBeenCalledWith("prj-1", true);
  });
});
