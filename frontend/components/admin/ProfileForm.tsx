"use client";

import { useState, useTransition } from "react";
import { Controller, useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import type { ActionResult } from "@/lib/action-result";
import type { Profile } from "@/lib/types/admin";
import { profileSchema, type ProfileFormValues } from "@/lib/validators";

import { ImageUploader } from "./ImageUploader";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Textarea } from "./ui/textarea";

const defaultStatus = null as { type: "success" | "error"; message: string } | null;

type ProfileFormProps = {
  profile: Profile;
  onSubmit: (values: ProfileFormValues) => Promise<ActionResult<Profile>>;
};

export function ProfileForm({ profile, onSubmit }: ProfileFormProps) {
  const [status, setStatus] = useState(defaultStatus);
  const [isPending, startTransition] = useTransition();

  const {
    register,
    handleSubmit,
    setError,
    control,
    reset,
    formState: { errors, isDirty },
  } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      name: profile.name,
      title: profile.title,
      bio: profile.bio ?? "",
      email: profile.email ?? "",
      phone: profile.phone ?? "",
      location: profile.location ?? "",
      avatar_url: profile.avatar_url ?? "",
    },
  });

  const avatarPreview = useWatch({ control, name: "avatar_url" });

  const submitHandler = (values: ProfileFormValues) => {
    setStatus(null);
    startTransition(async () => {
      const result = await onSubmit(values);
      if (result.success) {
        setStatus({ type: "success", message: result.message ?? "Profil tersimpan." });
        reset({ ...values });
        return;
      }
      setStatus({ type: "error", message: result.error });
      if (result.fieldErrors) {
        Object.entries(result.fieldErrors).forEach(([field, message]) => {
          setError(field as keyof ProfileFormValues, { type: "server", message });
        });
      }
    });
  };

  return (
    <form onSubmit={handleSubmit(submitHandler)} className="grid gap-6 lg:grid-cols-[1fr_minmax(280px,340px)]">
      <div className="space-y-4">
        <Field label="Nama" error={errors.name?.message}>
          <Input
            placeholder="Nama lengkap"
            {...register("name")}
            invalid={Boolean(errors.name)}
          />
        </Field>
        <Field label="Jabatan" error={errors.title?.message}>
          <Input
            placeholder="Contoh: Product Designer"
            {...register("title")}
            invalid={Boolean(errors.title)}
          />
        </Field>
        <Field label="Bio" error={errors.bio?.message}>
          <Textarea
            placeholder="Ringkasan tentang Anda"
            {...register("bio")}
            invalid={Boolean(errors.bio)}
          />
        </Field>
        <div className="grid gap-4 sm:grid-cols-2">
          <Field label="Email" error={errors.email?.message}>
            <Input
              type="email"
              placeholder="admin@example.com"
              {...register("email")}
              invalid={Boolean(errors.email)}
            />
          </Field>
          <Field label="Telepon" error={errors.phone?.message}>
            <Input
              placeholder="Nomor kontak"
              {...register("phone")}
              invalid={Boolean(errors.phone)}
            />
          </Field>
        </div>
        <Field label="Lokasi" error={errors.location?.message}>
          <Input
            placeholder="Kota, Negara"
            {...register("location")}
            invalid={Boolean(errors.location)}
          />
        </Field>
        <Field label="Avatar" error={errors.avatar_url?.message}>
          <Controller
            name="avatar_url"
            control={control}
            render={({ field }) => (
              <ImageUploader
                value={field.value ?? ""}
                onChange={(url) => field.onChange(url)}
                onBlur={field.onBlur}
                disabled={isPending}
              />
            )}
          />
        </Field>

        {status ? (
          <p
            className={`rounded-md border px-3 py-2 text-sm ${
              status.type === "success"
                ? "border-emerald-400/60 bg-emerald-500/10 text-emerald-300"
                : "border-rose-400/60 bg-rose-500/10 text-rose-200"
            }`}
            role="alert"
          >
            {status.message}
          </p>
        ) : null}

        <div className="flex justify-end">
          <Button type="submit" disabled={!isDirty || isPending}>
            {isPending ? "Menyimpan…" : "Simpan perubahan"}
          </Button>
        </div>
      </div>

      <aside className="space-y-4">
        <div className="rounded-2xl border border-border bg-card/90 p-4 text-sm text-muted-foreground shadow-sm transition-colors supports-[backdrop-filter]:bg-card/70 supports-[backdrop-filter]:backdrop-blur">
          <p className="font-semibold text-foreground">Pratinjau Avatar</p>
          <div className="mt-3 flex items-center justify-center">
            {avatarPreview ? (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                src={avatarPreview}
                alt="Avatar preview"
                className="h-32 w-32 rounded-full border border-border object-cover shadow-sm"
              />
            ) : (
              <div className="flex h-32 w-32 items-center justify-center rounded-full border border-dashed border-border text-xs text-muted-foreground">
                Tidak ada gambar
              </div>
            )}
          </div>
          <p className="mt-3 text-xs text-muted-foreground">
            Unggahan akan tersimpan otomatis ke storage dan menghasilkan URL publik.
          </p>
        </div>
      </aside>
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
    <label className="space-y-2 text-sm font-medium text-foreground">
      <span>{label}</span>
      <div>{children}</div>
      {error ? (
        <span className="block text-xs font-normal text-destructive" role="alert">
          {error}
        </span>
      ) : null}
    </label>
  );
}
