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
import { clsx } from "clsx";
import { ArrowDown, ArrowUp, Check, GripVertical, Pencil, Trash2, XCircle } from "lucide-react";

import type { ActionResult } from "@/lib/action-result";
import type { Skill } from "@/lib/types/admin";
import { type SkillFormValues } from "@/lib/validators";

import { Button } from "./ui/button";
import { Input } from "./ui/input";

const emptyStatus = null as { type: "success" | "error"; message: string } | null;

export type SkillActions = {
  createSkill: (values: SkillFormValues) => Promise<ActionResult<Skill>>;
  updateSkill: (args: { id: string; values: SkillFormValues }) => Promise<ActionResult<Skill>>;
  deleteSkill: (args: { id: string }) => Promise<ActionResult<null>>;
  reorderSkill: (args: { items: { id: string; order: number }[] }) => Promise<ActionResult<null>>;
};

type SkillsManagerProps = SkillActions & {
  initialSkills: Skill[];
};

export function SkillsManager({ initialSkills, createSkill, updateSkill, deleteSkill, reorderSkill }: SkillsManagerProps) {
  const [skills, setSkills] = useState(initialSkills);
  const [newSkill, setNewSkill] = useState("");
  const [status, setStatus] = useState(emptyStatus);
  const [isPending, startTransition] = useTransition();
  const [pendingId, setPendingId] = useState<string | null>(null);

  useEffect(() => {
    setSkills(initialSkills);
  }, [initialSkills]);

  const sensors = useSensors(useSensor(PointerSensor));

  const handleCreate = () => {
    if (!newSkill.trim()) {
      setStatus({ type: "error", message: "Nama skill wajib diisi." });
      return;
    }
    startTransition(async () => {
      setStatus(null);
      const result = await createSkill({ name: newSkill.trim() });
      if (result.success) {
        setSkills((prev) => [...prev, result.data]);
        setNewSkill("");
        setStatus({ type: "success", message: "Skill ditambahkan." });
        return;
      }
      setStatus({ type: "error", message: result.error });
    });
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) {
      return;
    }
    setSkills((items) => {
      const oldIndex = items.findIndex((item) => item.id === active.id);
      const newIndex = items.findIndex((item) => item.id === over.id);
      const reordered = arrayMove(items, oldIndex, newIndex);
      startTransition(async () => {
        setStatus(null);
        const payload = reordered.map((skill, index) => ({ id: skill.id, order: index }));
        const result = await reorderSkill({ items: payload });
        if (!result.success) {
          setStatus({ type: "error", message: result.error });
        }
      });
      return reordered;
    });
  };

  const handleUpdate = (id: string, name: string) => {
    setPendingId(id);
    startTransition(async () => {
      const trimmed = name.trim();
      if (!trimmed) {
        setStatus({ type: "error", message: "Nama skill wajib diisi." });
        setPendingId(null);
        return;
      }
      const result = await updateSkill({ id, values: { name: trimmed } });
      if (result.success) {
        setSkills((items) =>
          items.map((item) => (item.id === id ? { ...item, name: result.data.name } : item)),
        );
        setStatus({ type: "success", message: "Skill diperbarui." });
      } else {
        setStatus({ type: "error", message: result.error });
      }
      setPendingId(null);
    });
  };

  const handleDelete = (id: string) => {
    if (typeof window !== "undefined" && !window.confirm("Hapus skill ini?")) {
      return;
    }
    setPendingId(id);
    startTransition(async () => {
      const result = await deleteSkill({ id });
      if (result.success) {
        setSkills((items) => items.filter((item) => item.id !== id));
        setStatus({ type: "success", message: "Skill dihapus." });
      } else {
        setStatus({ type: "error", message: result.error });
      }
      setPendingId(null);
    });
  };

  const handleMove = (id: string, direction: "up" | "down") => {
    let nextOrder: Skill[] | null = null;
    setSkills((items) => {
      const index = items.findIndex((item) => item.id === id);
      if (index < 0) {
        return items;
      }
      const targetIndex = direction === "up" ? Math.max(0, index - 1) : Math.min(items.length - 1, index + 1);
      if (targetIndex === index) {
        return items;
      }
      nextOrder = arrayMove(items, index, targetIndex);
      return nextOrder;
    });
    if (nextOrder) {
      startTransition(async () => {
        setStatus(null);
        const payload = nextOrder!.map((skill, order) => ({ id: skill.id, order }));
        const result = await reorderSkill({ items: payload });
        if (!result.success) {
          setStatus({ type: "error", message: result.error });
        }
      });
    }
  };

  const items = useMemo(() => skills.map((skill) => skill.id), [skills]);

  return (
    <div className="space-y-4">
      <div className="rounded-2xl border border-slate-200 bg-white/80 p-4 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
        <h2 className="text-sm font-semibold text-slate-900 dark:text-slate-100">Tambah Skill</h2>
        <div className="mt-3 flex flex-col gap-2 sm:flex-row">
          <Input
            value={newSkill}
            onChange={(event) => setNewSkill(event.target.value)}
            placeholder="Contoh: Prompt Engineering"
            disabled={isPending}
          />
          <Button type="button" onClick={handleCreate} disabled={isPending}>
            Tambah
          </Button>
        </div>
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
        <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
          <table className="min-w-full divide-y divide-slate-200 text-sm dark:divide-slate-800">
            <thead className="bg-slate-100/80 dark:bg-slate-900/50">
              <tr>
                <th scope="col" className="px-4 py-3 text-left font-semibold text-slate-600 dark:text-slate-300">
                  Skill
                </th>
                <th scope="col" className="w-32 px-4 py-3 text-right font-semibold text-slate-600 dark:text-slate-300">
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody>
              <SortableContext items={items} strategy={verticalListSortingStrategy}>
                {skills.map((skill, index) => (
                  <SkillRow
                    key={skill.id}
                    skill={skill}
                    onUpdate={handleUpdate}
                    onDelete={handleDelete}
                    disabled={isPending && pendingId !== skill.id}
                    isPending={isPending && pendingId === skill.id}
                    onMove={handleMove}
                    canMoveUp={index > 0}
                    canMoveDown={index < skills.length - 1}
                  />
                ))}
              </SortableContext>
            </tbody>
          </table>
        </DndContext>
      </div>
    </div>
  );
}

