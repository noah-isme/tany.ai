"use server";

import { revalidatePath } from "next/cache";

import { persistThemePreference } from "@/lib/auth";
import { type ActionResult } from "@/lib/action-result";
import { apiKeySchema, type ApiKeyFormValues } from "@/lib/validators";
import { issuesToFieldErrors } from "@/lib/form-errors";

import crypto from "crypto";

type StoredApiKeys = {
  openai?: string;
  anthropic?: string;
  pinecone?: string;
  updatedAt?: string;
};

const globalStore = globalThis as unknown as { __tanyApiKeys?: StoredApiKeys };

export async function saveApiKeysAction(values: ApiKeyFormValues): Promise<ActionResult<StoredApiKeys>> {
  const parsed = apiKeySchema.safeParse(values);
  if (!parsed.success) {
    return {
      success: false,
      error: "Periksa kembali input API key.",
      fieldErrors: issuesToFieldErrors(parsed.error.issues),
    };
  }

  const now = new Date().toISOString();
  const hashed: StoredApiKeys = { updatedAt: now };
  if (parsed.data.openai) {
    hashed.openai = hashSecret(parsed.data.openai);
  }
  if (parsed.data.anthropic) {
    hashed.anthropic = hashSecret(parsed.data.anthropic);
  }
  if (parsed.data.pinecone) {
    hashed.pinecone = hashSecret(parsed.data.pinecone);
  }

  globalStore.__tanyApiKeys = {
    ...globalStore.__tanyApiKeys,
    ...hashed,
  };

  revalidatePath("/admin/settings");
  return { success: true, data: { ...globalStore.__tanyApiKeys } };
}

function hashSecret(value: string): string {
  return crypto.createHash("sha256").update(value).digest("hex");
}

type ThemeResult = { success: true };

export async function setThemePreferenceAction(theme: "light" | "dark" | "system"): Promise<ThemeResult> {
  await persistThemePreference(theme);
  revalidatePath("/admin/settings");
  return { success: true };
}

export async function getStoredApiKeys(): Promise<StoredApiKeys | undefined> {
  return globalStore.__tanyApiKeys;
}
