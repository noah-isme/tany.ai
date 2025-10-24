import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";

import { apiFetch } from "./api-client";
import { resolveApiUrl } from "./env";
import {
  askAssistant,
  createAssistantReplyFromText,
  createUserMessage,
  type ChatApiResponse,
} from "./chat";

vi.mock("./api-client", async (importOriginal) => {
  const actual = await importOriginal<typeof import("./api-client")>();
  return {
    ...actual,
    apiFetch: vi.fn(),
  };
});

describe("chat helpers", () => {
  beforeEach(() => {
    vi.mocked(apiFetch).mockReset();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("creates user messages with unique IDs", () => {
    const message = createUserMessage("Halo");
    expect(message.role).toBe("user");
    expect(message.id).toMatch(/user-/);
    expect(message.createdAt).toBeInstanceOf(Date);
  });

  it("creates assistant replies from plain text", () => {
    const message = createAssistantReplyFromText("Sukses");
    expect(message.role).toBe("assistant");
    expect(message.content).toBe("Sukses");
  });

  it("calls the chat API and maps the response", async () => {
    const payload: ChatApiResponse = {
      chatId: "123",
      answer: "Halo kembali",
      model: "mock-model",
      prompt: "prompt",
    };

    const isServerLike = typeof window === "undefined";
    let fetchMock: ReturnType<typeof vi.fn> | null = null;

    if (isServerLike) {
      vi.mocked(apiFetch).mockResolvedValueOnce(payload);
    } else {
      fetchMock = vi.fn().mockResolvedValue({
        ok: true,
        json: async () => payload,
        headers: new Headers({ "content-type": "application/json" }),
      } as Response);
      vi.stubGlobal("fetch", fetchMock);
    }

    const result = await askAssistant("Apa kabar?", null);

    if (isServerLike) {
      expect(apiFetch).toHaveBeenCalledWith("/api/v1/chat", {
        method: "POST",
        body: { question: "Apa kabar?" },
        withAuth: false,
      });
    } else if (fetchMock) {
      expect(fetchMock).toHaveBeenCalledWith(resolveApiUrl("/api/v1/chat"), {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify({ question: "Apa kabar?" }),
        cache: "no-store",
      });
    }

    expect(result.chatId).toBe("123");
    expect(result.message.role).toBe("assistant");
    expect(result.message.content).toBe("Halo kembali");
  });
});
