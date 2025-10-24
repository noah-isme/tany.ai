import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { decodeJwt, type JwtPayload } from "./jwt";
import { verifyAccessToken } from "./jwt-verify";

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
  const allowInsecureCookie = process.env.ALLOW_INSECURE_AUTH_COOKIE === "true";
  const cookieStore = await cookies();
  cookieStore.set({
    name: ACCESS_TOKEN_COOKIE,
    value: token,
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production" && !allowInsecureCookie,
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

function normalizeUserFromPayload(payload: JwtPayload): CurrentUser {
  const roles = Array.isArray(payload.roles)
    ? payload.roles.filter((role): role is string => typeof role === "string")
    : [];
  return {
    id: typeof payload.sub === "string" ? payload.sub : "",
    email: typeof payload.email === "string" ? payload.email : "",
    roles,
    expiresAt: typeof payload.exp === "number" ? payload.exp : undefined,
  };
}

export async function getCurrentUser(): Promise<CurrentUser | null> {
  const token = await getAccessToken();
  if (!token) {
    return null;
  }
  const payload = await verifyAccessToken(token);
  if (!payload) {
    return null;
  }
  return normalizeUserFromPayload(payload);
}

export async function ensureAdminUser(): Promise<CurrentUser> {
  const user = await getCurrentUser();
  if (!user) {
    throw new MissingAccessTokenError();
  }
  if (!user.roles.map((role) => role.toLowerCase()).includes("admin")) {
    throw new ForbiddenError();
  }
  return user;
}

export async function requireAdminOrRedirect(): Promise<CurrentUser> {
  const user = await getCurrentUser();
  if (!user) {
    redirect("/login");
  }
  if (!user.roles.map((role) => role.toLowerCase()).includes("admin")) {
    redirect("/403");
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
