"use server";

export async function resolveServerAuthToken(): Promise<string | null> {
  const { getAccessToken } = await import("./auth");
  return getAccessToken();
}
