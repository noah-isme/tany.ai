"use server";

import { revalidatePath } from "next/cache";

import {
  createProject,
  deleteProject,
  featureProject,
  reorderProjects,
  updateProject,
} from "@/lib/admin-api";
import { ApiError, isApiError } from "@/lib/api-client";
import { type ActionResult } from "@/lib/action-result";
import { issuesToFieldErrors } from "@/lib/form-errors";
import type { Project } from "@/lib/types/admin";
import { projectSchema, type ProjectFormValues } from "@/lib/validators";

export async function createProjectAction(values: ProjectFormValues): Promise<ActionResult<Project>> {
  const parsed = projectSchema.safeParse(values);
  if (!parsed.success) {
    return {
      success: false,
      error: "Periksa data proyek.",
      fieldErrors: issuesToFieldErrors(parsed.error.issues),
    };
  }

  try {
    const created = await createProject(buildProjectPayload(parsed.data));
    revalidatePath("/admin/projects");
    revalidatePath("/admin");
    return { success: true, data: created, message: "Proyek ditambahkan." };
  } catch (error) {
    return normalizeProjectError(error);
  }
}

export async function updateProjectAction(params: {
  id: string;
  values: ProjectFormValues;
}): Promise<ActionResult<Project>> {
  const parsed = projectSchema.safeParse(params.values);
  if (!parsed.success) {
    return {
      success: false,
      error: "Periksa data proyek.",
      fieldErrors: issuesToFieldErrors(parsed.error.issues),
    };
  }

  try {
    const updated = await updateProject(params.id, buildProjectPayload(parsed.data));
    revalidatePath("/admin/projects");
    revalidatePath("/admin");
    return { success: true, data: updated, message: "Proyek diperbarui." };
  } catch (error) {
    return normalizeProjectError(error);
  }
}

export async function deleteProjectAction(params: { id: string }): Promise<ActionResult<null>> {
  try {
    await deleteProject(params.id);
    revalidatePath("/admin/projects");
    revalidatePath("/admin");
    return { success: true, data: null, message: "Proyek dihapus." };
  } catch (error) {
    return normalizeProjectError(error);
  }
}

export async function reorderProjectAction(params: {
  items: { id: string; order: number }[];
}): Promise<ActionResult<null>> {
  try {
    await reorderProjects(params.items);
    revalidatePath("/admin/projects");
    revalidatePath("/admin");
    return { success: true, data: null };
  } catch (error) {
    return normalizeProjectError(error);
  }
}

export async function featureProjectAction(params: {
  id: string;
  is_featured: boolean;
}): Promise<ActionResult<Project>> {
  try {
    const updated = await featureProject(params.id, params.is_featured);
    revalidatePath("/admin/projects");
    revalidatePath("/admin");
    return { success: true, data: updated };
  } catch (error) {
    return normalizeProjectError(error);
  }
}

function buildProjectPayload(values: ProjectFormValues) {
  return {
    title: values.title,
    description: values.description ?? "",
    tech_stack: values.tech_stack,
    image_url: values.image_url ?? "",
    project_url: values.project_url ?? "",
    category: values.category ?? "",
    is_featured: values.is_featured ?? false,
  };
}

function normalizeProjectError(error: unknown): ActionResult<never> {
  if (isApiError(error)) {
    const apiError = error as ApiError;
    const fieldErrors =
      apiError.details && typeof apiError.details === "object"
        ? Object.fromEntries(
            Object.entries(apiError.details as Record<string, unknown>).filter(([, value]) => typeof value === "string") as [
              string,
              string
            ][],
          )
        : undefined;
    return { success: false, error: apiError.message, fieldErrors };
  }
  return { success: false, error: "Operasi proyek gagal." };
}