type SkillRowProps = {
  skill: Skill;
  onUpdate: (id: string, name: string) => void;
  onDelete: (id: string) => void;
  onMove: (id: string, direction: "up" | "down") => void;
  canMoveUp: boolean;
  canMoveDown: boolean;
  disabled: boolean;
  isPending: boolean;
};

function SkillRow({ skill, onUpdate, onDelete, onMove, canMoveUp, canMoveDown, disabled, isPending }: SkillRowProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: skill.id });
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };
  const [isEditing, setIsEditing] = useState(false);
  const [draft, setDraft] = useState(skill.name);

  const beginEdit = () => {
    setDraft(skill.name);
    setIsEditing(true);
  };

  const cancelEdit = () => {
    setDraft(skill.name);
    setIsEditing(false);
  };

  const save = () => {
    const nextValue = draft.trim();
    if (!nextValue) {
      return;
    }
    onUpdate(skill.id, nextValue);
    setDraft(nextValue);
    setIsEditing(false);
  };

  return (
    <tr
      ref={setNodeRef}
      style={style}
      className={clsx(
        "border-b border-slate-200/80 last:border-b-0 dark:border-slate-800/60",
        isDragging ? "bg-indigo-500/10" : "bg-white/0",
      )}
    >
      <td className="px-4 py-3">
        <div className="flex items-center gap-3">
          <button
            type="button"
            className="rounded-md border border-transparent p-1 text-slate-400 hover:text-indigo-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-400"
            {...listeners}
            {...attributes}
            aria-label="Ubah urutan"
          >
            <GripVertical className="h-4 w-4" />
          </button>
          {isEditing ? (
            <Input
              value={draft}
              onChange={(event) => setDraft(event.target.value)}
              disabled={isPending}
            />
          ) : (
            <span className="font-medium text-slate-800 dark:text-slate-100">{skill.name}</span>
          )}
        </div>
      </td>
      <td className="px-4 py-3">
        <div className="flex justify-end gap-2">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => onMove(skill.id, "up")}
            disabled={disabled || !canMoveUp}
            aria-label="Pindah ke atas"
          >
            <ArrowUp className="h-4 w-4" />
            <span className="sr-only">Pindah ke atas</span>
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => onMove(skill.id, "down")}
            disabled={disabled || !canMoveDown}
            aria-label="Pindah ke bawah"
          >
            <ArrowDown className="h-4 w-4" />
            <span className="sr-only">Pindah ke bawah</span>
          </Button>
          {isEditing ? (
            <>
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={save}
                disabled={isPending || draft.trim().length === 0}
              >
                <Check className="h-4 w-4" />
                <span className="sr-only">Simpan</span>
              </Button>
              <Button type="button" variant="ghost" size="sm" onClick={cancelEdit}>
                <XCircle className="h-4 w-4" />
                <span className="sr-only">Batal</span>
              </Button>
            </>
          ) : (
            <>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={beginEdit}
                disabled={disabled}
              >
                <Pencil className="h-4 w-4" />
                <span className="sr-only">Edit</span>
              </Button>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => onDelete(skill.id)}
                disabled={disabled}
                className="text-rose-400 hover:text-rose-500"
              >
                <Trash2 className="h-4 w-4" />
                <span className="sr-only">Hapus</span>
              </Button>
            </>
          )}
        </div>
      </td>
    </tr>
  );
}
