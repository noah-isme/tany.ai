"use server";

import { logoutRequest } from "@/lib/auth-api";
import { clearAccessTokenCookie } from "@/lib/auth";

export type LogoutResult = { success: true } | { success: false; error: string };

export async function logoutAction(): Promise<LogoutResult> {
  try {
    await logoutRequest();
  } catch {
    await clearAccessTokenCookie();
    return { success: false, error: "Gagal logout dari server." };
  }
  await clearAccessTokenCookie();
  return { success: true };
}
