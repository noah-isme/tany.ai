import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

import { ACCESS_TOKEN_COOKIE } from "@/lib/auth";
import { verifyAccessToken } from "@/lib/jwt-verify";

const ADMIN_PREFIX = "/admin";

export async function middleware(request: NextRequest) {
  // Generate nonce untuk CSP
  const nonce = Buffer.from(crypto.randomUUID()).toString("base64");
  
  // HTTPS redirect di production
  if (
    process.env.NODE_ENV === "production" &&
    request.headers.get("x-forwarded-proto") !== "https"
  ) {
    const url = request.nextUrl.clone();
    url.protocol = "https:";
    return NextResponse.redirect(url);
  }

  const { pathname } = request.nextUrl;

  const token = request.cookies.get(ACCESS_TOKEN_COOKIE)?.value ?? "";

  let response: NextResponse;

  if (pathname.startsWith(ADMIN_PREFIX)) {
    response = await handleAdminRoute(request, token);
  } else if (pathname === "/login") {
    response = await handleLoginRoute(request, token);
  } else {
    response = NextResponse.next();
  }

  // Set security headers dengan CSP yang aman untuk Next.js
  const apiOrigin = process.env.NEXT_PUBLIC_API_URL || 
                   process.env.NEXT_PUBLIC_API_BASE_URL || 
                   "http://localhost:8080";
  
  const cspHeader = `
    default-src 'self';
    script-src 'self' 'nonce-${nonce}' 'strict-dynamic' ${process.env.NODE_ENV === "development" ? "'unsafe-eval'" : ""};
    style-src 'self' 'unsafe-inline';
    img-src 'self' data: blob: https:;
    font-src 'self' data:;
    connect-src 'self' ${apiOrigin} https://generativelanguage.googleapis.com https://*.apn.leapcell.dev;
    frame-ancestors 'none';
    base-uri 'self';
    form-action 'self';
    object-src 'none';
    ${process.env.NODE_ENV === "production" ? "upgrade-insecure-requests;" : ""}
  `.replace(/\n/g, "").replace(/\s{2,}/g, " ").trim();

  response.headers.set("Content-Security-Policy", cspHeader);
  response.headers.set("X-Content-Type-Options", "nosniff");
  response.headers.set("X-Frame-Options", "DENY");
  response.headers.set("X-XSS-Protection", "1; mode=block");
  response.headers.set("Referrer-Policy", "strict-origin-when-cross-origin");
  response.headers.set("Permissions-Policy", "geolocation=(), microphone=(), camera=()");
  response.headers.set("x-nonce", nonce);

  return response;
}

export const config = {
  matcher: ["/(.*)"],
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
