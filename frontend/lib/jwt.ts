export type JwtPayload = {
  sub?: string;
  email?: string;
  roles?: string[];
  exp?: number;
  iat?: number;
  [key: string]: unknown;
};

function decodeBase64Url(segment: string): string {
  const normalized = segment.replace(/-/g, "+").replace(/_/g, "/");
  const padded = normalized.padEnd(normalized.length + ((4 - (normalized.length % 4)) % 4), "=");

  if (typeof atob === "function") {
    const binary = atob(padded);
    if (typeof TextDecoder !== "undefined") {
      const bytes = new Uint8Array(binary.length);
      for (let i = 0; i < binary.length; i += 1) {
        bytes[i] = binary.charCodeAt(i);
      }
      return new TextDecoder().decode(bytes);
    }
    return binary;
  }

  if (typeof Buffer !== "undefined") {
    return Buffer.from(padded, "base64").toString("utf-8");
  }
  throw new Error("No base64 decoder available in this environment");
}

export function decodeJwt(token: string): JwtPayload | null {
  try {
    const [, payloadSegment] = token.split(".");
    if (!payloadSegment) {
      return null;
    }
    const json = decodeBase64Url(payloadSegment);
    const payload = JSON.parse(json);
    if (payload && typeof payload === "object") {
      return payload as JwtPayload;
    }
    return null;
  } catch {
    return null;
  }
}

export function tokenHasRole(token: string, role: string): boolean {
  const payload = decodeJwt(token);
  if (!payload || !Array.isArray(payload.roles)) {
    return false;
  }
  return payload.roles.map((value) => value.toLowerCase()).includes(role.toLowerCase());
}
