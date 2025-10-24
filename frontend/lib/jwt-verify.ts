import { jwtVerify } from "jose";

import { decodeJwt, type JwtPayload } from "./jwt";

let cachedSecret: CryptoKey | Uint8Array | null = null;

async function getJwtSecret(): Promise<CryptoKey | Uint8Array> {
  if (cachedSecret) {
    return cachedSecret;
  }
  const secret = process.env.JWT_SECRET;
  if (!secret) {
    throw new Error("JWT_SECRET environment variable is not defined");
  }
  const bytes = new TextEncoder().encode(secret);
  const cryptoModule: Crypto | undefined = typeof globalThis !== "undefined" ? (globalThis as { crypto?: Crypto }).crypto : undefined;
  if (cryptoModule?.subtle && typeof cryptoModule.subtle.importKey === "function") {
    cachedSecret = await cryptoModule.subtle.importKey(
      "raw",
      bytes,
      { name: "HMAC", hash: "SHA-256" },
      false,
      ["sign", "verify"],
    );
  } else {
    cachedSecret = bytes;
  }
  return cachedSecret;
}

export async function verifyAccessToken(token: string): Promise<JwtPayload | null> {
  if (!token) {
    return null;
  }
  const decoded = decodeJwt(token);
  if (!decoded) {
    return null;
  }
  try {
    const secret = await getJwtSecret();
    const { payload } = await jwtVerify(token, secret);
    const exp = typeof payload?.exp === "number" ? payload.exp : decoded.exp;
    return { ...decoded, exp } satisfies JwtPayload;
  } catch {
    return null;
  }
}

export function resetJwtSecretCacheForTests(): void {
  cachedSecret = null;
}
