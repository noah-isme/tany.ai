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
import { useForm, Controller } from "react-hook-form";
import { BadgeCheck, GripVertical, Pencil, Trash2 } from "lucide-react";

import type { ActionResult } from "@/lib/action-result";
import type { Service } from "@/lib/types/admin";
import { serviceSchema, type ServiceFormValues } from "@/lib/validators";

import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Switch } from "./ui/switch";
import { Textarea } from "./ui/textarea";

const emptyStatus = null as { type: "success" | "error"; message: string } | null;

type ServiceActions = {
  createService: (values: ServiceFormValues) => Promise<ActionResult<Service>>;
  updateService: (args: { id: string; values: ServiceFormValues }) => Promise<ActionResult<Service>>;
  deleteService: (args: { id: string }) => Promise<ActionResult<null>>;
  reorderService: (args: { items: { id: string; order: number }[] }) => Promise<ActionResult<null>>;
  toggleService: (args: { id: string; is_active: boolean }) => Promise<ActionResult<Service>>;
};

type ServicesManagerProps = ServiceActions & {
  initialServices: Service[];
};

export function ServicesManager({
  initialServices,
  createService,
  updateService,
  deleteService,
  reorderService,
  toggleService,
}: ServicesManagerProps) {
  const [services, setServices] = useState(initialServices);
  const [status, setStatus] = useState(emptyStatus);
  const [editor, setEditor] = useState<{ mode: "create" | "edit"; service?: Service } | null>(null);
  const [isPending, startTransition] = useTransition();
  const [pendingId, setPendingId] = useState<string | null>(null);

  useEffect(() => {
    setServices(initialServices);
  }, [initialServices]);

  const sensors = useSensors(useSensor(PointerSensor));
  const items = useMemo(() => services.map((service) => service.id), [services]);

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) {
      return;
    }
    setServices((prev) => {
      const oldIndex = prev.findIndex((item) => item.id === active.id);
      const newIndex = prev.findIndex((item) => item.id === over.id);
      const reordered = arrayMove(prev, oldIndex, newIndex);
      startTransition(async () => {
        const payload = reordered.map((service, index) => ({ id: service.id, order: index }));
        const result = await reorderService({ items: payload });
        if (!result.success) {
          setStatus({ type: "error", message: result.error });
        }
      });
      return reordered;
    });
  };

  const handleDelete = (id: string) => {
    if (typeof window !== "undefined" && !window.confirm("Hapus layanan ini?")) {
      return;
    }
    setPendingId(id);
    startTransition(async () => {
      const result = await deleteService({ id });
      if (result.success) {
        setServices((prev) => prev.filter((service) => service.id !== id));
        setStatus({ type: "success", message: "Layanan dihapus." });
      } else {
        setStatus({ type: "error", message: result.error });
      }
      setPendingId(null);
    });
  };

  const handleToggle = (service: Service) => {
    setPendingId(service.id);
    startTransition(async () => {
      const result = await toggleService({ id: service.id, is_active: !service.is_active });
      if (result.success) {
        setServices((prev) =>
          prev.map((item) => (item.id === service.id ? { ...item, is_active: result.data.is_active } : item)),
        );
      } else {
        setStatus({ type: "error", message: result.error });
      }
      setPendingId(null);
    });
  };

  const handleCreate = (values: ServiceFormValues) => {
    startTransition(async () => {
      const result = await createService(values);
      if (result.success) {
        setServices((prev) => [...prev, result.data]);
        setStatus({ type: "success", message: "Layanan dibuat." });
        setEditor(null);
      } else {
        setStatus({ type: "error", message: result.error });
      }
    });
  };

  const handleUpdate = (service: Service, values: ServiceFormValues) => {
    startTransition(async () => {
      const result = await updateService({ id: service.id, values });
      if (result.success) {
        setServices((prev) =>
          prev.map((item) => (item.id === service.id ? { ...item, ...result.data } : item)),
        );
        setStatus({ type: "success", message: "Layanan diperbarui." });
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
          <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Daftar Layanan</h2>
          <p className="text-xs text-slate-500 dark:text-slate-400">
            Drag & drop untuk mengatur prioritas layanan di knowledge base.
          </p>
        </div>
        <Button type="button" onClick={() => setEditor({ mode: "create" })}>
          Tambah Layanan
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

      <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white/80 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
        <table className="min-w-full divide-y divide-slate-200 text-sm dark:divide-slate-800">
          <thead className="bg-slate-100/70 dark:bg-slate-900/40">
            <tr>
              <th className="px-4 py-3 text-left font-semibold text-slate-600 dark:text-slate-300">Layanan</th>
              <th className="px-4 py-3 text-left font-semibold text-slate-600 dark:text-slate-300">Harga</th>
              <th className="px-4 py-3 text-center font-semibold text-slate-600 dark:text-slate-300">Status</th>
              <th className="w-40 px-4 py-3 text-right font-semibold text-slate-600 dark:text-slate-300">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
              <SortableContext items={items} strategy={verticalListSortingStrategy}>
                {services.map((service) => (
                  <ServiceRow
                    key={service.id}
                    service={service}
                    onEdit={() => setEditor({ mode: "edit", service })}
                    onDelete={() => handleDelete(service.id)}
                    onToggle={() => handleToggle(service)}
                    disabled={isPending && pendingId !== service.id}
                    isPending={isPending && pendingId === service.id}
                  />
                ))}
              </SortableContext>
            </DndContext>
          </tbody>
        </table>
      </div>

      {editor ? (
        <ServiceEditor
          key={editor.service?.id ?? "create"}
          mode={editor.mode}
          service={editor.service}
          onSubmit={(values) =>
            editor.mode === "create" ? handleCreate(values) : handleUpdate(editor.service as Service, values)
          }
          onCancel={() => setEditor(null)}
          isPending={isPending}
        />
      ) : null}
    </div>
  );
}

type ServiceRowProps = {
  service: Service;
  onEdit: () => void;
  onDelete: () => void;
  onToggle: () => void;
  disabled: boolean;
  isPending: boolean;
};

function ServiceRow({ service, onEdit, onDelete, onToggle, disabled, isPending }: ServiceRowProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: service.id });
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  const priceRange = formatPriceRange(service);

  return (
    <tr
      ref={setNodeRef}
      style={style}
      className={clsx(
        "border-b border-slate-200/80 last:border-b-0 dark:border-slate-800/60",
        isDragging ? "bg-indigo-500/10" : "bg-white/0",
      )}
    >
      <td className="px-4 py-4">
        <div className="flex gap-3">
          <button
            type="button"
            className="rounded-md border border-transparent p-1 text-slate-400 hover:text-indigo-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-400"
            {...listeners}
            {...attributes}
            aria-label="Ubah urutan"
          >
            <GripVertical className="h-4 w-4" />
          </button>
          <div>
            <p className="font-semibold text-slate-900 dark:text-slate-100">{service.name}</p>
            {service.description ? (
              <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">{service.description}</p>
            ) : null}
            {service.duration_label ? (
              <p className="mt-1 inline-flex items-center gap-1 rounded-full border border-indigo-400/40 bg-indigo-500/10 px-2 py-0.5 text-[11px] text-indigo-300">
                <BadgeCheck className="h-3 w-3" /> {service.duration_label}
              </p>
            ) : null}
          </div>
        </div>
      </td>
      <td className="px-4 py-4 text-sm text-slate-600 dark:text-slate-300">{priceRange}</td>
      <td className="px-4 py-4 text-center">
        <Switch
          checked={service.is_active}
          onChange={(event) => {
            event.preventDefault();
            onToggle();
          }}
          disabled={disabled || isPending}
          aria-label={`Toggle status layanan ${service.name}`}
        />
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
            className="text-rose-400 hover:text-rose-500"
          >
            <Trash2 className="h-4 w-4" />
            <span className="sr-only">Hapus</span>
          </Button>
        </div>
      </td>
    </tr>
  );
}

