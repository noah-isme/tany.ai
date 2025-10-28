"use server";

import { revalidatePath } from "next/cache";

import { setExternalItemVisibility, syncExternalSource } from "@/lib/admin-api";
import { ApiError, isApiError } from "@/lib/api-client";
import { type ActionResult } from "@/lib/action-result";
import type { ExternalItem } from "@/lib/types/admin";

export async function syncExternalSourceAction(
  id: string,
): Promise<ActionResult<{ itemsUpserted: number; message?: string }>> {
  try {
    const result = await syncExternalSource(id);
    revalidatePath("/admin/integrations");
    revalidatePath("/admin");
    return { success: true, data: result, message: result.message ?? "Sinkronisasi selesai." };
  } catch (error) {
    return normalizeIntegrationError(error, "Gagal sinkronisasi sumber.");
  }
}

export async function toggleExternalItemVisibilityAction(params: {
  id: string;
  visible: boolean;
}): Promise<ActionResult<ExternalItem>> {
  try {
    const updated = await setExternalItemVisibility(params.id, params.visible);
    revalidatePath("/admin/integrations");
    return { success: true, data: updated };
  } catch (error) {
    return normalizeIntegrationError(error, "Gagal memperbarui visibilitas konten.");
  }
}

function normalizeIntegrationError(error: unknown, fallback: string): ActionResult<never> {
  if (isApiError(error)) {
    const apiError = error as ApiError;
    return { success: false, error: apiError.message || fallback };
  }
  return { success: false, error: fallback };
}
