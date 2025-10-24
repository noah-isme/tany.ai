import { beforeEach, describe, expect, it, vi } from "vitest";

import { ApiError } from "@/lib/api-client";

import {
  createServiceAction,
  toggleServiceAction,
} from "./actions";

const createServiceMock = vi.fn();
const toggleServiceMock = vi.fn();

vi.mock("next/cache", () => ({
  revalidatePath: vi.fn(),
}));

vi.mock("@/lib/admin-api", () => ({
  createService: (...args: unknown[]) => createServiceMock(...args),
  updateService: vi.fn(),
  deleteService: vi.fn(),
  reorderServices: vi.fn(),
  toggleService: (...args: unknown[]) => toggleServiceMock(...args),
}));

describe("services actions", () => {
  beforeEach(() => {
    createServiceMock.mockReset();
    toggleServiceMock.mockReset();
  });

  it("rejects when price_max less than price_min", async () => {
    const result = await createServiceAction({
      name: "Consulting",
      description: "",
      price_min: 5000,
      price_max: 1000,
      currency: "idr",
      duration_label: "1 minggu",
      is_active: true,
    });

    expect(result.success).toBe(false);
    expect(result.fieldErrors?.price_max).toBe("Harga maksimal harus >= harga minimal");
    expect(createServiceMock).not.toHaveBeenCalled();
  });

  it("creates service with normalized payload", async () => {
    const service = {
      id: "svc-1",
      name: "Consulting",
      description: "",
      price_min: 1000,
      price_max: 5000,
      currency: "IDR",
      duration_label: "1 minggu",
      is_active: true,
      order: 0,
    };
    createServiceMock.mockResolvedValue(service);

    const result = await createServiceAction({
      name: "Consulting",
      description: "",
      price_min: 1000,
      price_max: 5000,
      currency: "idr",
      duration_label: "1 minggu",
      is_active: true,
    });

    expect(result.success).toBe(true);
    expect(result.data).toEqual(service);
    expect(createServiceMock).toHaveBeenCalledWith({
      name: "Consulting",
      description: "",
      price_min: 1000,
      price_max: 5000,
      currency: "IDR",
      duration_label: "1 minggu",
      is_active: true,
    });
  });

  it("normalizes toggle error", async () => {
    toggleServiceMock.mockRejectedValue(new ApiError("Tidak ditemukan", 404));

    const result = await toggleServiceAction({ id: "svc-1", is_active: false });

    expect(result.success).toBe(false);
    expect(result.error).toBe("Tidak ditemukan");
  });

  it("toggles service status", async () => {
    const toggled = {
      id: "svc-1",
      name: "Consulting",
      description: "",
      price_min: 1000,
      price_max: 5000,
      currency: "IDR",
      duration_label: "1 minggu",
      is_active: false,
      order: 0,
    };
    toggleServiceMock.mockResolvedValue(toggled);

    const result = await toggleServiceAction({ id: "svc-1", is_active: false });

    expect(result.success).toBe(true);
    expect(result.data).toEqual(toggled);
    expect(toggleServiceMock).toHaveBeenCalledWith("svc-1", false);
  });
});
