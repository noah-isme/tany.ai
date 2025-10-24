"use server";

import { revalidatePath } from "next/cache";

import { updateProfile } from "@/lib/admin-api";
import { ApiError, isApiError } from "@/lib/api-client";
import { type ActionResult } from "@/lib/action-result";
import { issuesToFieldErrors } from "@/lib/form-errors";
import type { Profile } from "@/lib/types/admin";
import { profileSchema, type ProfileFormValues } from "@/lib/validators";

export async function updateProfileAction(values: ProfileFormValues): Promise<ActionResult<Profile>> {
  const parsed = profileSchema.safeParse(values);
  if (!parsed.success) {
    return {
      success: false,
      error: "Periksa kembali data yang diisi.",
      fieldErrors: issuesToFieldErrors(parsed.error.issues),
    };
  }

  const payload = {
    name: parsed.data.name,
    title: parsed.data.title,
    bio: parsed.data.bio ?? "",
    email: parsed.data.email ?? "",
    phone: parsed.data.phone ?? "",
    location: parsed.data.location ?? "",
    avatar_url: parsed.data.avatar_url ?? "",
  };

  try {
    const updated = await updateProfile(payload);
    revalidatePath("/admin/profile");
    revalidatePath("/admin");
    return { success: true, data: updated, message: "Profil berhasil diperbarui" };
  } catch (error) {
    if (isApiError(error)) {
      const apiError = error as ApiError;
      const fieldErrors =
        apiError.details && typeof apiError.details === "object"
          ? Object.fromEntries(
              Object.entries(apiError.details as Record<string, unknown>).filter(
                ([, value]) => typeof value === "string",
              ) as [string, string][]
            )
          : undefined;
      return {
        success: false,
        error: apiError.message,
        fieldErrors,
      };
    }
    return { success: false, error: "Terjadi kesalahan saat memperbarui profil." };
  }
}
