import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { decodeJwt, JwtPayload, tokenHasRole } from "./jwt";

export const ACCESS_TOKEN_COOKIE = "ta_access";
export const THEME_COOKIE = "ta_theme";

export type CurrentUser = {
  id: string;
  email: string;
  roles: string[];
  expiresAt?: number;
};

export class MissingAccessTokenError extends Error {
  constructor(message = "missing access token") {
    super(message);
    this.name = "MissingAccessTokenError";
  }
}

export class ForbiddenError extends Error {
  constructor(message = "forbidden") {
    super(message);
    this.name = "ForbiddenError";
  }
}

function computeMaxAge(payload: JwtPayload | null): number | undefined {
  if (!payload?.exp) {
    return undefined;
  }
  const nowSeconds = Math.floor(Date.now() / 1000);
  const remaining = payload.exp - nowSeconds;
  if (remaining <= 0) {
    return undefined;
  }
  return remaining;
}

export async function setAccessTokenCookie(token: string): Promise<void> {
  const payload = decodeJwt(token);
  const maxAge = computeMaxAge(payload);
  const cookieStore = await cookies();
  cookieStore.set({
    name: ACCESS_TOKEN_COOKIE,
    value: token,
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge,
  });
}

export async function clearAccessTokenCookie(): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.delete(ACCESS_TOKEN_COOKIE);
}

export async function getAccessToken(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get(ACCESS_TOKEN_COOKIE)?.value ?? null;
}

export async function getCurrentUser(): Promise<CurrentUser | null> {
  const token = await getAccessToken();
  if (!token) {
    return null;
  }
  const payload = decodeJwt(token);
  if (!payload) {
    return null;
  }
  return {
    id: typeof payload.sub === "string" ? payload.sub : "",
    email: typeof payload.email === "string" ? payload.email : "",
    roles: Array.isArray(payload.roles) ? payload.roles : [],
    expiresAt: typeof payload.exp === "number" ? payload.exp : undefined,
  };
}

export async function ensureAdminUser(): Promise<CurrentUser> {
  const user = await getCurrentUser();
  if (!user) {
    throw new MissingAccessTokenError();
  }
  const token = await getAccessToken();
  if (!token || !tokenHasRole(token, "admin")) {
    throw new ForbiddenError();
  }
  return user;
}

export async function requireAdminOrRedirect(): Promise<CurrentUser> {
  const user = await getCurrentUser();
  if (!user) {
    redirect("/login");
  }
  const token = await getAccessToken();
  if (!token || !tokenHasRole(token, "admin")) {
    redirect("/login");
  }
  return user;
}

export async function getStoredTheme(): Promise<"light" | "dark" | null> {
  const cookieStore = await cookies();
  const theme = cookieStore.get(THEME_COOKIE)?.value;
  if (theme === "light" || theme === "dark") {
    return theme;
  }
  return null;
}

export async function persistThemePreference(theme: "light" | "dark"): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.set({
    name: THEME_COOKIE,
    value: theme,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: 60 * 60 * 24 * 365,
  });
}
