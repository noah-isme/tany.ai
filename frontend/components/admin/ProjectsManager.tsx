"use client";

import { useEffect, useMemo, useState, useTransition } from "react";
import {
  DndContext,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  verticalListSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { zodResolver } from "@hookform/resolvers/zod";
import { clsx } from "clsx";
import { useForm, Controller, useFieldArray, useWatch, type FieldArrayPath } from "react-hook-form";
import { GripVertical, Pencil, Sparkles, Star, Trash2 } from "lucide-react";

import type { ActionResult } from "@/lib/action-result";
import type { Project } from "@/lib/types/admin";
import { projectSchema, type ProjectFormValues } from "@/lib/validators";

import { ImageUploader } from "./ImageUploader";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Switch } from "./ui/switch";
import { Textarea } from "./ui/textarea";

const emptyStatus = null as { type: "success" | "error"; message: string } | null;

type ProjectActions = {
  createProject: (values: ProjectFormValues) => Promise<ActionResult<Project>>;
  updateProject: (args: { id: string; values: ProjectFormValues }) => Promise<ActionResult<Project>>;
  deleteProject: (args: { id: string }) => Promise<ActionResult<null>>;
  reorderProject: (args: { items: { id: string; order: number }[] }) => Promise<ActionResult<null>>;
  featureProject: (args: { id: string; is_featured: boolean }) => Promise<ActionResult<Project>>;
};

type ProjectsManagerProps = ProjectActions & {
  initialProjects: Project[];
};

