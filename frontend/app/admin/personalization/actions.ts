"use server";

import { revalidatePath } from "next/cache";

import {
  reindexPersonalization,
  resetPersonalization,
  updatePersonalizationWeight,
} from "@/lib/admin-api";
import { ApiError, isApiError } from "@/lib/api-client";
import { type ActionResult } from "@/lib/action-result";

export async function reindexPersonalizationAction(): Promise<ActionResult<{ indexed: number }>> {
  try {
    const result = await reindexPersonalization();
    revalidatePath("/admin/personalization");
    revalidatePath("/admin");
    return { success: true, data: result, message: "Embedding diperbarui." };
  } catch (error) {
    return normalizePersonalizationError(error, "Gagal melakukan reindex personalisasi.");
  }
}

export async function resetPersonalizationAction(): Promise<ActionResult<null>> {
  try {
    await resetPersonalization();
    revalidatePath("/admin/personalization");
    return { success: true, data: null, message: "Embedding direset." };
  } catch (error) {
    return normalizePersonalizationError(error, "Gagal mereset embedding.");
  }
}

export async function updatePersonalizationWeightAction(weight: number): Promise<ActionResult<number>> {
  try {
    const updated = await updatePersonalizationWeight(Number(weight.toFixed(2)));
    revalidatePath("/admin/personalization");
    return { success: true, data: updated, message: "Bobot personalisasi diperbarui." };
  } catch (error) {
    return normalizePersonalizationError(error, "Gagal memperbarui bobot personalisasi.");
  }
}

function normalizePersonalizationError(error: unknown, fallback: string): ActionResult<never> {
  if (isApiError(error)) {
    const apiError = error as ApiError;
    return { success: false, error: apiError.message ?? fallback };
  }
  return { success: false, error: fallback };
}
