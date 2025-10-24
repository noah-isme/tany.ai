import { resolveApiUrl } from "./env";

export type ApiRequestInit = {
  method?: string;
  body?: unknown;
  token?: string | null;
  headers?: HeadersInit;
  withAuth?: boolean;
  cache?: RequestCache;
};

export type ApiErrorBody = {
  code?: string;
  message?: string;
  details?: unknown;
};

export class ApiError extends Error {
  status: number;
  code?: string;
  details?: unknown;

  constructor(message: string, status: number, body?: ApiErrorBody) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = body?.code;
    this.details = body?.details;
  }
}

export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError;
}

export async function apiFetch<T = unknown>(path: string, init: ApiRequestInit = {}): Promise<T> {
  const {
    method = "GET",
    body,
    token,
    headers,
    withAuth = true,
    cache = "no-store",
  } = init;

  const requestHeaders = new Headers(headers ?? {});
  requestHeaders.set("Accept", "application/json");

  let requestBody: BodyInit | undefined;
  if (body !== undefined) {
    if (body instanceof FormData || body instanceof URLSearchParams) {
      requestBody = body;
    } else {
      requestBody = JSON.stringify(body);
      if (!requestHeaders.has("Content-Type")) {
        requestHeaders.set("Content-Type", "application/json");
      }
    }
  }

  const bearer = await resolveAuthToken(token, withAuth);
  if (bearer) {
    requestHeaders.set("Authorization", `Bearer ${bearer}`);
  }

  const response = await fetch(resolveApiUrl(path), {
    method,
    headers: requestHeaders,
    body: requestBody,
    cache,
  });

  const contentType = response.headers.get("content-type") ?? "";
  const isJson = contentType.includes("application/json");
  const payload = isJson ? await response.json() : await response.text();

  if (!response.ok) {
    if (isJson && payload && typeof payload === "object" && "error" in payload) {
      const body = (payload as { error: ApiErrorBody }).error;
      throw new ApiError(body.message ?? "request failed", response.status, body);
    }
    const message = typeof payload === "string" && payload ? payload : "request failed";
    throw new ApiError(message, response.status);
  }

  return payload as T;
}

async function resolveAuthToken(
  token: string | null | undefined,
  withAuth: boolean,
): Promise<string | null> {
  if (token !== undefined) {
    return token;
  }
  if (!withAuth) {
    return null;
  }
  if (typeof window !== "undefined") {
    return null;
  }
  const { resolveServerAuthToken } = await import("./api-client-server");
  return resolveServerAuthToken();
}
