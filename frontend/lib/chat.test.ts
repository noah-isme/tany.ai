import { describe, expect, it } from "vitest";

import { buildSystemPrompt, knowledgeBase } from "@/data/knowledge";

import {
  askAssistant,
  createAssistantMessage,
  createAssistantReplyFromText,
  createUserMessage,
  initialMessages,
  systemMessage,
} from "./chat";

describe("chat helpers", () => {
  it("membangun pesan system yang konsisten dengan knowledge base", () => {
    expect(systemMessage.content).toEqual(buildSystemPrompt(knowledgeBase));
  });

  it("membuat pesan user yang memiliki ID dan timestamp", () => {
    const message = createUserMessage("Halo");

    expect(message.role).toBe("user");
    expect(message.id).toMatch(/user-/);
    expect(message.createdAt).toBeInstanceOf(Date);
  });

  it("menghasilkan jawaban AI mock yang menyertakan pertanyaan pengguna", () => {
    const message = createAssistantMessage("Apa saja layananmu?");

    expect(message.role).toBe("assistant");
    expect(message.content).toContain("Pertanyaan diterima");
    expect(message.content).toContain("Apa saja layananmu?");
  });

  it("dapat membuat balasan manual untuk skenario error", () => {
    const custom = createAssistantReplyFromText("Mohon coba lagi");
    expect(custom.content).toBe("Mohon coba lagi");
    expect(custom.role).toBe("assistant");
  });

  it("askAssistant mengembalikan pesan asisten", async () => {
    const result = await askAssistant("Ada paket apa saja?");
    expect(result.role).toBe("assistant");
  });

  it("pesan awal hanya terdiri dari system message", () => {
    expect(initialMessages).toHaveLength(1);
    expect(initialMessages[0].role).toBe("system");
  });
});
