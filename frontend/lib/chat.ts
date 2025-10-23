import { buildSystemPrompt, createMockAnswer, knowledgeBase } from "@/data/knowledge";

export type ChatRole = "user" | "assistant" | "system";

export type ChatMessage = {
  id: string;
  role: ChatRole;
  content: string;
  createdAt: Date;
};

export const systemMessage: ChatMessage = {
  id: "system-message",
  role: "system",
  content: buildSystemPrompt(knowledgeBase),
  createdAt: new Date(),
};

function createId(prefix: string): string {
  if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
    return `${prefix}-${crypto.randomUUID()}`;
  }
  return `${prefix}-${Math.random().toString(36).slice(2)}`;
}

export function createUserMessage(content: string): ChatMessage {
  return {
    id: createId("user"),
    role: "user",
    content,
    createdAt: new Date(),
  };
}

export function createAssistantMessage(question: string): ChatMessage {
  return createAssistantReplyFromText(
    createMockAnswer(question, knowledgeBase),
  );
}

export function createAssistantReplyFromText(content: string): ChatMessage {
  return {
    id: createId("assistant"),
    role: "assistant",
    content,
    createdAt: new Date(),
  };
}

export async function askAssistant(question: string): Promise<ChatMessage> {
  await new Promise((resolve) => setTimeout(resolve, 400));
  return createAssistantMessage(question);
}

export const initialMessages: ChatMessage[] = [systemMessage];
