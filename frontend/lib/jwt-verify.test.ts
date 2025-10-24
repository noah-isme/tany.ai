import { createHmac } from "node:crypto";
import { describe, expect, it, beforeEach } from "vitest";

import { verifyAccessToken, resetJwtSecretCacheForTests } from "./jwt-verify";

const SECRET = "unit-test-secret";

function signToken(payload: Record<string, unknown>): string {
  const header = Buffer.from(JSON.stringify({ alg: "HS256", typ: "JWT" })).toString("base64url");
  const body = Buffer.from(JSON.stringify(payload)).toString("base64url");
  const signature = createHmac("sha256", SECRET).update(`${header}.${body}`).digest("base64url");
  return `${header}.${body}.${signature}`;
}

describe("verifyAccessToken", () => {
  beforeEach(() => {
    process.env.JWT_SECRET = SECRET;
    resetJwtSecretCacheForTests();
  });

  it("returns payload for valid token", async () => {
    const token = signToken({ sub: "1", email: "user@example.com", roles: ["admin"], exp: Math.floor(Date.now() / 1000) + 60 });
    const payload = await verifyAccessToken(token);
    expect(payload?.sub).toBe("1");
    expect(payload?.email).toBe("user@example.com");
    expect(payload?.roles).toContain("admin");
  });

  it("returns null for invalid token", async () => {
    const token = "invalid.token.value";
    const payload = await verifyAccessToken(token);
    expect(payload).toBeNull();
  });
});