export function ProjectsManager({
  initialProjects,
  createProject,
  updateProject,
  deleteProject,
  reorderProject,
  featureProject,
}: ProjectsManagerProps) {
  const [projects, setProjects] = useState(initialProjects);
  const [status, setStatus] = useState(emptyStatus);
  const [editor, setEditor] = useState<{ mode: "create" | "edit"; project?: Project } | null>(null);
  const [isPending, startTransition] = useTransition();
  const [pendingId, setPendingId] = useState<string | null>(null);

  useEffect(() => {
    setProjects(initialProjects);
  }, [initialProjects]);

  const sensors = useSensors(useSensor(PointerSensor));
  const items = useMemo(() => projects.map((project) => project.id), [projects]);

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) {
      return;
    }
    setProjects((prev) => {
      const oldIndex = prev.findIndex((item) => item.id === active.id);
      const newIndex = prev.findIndex((item) => item.id === over.id);
      const reordered = arrayMove(prev, oldIndex, newIndex);
      startTransition(async () => {
        const payload = reordered.map((project, index) => ({ id: project.id, order: index }));
        const result = await reorderProject({ items: payload });
        if (!result.success) {
          setStatus({ type: "error", message: result.error });
        }
      });
      return reordered;
    });
  };

  const handleDelete = (id: string) => {
    if (typeof window !== "undefined" && !window.confirm("Hapus proyek ini?")) {
      return;
    }
    setPendingId(id);
    startTransition(async () => {
      const result = await deleteProject({ id });
      if (result.success) {
        setProjects((prev) => prev.filter((project) => project.id !== id));
        setStatus({ type: "success", message: "Proyek dihapus." });
      } else {
        setStatus({ type: "error", message: result.error });
      }
      setPendingId(null);
    });
  };

  const handleFeature = (project: Project, isFeatured: boolean) => {
    setPendingId(project.id);
    startTransition(async () => {
      const result = await featureProject({ id: project.id, is_featured: isFeatured });
      if (result.success) {
        setProjects((prev) =>
          prev.map((item) => (item.id === project.id ? { ...item, is_featured: result.data.is_featured } : item)),
        );
      } else {
        setStatus({ type: "error", message: result.error });
      }
      setPendingId(null);
    });
  };

  const handleCreate = (values: ProjectFormValues) => {
    startTransition(async () => {
      const result = await createProject(values);
      if (result.success) {
        setProjects((prev) => [...prev, result.data]);
        setStatus({ type: "success", message: "Proyek ditambahkan." });
        setEditor(null);
      } else {
        setStatus({ type: "error", message: result.error });
      }
    });
  };

  const handleUpdate = (project: Project, values: ProjectFormValues) => {
    startTransition(async () => {
      const result = await updateProject({ id: project.id, values });
      if (result.success) {
        setProjects((prev) =>
          prev.map((item) => (item.id === project.id ? { ...item, ...result.data } : item)),
        );
        setStatus({ type: "success", message: "Proyek diperbarui." });
        setEditor(null);
      } else {
        setStatus({ type: "error", message: result.error });
      }
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="space-y-1">
          <h2 className="text-base font-semibold text-foreground">Portfolio Projects</h2>
          <p className="text-xs text-muted-foreground">
            Susun proyek unggulan untuk dijadikan rujukan AI saat menjawab studi kasus.
          </p>
        </div>
        <Button type="button" onClick={() => setEditor({ mode: "create" })}>
          Tambah Proyek
        </Button>
      </div>

      {status ? (
        <p
          className={clsx(
            "rounded-md border px-3 py-2 text-sm",
            status.type === "success"
              ? "border-emerald-400/60 bg-emerald-500/10 text-emerald-300"
              : "border-rose-400/60 bg-rose-500/10 text-rose-200",
          )}
        >
          {status.message}
        </p>
      ) : null}

      <div className="overflow-hidden rounded-2xl border border-border bg-card/95 shadow-sm transition-colors supports-[backdrop-filter]:bg-card/80 supports-[backdrop-filter]:backdrop-blur">
        <table className="min-w-full divide-y divide-border/60 text-sm">
          <thead className="bg-muted/80">
            <tr>
              <th className="px-4 py-3 text-left font-semibold text-slate-600 dark:text-slate-300">Proyek</th>
              <th className="px-4 py-3 text-left font-semibold text-slate-600 dark:text-slate-300">Ringkasan</th>
              <th className="px-4 py-3 text-center font-semibold text-slate-600 dark:text-slate-300">Featured</th>
              <th className="w-40 px-4 py-3 text-right font-semibold text-slate-600 dark:text-slate-300">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
              <SortableContext items={items} strategy={verticalListSortingStrategy}>
                {projects.map((project) => (
                  <ProjectRow
                    key={project.id}
                    project={project}
                    onEdit={() => setEditor({ mode: "edit", project })}
                    onDelete={() => handleDelete(project.id)}
                    onFeature={(state) => handleFeature(project, state)}
                    disabled={isPending && pendingId !== project.id}
                    isPending={isPending && pendingId === project.id}
                  />
                ))}
              </SortableContext>
            </DndContext>
          </tbody>
        </table>
      </div>

      {editor ? (
        <ProjectEditor
          key={editor.project?.id ?? "create"}
          mode={editor.mode}
          project={editor.project}
          onSubmit={(values) =>
            editor.mode === "create" ? handleCreate(values) : handleUpdate(editor.project as Project, values)
          }
          onCancel={() => setEditor(null)}
          isPending={isPending}
        />
      ) : null}
    </div>
  );
}

type ProjectRowProps = {
  project: Project;
  onEdit: () => void;
  onDelete: () => void;
  onFeature: (state: boolean) => void;
  disabled: boolean;
  isPending: boolean;
};

function ProjectRow({ project, onEdit, onDelete, onFeature, disabled, isPending }: ProjectRowProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: project.id });
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };
  const badges = [
    project.category ? { label: project.category, tone: "neutral" as const } : null,
    project.duration_label ? { label: `Durasi ${project.duration_label}`, tone: "neutral" as const } : null,
    project.price_label ? { label: project.price_label, tone: "accent" as const } : null,
    project.budget_label ? { label: project.budget_label, tone: "accent" as const } : null,
  ].filter(Boolean) as { label: string; tone: "neutral" | "accent" }[];

  return (
    <tr
      ref={setNodeRef}
      style={style}
      className={clsx(
        "border-b border-border/60 last:border-b-0",
        isDragging ? "bg-primary/10" : "bg-transparent",
      )}
    >
      <td className="px-4 py-4">
        <div className="flex gap-3">
          <button
            type="button"
            className="rounded-md border border-transparent p-1 text-muted-foreground transition-colors hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40"
            {...listeners}
            {...attributes}
            aria-label="Ubah urutan"
          >
            <GripVertical className="h-4 w-4" />
          </button>
          <div>
            <p className="font-semibold text-foreground">{project.title}</p>
            {project.description ? (
              <p className="mt-1 text-xs text-muted-foreground">{project.description}</p>
            ) : null}
            {project.tech_stack.length ? (
              <div className="mt-2 flex flex-wrap gap-2">
                {project.tech_stack.map((tech) => (
                  <span
                    key={tech}
                    className="rounded-full border border-primary/40 bg-primary/10 px-2 py-0.5 text-[11px] text-primary"
                  >
                    {tech}
                  </span>
                ))}
              </div>
            ) : null}
          </div>
        </div>
      </td>
      <td className="px-4 py-4">
        {badges.length ? (
          <div className="flex flex-wrap gap-2 text-xs">
            {badges.map((badge) => (
              <span
                key={badge.label}
                className={clsx(
                  "rounded-full px-3 py-1",
                  badge.tone === "accent"
                    ? "bg-emerald-500/10 text-emerald-200"
                    : "border border-slate-200/70 text-slate-600 dark:border-slate-700 dark:text-slate-300",
                )}
              >
                {badge.label}
              </span>
            ))}
          </div>
        ) : (
          <span className="text-xs text-slate-400 dark:text-slate-600">-</span>
        )}
      </td>
      <td className="px-4 py-4 text-center">
        <Button
          type="button"
          variant={project.is_featured ? "secondary" : "ghost"}
          size="sm"
          onClick={() => onFeature(!project.is_featured)}
          disabled={disabled || isPending}
        >
          {project.is_featured ? (
            <>
              <Star className="h-4 w-4 text-amber-400" />
              <span className="sr-only">Hapus featured</span>
            </>
          ) : (
            <>
              <Sparkles className="h-4 w-4" />
              <span className="sr-only">Jadikan featured</span>
            </>
          )}
        </Button>
      </td>
      <td className="px-4 py-4">
        <div className="flex justify-end gap-2">
          <Button type="button" variant="ghost" size="sm" onClick={onEdit} disabled={disabled}>
            <Pencil className="h-4 w-4" />
            <span className="sr-only">Edit</span>
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={onDelete}
            disabled={disabled}
            className="text-destructive hover:text-destructive"
          >
            <Trash2 className="h-4 w-4" />
            <span className="sr-only">Hapus</span>
          </Button>
        </div>
      </td>
    </tr>
  );
}

