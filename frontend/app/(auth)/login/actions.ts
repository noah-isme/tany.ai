"use server";

import { loginRequest } from "@/lib/auth-api";
import { ApiError, isApiError } from "@/lib/api-client";
import { clearAccessTokenCookie, setAccessTokenCookie } from "@/lib/auth";
import { loginSchema } from "@/lib/validators";

export type LoginActionResult =
  | { success: true }
  | { success: false; error: string; fieldErrors?: Record<string, string> };

export async function loginAction(formData: FormData): Promise<LoginActionResult> {
  const email = (formData.get("email") ?? "") as string;
  const password = (formData.get("password") ?? "") as string;

  const parsed = loginSchema.safeParse({ email, password });
  if (!parsed.success) {
    const fieldErrors: Record<string, string> = {};
    for (const issue of parsed.error.issues) {
      if (issue.path.length > 0) {
        fieldErrors[issue.path[0] as string] = issue.message;
      }
    }
    return { success: false, error: "Periksa kembali data yang diisi.", fieldErrors };
  }

  try {
    const response = await loginRequest(parsed.data.email, parsed.data.password);
    await setAccessTokenCookie(response.accessToken);
    return { success: true };
  } catch (error) {
    await clearAccessTokenCookie();
    if (isApiError(error)) {
      const apiError = error as ApiError;
      const message = apiError.code === "UNAUTHORIZED" ? "Email atau password salah." : apiError.message;
      return { success: false, error: message };
    }
    return { success: false, error: "Terjadi kesalahan tak terduga." };
  }
}