type ServiceEditorProps = {
  mode: "create" | "edit";
  service?: Service;
  onSubmit: (values: ServiceFormValues) => void;
  onCancel: () => void;
  isPending: boolean;
};

function ServiceEditor({ mode, service, onSubmit, onCancel, isPending }: ServiceEditorProps) {
  const form = useForm<ServiceFormValues>({
    resolver: zodResolver(serviceSchema),
    defaultValues: {
      name: service?.name ?? "",
      description: service?.description ?? "",
      price_min: service?.price_min ?? null,
      price_max: service?.price_max ?? null,
      currency: service?.currency ?? "IDR",
      duration_label: service?.duration_label ?? "",
      is_active: service?.is_active ?? true,
    },
  });

  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isDirty },
  } = form;

  const submit = (values: ServiceFormValues) => {
    onSubmit(values);
  };

  return (
    <form
      onSubmit={handleSubmit(submit)}
      className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70"
    >
      <div className="flex items-center justify-between">
        <h3 className="text-base font-semibold text-slate-900 dark:text-slate-100">
          {mode === "create" ? "Layanan baru" : "Edit layanan"}
        </h3>
        <Button type="button" variant="ghost" size="sm" onClick={onCancel}>
          Batal
        </Button>
      </div>

      <div className="mt-4 grid gap-4 sm:grid-cols-2">
        <Field label="Nama" error={errors.name?.message}>
          <Input
            placeholder="Contoh: Website Development"
            {...register("name")}
            invalid={Boolean(errors.name)}
          />
        </Field>
        <Field label="Durasi" error={errors.duration_label?.message}>
          <Input
            placeholder="3-6 minggu"
            {...register("duration_label")}
            invalid={Boolean(errors.duration_label)}
          />
        </Field>
        <Field label="Deskripsi" error={errors.description?.message} className="sm:col-span-2">
          <Textarea
            placeholder="Ringkasan layanan"
            {...register("description")}
            invalid={Boolean(errors.description)}
          />
        </Field>
        <Field label="Harga Minimum" error={errors.price_min?.message}>
          <Controller
            name="price_min"
            control={control}
            render={({ field }) => (
              <Input
                type="number"
                min={0}
                step="1000"
                value={field.value ?? ""}
                onChange={(event) => field.onChange(event.target.value === "" ? null : Number(event.target.value))}
              />
            )}
          />
        </Field>
        <Field label="Harga Maksimum" error={errors.price_max?.message}>
          <Controller
            name="price_max"
            control={control}
            render={({ field }) => (
              <Input
                type="number"
                min={0}
                step="1000"
                value={field.value ?? ""}
                onChange={(event) => field.onChange(event.target.value === "" ? null : Number(event.target.value))}
              />
            )}
          />
        </Field>
        <Field label="Mata Uang" error={errors.currency?.message}>
          <Input
            placeholder="IDR"
            {...register("currency")}
            invalid={Boolean(errors.currency)}
          />
        </Field>
        <Field label="Aktif" className="flex items-center gap-3" error={errors.is_active?.message}>
          <Controller
            name="is_active"
            control={control}
            render={({ field }) => (
              <Switch
                checked={field.value ?? true}
                onChange={(event) => field.onChange(event.target.checked)}
              />
            )}
          />
        </Field>
      </div>

      <div className="mt-6 flex justify-end gap-2">
        <Button type="submit" disabled={isPending || !isDirty}>
          {mode === "create" ? "Simpan layanan" : "Perbarui layanan"}
        </Button>
      </div>
    </form>
  );
}

type FieldProps = {
  label: string;
  error?: string;
  children: React.ReactNode;
  className?: string;
};

function Field({ label, error, children, className }: FieldProps) {
  return (
    <label className={clsx("space-y-2 text-sm font-medium text-slate-700 dark:text-slate-200", className)}>
      <span>{label}</span>
      <div>{children}</div>
      {error ? (
        <span className="block text-xs font-normal text-rose-300" role="alert">
          {error}
        </span>
      ) : null}
    </label>
  );
}

function formatPriceRange(service: Service): string {
  const hasMin = typeof service.price_min === "number";
  const hasMax = typeof service.price_max === "number";
  const currency = service.currency || "IDR";
  const formatter = new Intl.NumberFormat("id-ID");
  if (hasMin && hasMax) {
    return `${currency} ${formatter.format(service.price_min!)} – ${formatter.format(service.price_max!)}`;
  }
  if (hasMin) {
    return `Mulai ${currency} ${formatter.format(service.price_min!)}`;
  }
  if (hasMax) {
    return `Hingga ${currency} ${formatter.format(service.price_max!)}`;
  }
  return "Hubungi untuk harga";
}
