"use client";

import { useState, useTransition } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { saveApiKeysAction } from "@/app/admin/settings/actions";
import type { ApiKeyFormValues } from "@/lib/validators";
import { apiKeySchema } from "@/lib/validators";

import { Button } from "./ui/button";
import { Input } from "./ui/input";

export type StoredApiKeySnapshot = {
  openai?: string;
  anthropic?: string;
  pinecone?: string;
  updatedAt?: string;
};

type ApiKeysFormProps = {
  initial?: StoredApiKeySnapshot;
};

export function ApiKeysForm({ initial }: ApiKeysFormProps) {
  const [status, setStatus] = useState<{ type: "success" | "error"; message: string } | null>(null);
  const [stored, setStored] = useState<StoredApiKeySnapshot | undefined>(initial);
  const [isPending, startTransition] = useTransition();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
    setError,
  } = useForm<ApiKeyFormValues>({
    resolver: zodResolver(apiKeySchema),
    defaultValues: { openai: "", anthropic: "", pinecone: "" },
  });

  const onSubmit = (values: ApiKeyFormValues) => {
    setStatus(null);
    startTransition(async () => {
      const result = await saveApiKeysAction(values);
      if (result.success) {
        setStored(result.data);
        setStatus({ type: "success", message: "API keys tersimpan secara aman." });
        reset({ openai: "", anthropic: "", pinecone: "" });
      } else {
        setStatus({ type: "error", message: result.error });
        if (result.fieldErrors) {
          Object.entries(result.fieldErrors).forEach(([field, message]) => {
            setError(field as keyof ApiKeyFormValues, { type: "server", message });
          });
        }
      }
    });
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <Field label="OpenAI API Key" error={errors.openai?.message}>
        <Input
          placeholder="Masukkan key baru"
          autoComplete="off"
          {...register("openai")}
          invalid={Boolean(errors.openai)}
        />
      </Field>
      <Field label="Anthropic Key" error={errors.anthropic?.message}>
        <Input
          placeholder="Opsional"
          autoComplete="off"
          {...register("anthropic")}
          invalid={Boolean(errors.anthropic)}
        />
      </Field>
      <Field label="Pinecone Key" error={errors.pinecone?.message}>
        <Input
          placeholder="Opsional"
          autoComplete="off"
          {...register("pinecone")}
          invalid={Boolean(errors.pinecone)}
        />
      </Field>

      {status ? (
        <p
          className={`rounded-md border px-3 py-2 text-sm ${
            status.type === "success"
              ? "border-emerald-400/60 bg-emerald-500/10 text-emerald-300"
              : "border-rose-400/60 bg-rose-500/10 text-rose-200"
          }`}
        >
          {status.message}
        </p>
      ) : null}

      <Button type="submit" disabled={isPending}>
        Simpan API keys
      </Button>

      <StoredKeysSummary stored={stored} />
    </form>
  );
}

type FieldProps = {
  label: string;
  error?: string;
  children: React.ReactNode;
};

function Field({ label, error, children }: FieldProps) {
  return (
    <label className="flex flex-col gap-2 text-sm font-medium text-slate-700 dark:text-slate-200">
      <span>{label}</span>
      <div>{children}</div>
      {error ? (
        <span className="text-xs font-normal text-rose-300" role="alert">
          {error}
        </span>
      ) : null}
    </label>
  );
}

type StoredKeysSummaryProps = {
  stored?: StoredApiKeySnapshot;
};

function StoredKeysSummary({ stored }: StoredKeysSummaryProps) {
  if (!stored) {
    return (
      <p className="text-xs text-slate-500 dark:text-slate-400">
        Belum ada kredensial yang disimpan. Nilai akan di-hash dan tidak ditampilkan ulang setelah tersimpan.
      </p>
    );
  }

  const badge = (value?: string) => (value ? `${value.slice(0, 6)}…${value.slice(-4)}` : "-" );

  return (
    <div className="rounded-2xl border border-slate-200 bg-white/80 p-4 text-xs text-slate-500 shadow-sm dark:border-slate-800 dark:bg-slate-900/70 dark:text-slate-400">
      <p className="font-semibold text-slate-900 dark:text-slate-100">Status penyimpanan</p>
      <ul className="mt-3 space-y-1">
        <li>OpenAI: {badge(stored.openai)}</li>
        <li>Anthropic: {badge(stored.anthropic)}</li>
        <li>Pinecone: {badge(stored.pinecone)}</li>
      </ul>
      {stored.updatedAt ? <p className="mt-2">Terakhir diperbarui: {new Date(stored.updatedAt).toLocaleString("id-ID")}</p> : null}
    </div>
  );
}
