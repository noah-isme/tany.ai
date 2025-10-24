import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

import { ACCESS_TOKEN_COOKIE } from "@/lib/auth";
import { decodeJwt } from "@/lib/jwt";

const ADMIN_PREFIX = "/admin";

function tokenExpired(exp?: number): boolean {
  if (!exp) {
    return false;
  }
  const nowSeconds = Math.floor(Date.now() / 1000);
  return exp <= nowSeconds;
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  const token = request.cookies.get(ACCESS_TOKEN_COOKIE)?.value ?? "";
  const payload = token ? decodeJwt(token) : null;
  const isAdmin = Boolean(payload?.roles?.includes("admin")) && !tokenExpired(payload?.exp);

  if (pathname.startsWith(ADMIN_PREFIX)) {
    if (!token || !payload || !isAdmin) {
      const redirectUrl = new URL("/login", request.url);
      const response = NextResponse.redirect(redirectUrl);
      response.cookies.delete(ACCESS_TOKEN_COOKIE);
      return response;
    }
    return NextResponse.next();
  }

  if (pathname === "/login" && isAdmin) {
    return NextResponse.redirect(new URL(ADMIN_PREFIX, request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/admin/:path*", "/login"],
};
