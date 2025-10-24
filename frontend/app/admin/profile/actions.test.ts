import { beforeEach, describe, expect, it, vi } from "vitest";

import { ApiError } from "@/lib/api-client";

import { updateProfileAction } from "./actions";

const updateProfileMock = vi.fn();

vi.mock("next/cache", () => ({
  revalidatePath: vi.fn(),
}));

vi.mock("@/lib/admin-api", () => ({
  updateProfile: (...args: unknown[]) => updateProfileMock(...args),
  fetchProfile: vi.fn(),
}));

describe("updateProfileAction", () => {
  beforeEach(() => {
    updateProfileMock.mockReset();
  });

  it("returns validation error for invalid email", async () => {
    const result = await updateProfileAction({
      name: "Tanya",
      title: "AI Lead",
      bio: "",
      email: "invalid-email",
      phone: "",
      location: "",
      avatar_url: "https://example.com/avatar.png",
    });

    expect(result.success).toBe(false);
    expect(result.fieldErrors?.email).toBe("Format email tidak valid");
    expect(updateProfileMock).not.toHaveBeenCalled();
  });

  it("updates profile when payload valid", async () => {
    const payload = {
      id: "profile-1",
      name: "Tanya",
      title: "AI Lead",
      bio: "Bio",
      email: "admin@example.com",
      phone: "123",
      location: "Jakarta",
      avatar_url: "https://example.com/avatar.png",
    };
    updateProfileMock.mockResolvedValue(payload);

    const result = await updateProfileAction(payload);

    expect(result.success).toBe(true);
    expect(result.data).toEqual(payload);
    expect(updateProfileMock).toHaveBeenCalledWith({
      name: payload.name,
      title: payload.title,
      bio: payload.bio,
      email: payload.email,
      phone: payload.phone,
      location: payload.location,
      avatar_url: payload.avatar_url,
    });
  });

  it("normalizes api error response", async () => {
    updateProfileMock.mockRejectedValue(new ApiError("Gagal", 400));

    const result = await updateProfileAction({
      name: "Tanya",
      title: "AI Lead",
      bio: "",
      email: "admin@example.com",
      phone: "",
      location: "",
      avatar_url: "https://example.com/avatar.png",
    });

    expect(result.success).toBe(false);
    expect(result.error).toBe("Gagal");
  });
});
