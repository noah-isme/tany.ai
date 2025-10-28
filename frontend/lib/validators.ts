import { z } from "zod";

const optionalString = (schema: z.ZodString) =>
  z.union([schema, z.literal("")]).transform((value) => (value === "" ? "" : value));

const priceField = (message: string) =>
  z
    .union([z.string(), z.number(), z.null()])
    .optional()
    .transform((value) => {
      if (value === null || value === undefined || value === "") {
        return null;
      }
      if (typeof value === "number") {
        return value;
      }
      const parsed = Number(value);
      return parsed;
    })
    .refine(
      (value) =>
        value === null || (typeof value === "number" && Number.isFinite(value) && value >= 0),
      message,
    );

export const loginSchema = z.object({
  email: z.string({ required_error: "Email wajib diisi" }).min(1, "Email wajib diisi").email("Email tidak valid"),
  password: z
    .string({ required_error: "Password wajib diisi" })
    .min(8, "Password minimal 8 karakter"),
});

export const profileSchema = z.object({
  name: z.string().min(2, "Nama minimal 2 karakter").max(120, "Nama maksimal 120 karakter"),
  title: z.string().min(2, "Judul minimal 2 karakter").max(160, "Judul maksimal 160 karakter"),
  bio: optionalString(z.string().max(2000, "Bio maksimal 2000 karakter")),
  email: optionalString(z.string().email("Format email tidak valid")),
  phone: optionalString(z.string().max(64, "Telepon maksimal 64 karakter")),
  location: optionalString(z.string().max(160, "Lokasi maksimal 160 karakter")),
  avatar_url: optionalString(z.string().url("URL avatar tidak valid")),
});

export const skillSchema = z.object({
  name: z.string().min(2, "Nama skill minimal 2 karakter").max(80, "Nama skill maksimal 80 karakter"),
});

export const serviceSchema = z
  .object({
    name: z.string().min(2, "Nama layanan minimal 2 karakter").max(120),
    description: optionalString(z.string().max(2000, "Deskripsi maksimal 2000 karakter")),
    price_min: priceField("Harga minimal tidak valid"),
    price_max: priceField("Harga maksimal tidak valid"),
    currency: optionalString(z.string().max(3, "Maksimal 3 karakter")),
    duration_label: optionalString(z.string().max(80, "Durasi maksimal 80 karakter")),
    is_active: z.boolean().optional().default(true),
  })
  .refine(
    (data) => {
      if (data.price_min !== null && data.price_max !== null) {
        return data.price_max >= data.price_min;
      }
      return true;
    },
    { message: "Harga maksimal harus >= harga minimal", path: ["price_max"] },
  );

export const projectSchema = z.object({
  title: z.string().min(2, "Judul minimal 2 karakter").max(160),
  description: optionalString(z.string().max(4000, "Deskripsi maksimal 4000 karakter")),
  tech_stack: z
    .array(optionalString(z.string().max(32, "Tech stack maksimal 32 karakter")))
    .transform((items) => items.map((item) => item.trim()).filter((item) => item !== "")),
  image_url: optionalString(z.string().url("URL gambar tidak valid")),
  project_url: optionalString(z.string().url("URL proyek tidak valid")),
  category: optionalString(z.string().max(80, "Kategori maksimal 80 karakter")),
  duration_label: optionalString(z.string().max(80, "Durasi maksimal 80 karakter")),
  price_label: optionalString(z.string().max(120, "Label harga maksimal 120 karakter")),
  budget_label: optionalString(z.string().max(120, "Label budget maksimal 120 karakter")),
  is_featured: z.boolean().optional().default(false),
});

export const apiKeySchema = z.object({
  openai: optionalString(z.string()),
  anthropic: optionalString(z.string()),
  pinecone: optionalString(z.string()),
});

export type LoginFormValues = z.infer<typeof loginSchema>;
export type ProfileFormValues = z.infer<typeof profileSchema>;
export type SkillFormValues = z.infer<typeof skillSchema>;
export type ServiceFormValues = z.infer<typeof serviceSchema>;
export type ProjectFormValues = z.infer<typeof projectSchema>;
export type ApiKeyFormValues = z.infer<typeof apiKeySchema>;
