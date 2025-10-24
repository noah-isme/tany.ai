import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

import { ACCESS_TOKEN_COOKIE } from "@/lib/auth";
import { verifyAccessToken } from "@/lib/jwt-verify";

const ADMIN_PREFIX = "/admin";

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  const token = request.cookies.get(ACCESS_TOKEN_COOKIE)?.value ?? "";

  if (pathname.startsWith(ADMIN_PREFIX)) {
    return handleAdminRoute(request, token);
  }

  if (pathname === "/login") {
    return handleLoginRoute(request, token);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/admin/:path*", "/login"],
};

function createLoginRedirect(request: NextRequest): NextResponse {
  const redirectUrl = new URL("/login", request.url);
  const response = NextResponse.redirect(redirectUrl);
  response.cookies.delete(ACCESS_TOKEN_COOKIE);
  return response;
}

async function handleAdminRoute(request: NextRequest, token: string): Promise<NextResponse> {
  if (!token) {
    return createLoginRedirect(request);
  }

  const payload = await verifyAccessToken(token);
  if (!payload) {
    return createLoginRedirect(request);
  }

  const roles = Array.isArray(payload.roles) ? payload.roles.map((role) => String(role).toLowerCase()) : [];
  if (!roles.includes("admin")) {
    const forbidden = NextResponse.rewrite(new URL("/403", request.url), { status: 403 });
    forbidden.headers.set("x-middleware-cache", "no-cache");
    return forbidden;
  }

  return NextResponse.next();
}

async function handleLoginRoute(request: NextRequest, token: string): Promise<NextResponse> {
  if (!token) {
    return NextResponse.next();
  }
  const payload = await verifyAccessToken(token);
  if (!payload) {
    const response = NextResponse.next();
    response.cookies.delete(ACCESS_TOKEN_COOKIE);
    return response;
  }
  const roles = Array.isArray(payload?.roles) ? payload.roles.map((role) => String(role).toLowerCase()) : [];
  if (roles.includes("admin")) {
    return NextResponse.redirect(new URL(ADMIN_PREFIX, request.url));
  }
  return NextResponse.next();
}
