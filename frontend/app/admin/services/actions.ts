"use server";

import { revalidatePath } from "next/cache";

import {
  createService,
  deleteService,
  reorderServices,
  toggleService,
  updateService,
} from "@/lib/admin-api";
import { ApiError, isApiError } from "@/lib/api-client";
import { type ActionResult } from "@/lib/action-result";
import { issuesToFieldErrors } from "@/lib/form-errors";
import type { Service } from "@/lib/types/admin";
import { serviceSchema, type ServiceFormValues } from "@/lib/validators";

export async function createServiceAction(values: ServiceFormValues): Promise<ActionResult<Service>> {
  const parsed = serviceSchema.safeParse(values);
  if (!parsed.success) {
    return {
      success: false,
      error: "Periksa data layanan.",
      fieldErrors: issuesToFieldErrors(parsed.error.issues),
    };
  }

  try {
    const created = await createService(buildServicePayload(parsed.data));
    revalidatePath("/admin/services");
    revalidatePath("/admin");
    return { success: true, data: created, message: "Layanan dibuat." };
  } catch (error) {
    return normalizeServiceError(error);
  }
}

export async function updateServiceAction(params: {
  id: string;
  values: ServiceFormValues;
}): Promise<ActionResult<Service>> {
  const parsed = serviceSchema.safeParse(params.values);
  if (!parsed.success) {
    return {
      success: false,
      error: "Periksa data layanan.",
      fieldErrors: issuesToFieldErrors(parsed.error.issues),
    };
  }

  try {
    const updated = await updateService(params.id, buildServicePayload(parsed.data));
    revalidatePath("/admin/services");
    revalidatePath("/admin");
    return { success: true, data: updated, message: "Layanan diperbarui." };
  } catch (error) {
    return normalizeServiceError(error);
  }
}

export async function deleteServiceAction(params: { id: string }): Promise<ActionResult<null>> {
  try {
    await deleteService(params.id);
    revalidatePath("/admin/services");
    revalidatePath("/admin");
    return { success: true, data: null, message: "Layanan dihapus." };
  } catch (error) {
    return normalizeServiceError(error);
  }
}

export async function toggleServiceAction(params: {
  id: string;
  is_active: boolean;
}): Promise<ActionResult<Service>> {
  try {
    const toggled = await toggleService(params.id, params.is_active);
    revalidatePath("/admin/services");
    return { success: true, data: toggled };
  } catch (error) {
    return normalizeServiceError(error);
  }
}

export async function reorderServiceAction(params: {
  items: { id: string; order: number }[];
}): Promise<ActionResult<null>> {
  try {
    await reorderServices(params.items);
    revalidatePath("/admin/services");
    revalidatePath("/admin");
    return { success: true, data: null };
  } catch (error) {
    return normalizeServiceError(error);
  }
}

function buildServicePayload(values: ServiceFormValues) {
  return {
    name: values.name,
    description: values.description ?? "",
    price_min: values.price_min ?? null,
    price_max: values.price_max ?? null,
    currency: values.currency?.toUpperCase() ?? "",
    duration_label: values.duration_label ?? "",
    is_active: values.is_active ?? true,
  };
}

function normalizeServiceError(error: unknown): ActionResult<never> {
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
  return { success: false, error: "Operasi layanan gagal." };
}