type ProjectEditorProps = {
  mode: "create" | "edit";
  project?: Project;
  onSubmit: (values: ProjectFormValues) => void;
  onCancel: () => void;
  isPending: boolean;
};

function ProjectEditor({ mode, project, onSubmit, onCancel, isPending }: ProjectEditorProps) {
  const form = useForm<ProjectFormValues>({
    resolver: zodResolver(projectSchema),
    defaultValues: {
      title: project?.title ?? "",
      description: project?.description ?? "",
      tech_stack: project?.tech_stack ?? [],
      image_url: project?.image_url ?? "",
      project_url: project?.project_url ?? "",
      category: project?.category ?? "",
      duration_label: project?.duration_label ?? "",
      price_label: project?.price_label ?? "",
      budget_label: project?.budget_label ?? "",
      is_featured: project?.is_featured ?? false,
    },
  });

  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isDirty },
  } = form;

  const { fields, append, remove } = useFieldArray<ProjectFormValues>({
    name: "tech_stack" as FieldArrayPath<ProjectFormValues>,
    control,
  });

  const previewUrl = useWatch({ control, name: "image_url" });

  const submit = (values: ProjectFormValues) => {
    onSubmit(values);
  };

  return (
    <form
      onSubmit={handleSubmit(submit)}
      className="rounded-2xl border border-border bg-card/95 p-6 shadow-sm transition-colors supports-[backdrop-filter]:bg-card/80 supports-[backdrop-filter]:backdrop-blur"
    >
      <div className="flex items-center justify-between">
        <h3 className="text-base font-semibold text-foreground">
          {mode === "create" ? "Proyek baru" : "Edit proyek"}
        </h3>
        <Button type="button" variant="ghost" size="sm" onClick={onCancel}>
          Batal
        </Button>
      </div>

      <div className="mt-4 grid gap-4 lg:grid-cols-[2fr_1fr]">
        <div className="space-y-4">
          <Field label="Judul" error={errors.title?.message}>
            <Input
              placeholder="Contoh: Website SaaS"
              {...register("title")}
              invalid={Boolean(errors.title)}
            />
          </Field>
          <Field label="Deskripsi" error={errors.description?.message}>
            <Textarea
              placeholder="Ringkasan singkat hasil dan dampak proyek"
              {...register("description")}
              invalid={Boolean(errors.description)}
            />
          </Field>
          <Field label="Kategori" error={errors.category?.message}>
            <Input
              placeholder="Contoh: SaaS / Fintech"
              {...register("category")}
              invalid={Boolean(errors.category)}
            />
          </Field>
          <div className="grid gap-4 sm:grid-cols-2">
            <Field label="Durasi" error={errors.duration_label?.message}>
              <Input
                placeholder="Contoh: 6 minggu"
                {...register("duration_label")}
                invalid={Boolean(errors.duration_label)}
              />
            </Field>
            <Field label="Label harga" error={errors.price_label?.message}>
              <Input
                placeholder="Contoh: Growth package"
                {...register("price_label")}
                invalid={Boolean(errors.price_label)}
              />
            </Field>
            <div className="sm:col-span-2">
              <Field label="Label budget" error={errors.budget_label?.message}>
                <Input
                  placeholder="Contoh: IDR 150Jt"
                  {...register("budget_label")}
                  invalid={Boolean(errors.budget_label)}
                />
              </Field>
            </div>
          </div>
          <div className="space-y-2">
            <span className="text-sm font-medium text-foreground">Tech Stack</span>
            <div className="space-y-2">
              {fields.map((field, index) => (
                <div key={field.id} className="flex items-center gap-2">
                  <Input
                    placeholder="Contoh: Next.js"
                    {...register(`tech_stack.${index}` as const)}
                    invalid={Boolean(errors.tech_stack?.[index])}
                  />
                  <Button type="button" variant="ghost" size="sm" onClick={() => remove(index)}>
                    Hapus
                  </Button>
                </div>
              ))}
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={() => append("")}
                className="text-xs"
              >
                Tambah Tech
              </Button>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          <Field label="URL Proyek" error={errors.project_url?.message}>
            <Input
              placeholder="https://"
              {...register("project_url")}
              invalid={Boolean(errors.project_url)}
            />
          </Field>
          <Field label="Gambar Proyek" error={errors.image_url?.message}>
            <Controller
              name="image_url"
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
          <div className="rounded-2xl border border-border bg-card/90 p-4 text-sm text-muted-foreground shadow-sm transition-colors supports-[backdrop-filter]:bg-card/70 supports-[backdrop-filter]:backdrop-blur">
            <p className="font-semibold text-foreground">Pratinjau</p>
            <div className="mt-3 flex items-center justify-center">
              {previewUrl ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={previewUrl}
                  alt="Preview proyek"
                  className="h-40 w-full rounded-lg border border-border object-cover"
                />
              ) : (
                <div className="flex h-40 w-full items-center justify-center rounded-lg border border-dashed border-border text-xs text-muted-foreground">
                  Tidak ada gambar
                </div>
              )}
            </div>
          </div>
          <Field label="Featured" error={errors.is_featured?.message}>
            <Controller
              name="is_featured"
              control={control}
              render={({ field }) => (
                <Switch
                  checked={field.value ?? false}
                  onChange={(event) => field.onChange(event.target.checked)}
                />
              )}
            />
          </Field>
        </div>
      </div>

      <div className="mt-6 flex justify-end gap-2">
        <Button type="submit" disabled={isPending || !isDirty}>
          {mode === "create" ? "Simpan proyek" : "Perbarui proyek"}
        </Button>
      </div>
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
