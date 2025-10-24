import { describe, expect, it } from "vitest";

import { decodeJwt, tokenHasRole } from "./jwt";

function createToken(payload: object): string {
  const header = Buffer.from(JSON.stringify({ alg: "HS256", typ: "JWT" })).toString("base64url");
  const body = Buffer.from(JSON.stringify(payload)).toString("base64url");
  return `${header}.${body}.signature`;
}

describe("decodeJwt", () => {
  it("decodes base64 payload", () => {
    const token = createToken({ sub: "123", email: "user@example.com" });
    const result = decodeJwt(token);
    expect(result?.sub).toBe("123");
    expect(result?.email).toBe("user@example.com");
  });

  it("returns null for malformed token", () => {
    expect(decodeJwt("invalid")).toBeNull();
  });
});

describe("tokenHasRole", () => {
  it("returns true when role exists", () => {
    const token = createToken({ roles: ["ADMIN", "editor"] });
    expect(tokenHasRole(token, "admin")).toBe(true);
    expect(tokenHasRole(token, "Editor")).toBe(true);
  });

  it("returns false when role absent", () => {
    const token = createToken({ roles: ["viewer"] });
    expect(tokenHasRole(token, "admin")).toBe(false);
  });
});
