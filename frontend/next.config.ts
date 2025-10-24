import type { NextConfig } from "next";

const dev = process.env.NODE_ENV !== "production";
const connectSrc = [
  "'self'",
  "https:",
  "wss:",
  "http://localhost:4000",
  "http://127.0.0.1:4000",
];

if (dev) {
  connectSrc.push("ws://localhost:3000", "ws://127.0.0.1:3000");
}

const scriptSrc = ["'self'"];
if (dev) {
  scriptSrc.push("'unsafe-eval'", "'unsafe-inline'");
}

const cspDirectives = [
  "default-src 'self'",
  "img-src 'self' data:",
  `script-src ${scriptSrc.join(" ")}`,
  "style-src 'self' 'unsafe-inline'",
  `connect-src ${connectSrc.join(" ")}`,
];

const securityHeaders = [
  { key: "X-Frame-Options", value: "DENY" },
  { key: "X-Content-Type-Options", value: "nosniff" },
  { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
  {
    key: "Strict-Transport-Security",
    value: "max-age=63072000; includeSubDomains; preload",
  },
  {
    key: "Content-Security-Policy",
    value: cspDirectives.join("; "),
  },
];

const nextConfig: NextConfig = {
  experimental: {
    turbopackUseSystemTlsCerts: true,
  },
  env: {
    NEXT_PUBLIC_API_BASE_URL:
      process.env.NEXT_PUBLIC_API_BASE_URL ?? process.env.API_BASE_URL ?? "http://localhost:8080",
  },
  allowedDevOrigins: [
    "http://localhost:3000",
    "http://127.0.0.1:3000",
    "http://localhost:4000",
    "http://127.0.0.1:4000",
  ],
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: securityHeaders,
      },
    ];
  },
};

export default nextConfig;
