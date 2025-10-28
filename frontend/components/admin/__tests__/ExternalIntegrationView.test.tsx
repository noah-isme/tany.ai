import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { ExternalIntegrationView } from "../ExternalIntegrationView";

const refreshMock = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ refresh: refreshMock }),
}));

describe("ExternalIntegrationView", () => {
  const baseSource = {
    id: "source-1",
    name: "noahis.me",
    baseUrl: "https://noahis.me",
    sourceType: "auto",
    enabled: true,
    lastSyncedAt: undefined,
    lastModified: undefined,
  };

  const baseItem = {
    id: "item-1",
    sourceName: "noahis.me",
    kind: "post",
    title: "Membangun sesuatu",
    summary: "Ringkasan",
    url: "https://noahis.me/post",
    visible: true,
    publishedAt: undefined,
    metadata: {},
  };

  beforeEach(() => {
    refreshMock.mockReset();
  });

  it("syncs source and shows success message", async () => {
    const syncSourceMock = vi
      .fn()
      .mockResolvedValue({ success: true, data: { itemsUpserted: 2, message: "Sinkron berhasil" } });
    const toggleItemMock = vi
      .fn()
      .mockResolvedValue({ success: true, data: { ...baseItem, visible: true } });

    render(
      <ExternalIntegrationView
        initialSources={[baseSource]}
        initialItems={[baseItem]}
        syncSource={syncSourceMock}
        toggleItem={toggleItemMock}
      />,
    );

    await userEvent.click(screen.getByRole("button", { name: /sinkron sekarang/i }));

    await waitFor(() => {
      expect(syncSourceMock).toHaveBeenCalledWith("source-1");
    });

    expect(await screen.findByText(/sinkron berhasil/i)).toBeInTheDocument();
    expect(refreshMock).toHaveBeenCalled();
  });

  it("toggles item visibility", async () => {
    const syncSourceMock = vi.fn().mockResolvedValue({ success: true, data: { itemsUpserted: 0 } });
    const toggleItemMock = vi.fn().mockResolvedValue({ success: true, data: { ...baseItem, visible: false } });

    render(
      <ExternalIntegrationView
        initialSources={[baseSource]}
        initialItems={[baseItem]}
        syncSource={syncSourceMock}
        toggleItem={toggleItemMock}
      />,
    );

    const switchInput = screen.getByLabelText(/atur visibilitas/i) as HTMLInputElement;
    expect(switchInput.checked).toBe(true);

    await userEvent.click(switchInput);

    await waitFor(() => {
      expect(toggleItemMock).toHaveBeenCalledWith({ id: "item-1", visible: false });
    });

    expect(switchInput.checked).toBe(false);
  });
});
