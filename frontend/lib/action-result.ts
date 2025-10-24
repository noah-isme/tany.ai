export type ActionResult<T> =
  | { success: true; data: T; message?: string }
  | { success: false; error: string; fieldErrors?: Record<string, string> };

export function createFieldErrors(error: unknown): Record<string, string> | undefined {
  if (!error || typeof error !== "object") {
    return undefined;
  }
  if (error instanceof Map) {
    const result: Record<string, string> = {};
    error.forEach((value, key) => {
      result[String(key)] = String(value);
    });
    return result;
  }
  if ("details" in error && typeof error.details === "object" && error.details) {
    const details = error.details as Record<string, unknown>;
    const result: Record<string, string> = {};
    for (const [field, value] of Object.entries(details)) {
      if (typeof value === "string") {
        result[field] = value;
      }
    }
    return Object.keys(result).length > 0 ? result : undefined;
  }
  return undefined;
}
