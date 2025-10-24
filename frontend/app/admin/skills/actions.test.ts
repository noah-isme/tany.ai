import { describe, expect, it, beforeEach, vi } from "vitest";

import { ApiError } from "@/lib/api-client";

import { createSkillAction } from "./actions";

const createSkillMock = vi.fn();

vi.mock("next/cache", () => ({
  revalidatePath: vi.fn(),
}));

vi.mock("@/lib/admin-api", () => ({
  createSkill: (...args: unknown[]) => createSkillMock(...args),
  updateSkill: vi.fn(),
  deleteSkill: vi.fn(),
  reorderSkills: vi.fn(),
}));

describe("createSkillAction", () => {
  beforeEach(() => {
    createSkillMock.mockReset();
  });

  it("returns validation error when name empty", async () => {
    const result = await createSkillAction({ name: "" });
    expect(result.success).toBe(false);
    expect(result.fieldErrors?.name).toBeDefined();
    expect(createSkillMock).not.toHaveBeenCalled();
  });

  it("creates skill when data valid", async () => {
    const skill = { id: "1", name: "Prompt Engineering", order: 0 };
    createSkillMock.mockResolvedValue(skill);

    const result = await createSkillAction({ name: "Prompt Engineering" });

    expect(result.success).toBe(true);
    expect(result.data).toEqual(skill);
    expect(createSkillMock).toHaveBeenCalledWith({ name: "Prompt Engineering" });
  });

  it("normalizes api error message", async () => {
    createSkillMock.mockRejectedValue(new ApiError("Nama sudah digunakan", 409));

    const result = await createSkillAction({ name: "Duplicate" });

    expect(result.success).toBe(false);
    expect(result.error).toBe("Nama sudah digunakan");
  });
});
