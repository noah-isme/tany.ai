export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? 
  process.env.API_BASE_URL ?? 
  process.env.NEXT_PUBLIC_API_BASE_URL ?? 
  "http://localhost:8080";

export function resolveApiUrl(path: string): string {
  if (!path.startsWith("/")) {
    throw new Error(`API path must start with / (received ${path})`);
  }
  return `${API_BASE_URL}${path}`;
}
