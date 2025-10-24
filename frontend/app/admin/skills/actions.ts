"use server";

import { revalidatePath } from "next/cache";

import {
  createSkill,
  deleteSkill,
  reorderSkills,
  updateSkill,
} from "@/lib/admin-api";
import { ApiError, isApiError } from "@/lib/api-client";
import { type ActionResult } from "@/lib/action-result";
import { issuesToFieldErrors } from "@/lib/form-errors";
import type { Skill } from "@/lib/types/admin";
import { skillSchema, type SkillFormValues } from "@/lib/validators";

export async function createSkillAction(values: SkillFormValues): Promise<ActionResult<Skill>> {
  const parsed = skillSchema.safeParse(values);
  if (!parsed.success) {
    return {
      success: false,
      error: "Nama skill belum valid.",
      fieldErrors: issuesToFieldErrors(parsed.error.issues),
    };
  }

  try {
    const created = await createSkill({ name: parsed.data.name });
    revalidatePath("/admin/skills");
    revalidatePath("/admin");
    return { success: true, data: created, message: "Skill ditambahkan." };
  } catch (error) {
    return normalizeSkillError(error);
  }
}

export async function updateSkillAction(params: {
  id: string;
  values: SkillFormValues;
}): Promise<ActionResult<Skill>> {
  const parsed = skillSchema.safeParse(params.values);
  if (!parsed.success) {
    return {
      success: false,
      error: "Nama skill belum valid.",
      fieldErrors: issuesToFieldErrors(parsed.error.issues),
    };
  }
  try {
    const updated = await updateSkill(params.id, { name: parsed.data.name });
    revalidatePath("/admin/skills");
    revalidatePath("/admin");
    return { success: true, data: updated, message: "Skill diperbarui." };
  } catch (error) {
    return normalizeSkillError(error);
  }
}

export async function deleteSkillAction(params: { id: string }): Promise<ActionResult<null>> {
  try {
    await deleteSkill(params.id);
    revalidatePath("/admin/skills");
    revalidatePath("/admin");
    return { success: true, data: null, message: "Skill dihapus." };
  } catch (error) {
    return normalizeSkillError(error);
  }
}

export async function reorderSkillAction(params: {
  items: { id: string; order: number }[];
}): Promise<ActionResult<null>> {
  try {
    await reorderSkills(params.items);
    revalidatePath("/admin/skills");
    revalidatePath("/admin");
    return { success: true, data: null };
  } catch (error) {
    return normalizeSkillError(error);
  }
}

function normalizeSkillError(error: unknown): ActionResult<never> {
  if (isApiError(error)) {
    const apiError = error as ApiError;
    return { success: false, error: apiError.message };
  }
  return { success: false, error: "Operasi skill gagal diproses." };
}
