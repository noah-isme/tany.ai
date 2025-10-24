import { resolveApiUrl } from "./env";

export type ChatRole = "user" | "assistant" | "system";

export type ChatMessage = {
  id: string;
  role: ChatRole;
  content: string;
  createdAt: Date;
};

export type ChatApiResponse = {
  chatId: string;
  answer: string;
  model: string;
  prompt: string;
};

export type AssistantReply = {
  chatId: string;
  prompt: string;
  model: string;
  message: ChatMessage;
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

export function createAssistantReplyFromText(content: string): ChatMessage {
  return {
    id: createId("assistant"),
    role: "assistant",
    content,
    createdAt: new Date(),
  };
}

export async function askAssistant(
  question: string,
  chatId?: string | null,
): Promise<AssistantReply> {
  const payload: { question: string; chatId?: string } = { question };
  if (chatId) {
    payload.chatId = chatId;
  }

  if (typeof window === "undefined") {
    const { apiFetch } = await import("./api-client");
    const response = await apiFetch<ChatApiResponse>("/api/v1/chat", {
      method: "POST",
      body: payload,
      withAuth: false,
    });

    return {
      chatId: response.chatId,
      prompt: response.prompt,
      model: response.model,
      message: createAssistantReplyFromText(response.answer),
    };
  }

  const response = await fetch(resolveApiUrl("/api/v1/chat"), {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify(payload),
    cache: "no-store",
  });
  if (!response.ok) {
    throw new Error(`Failed to call chat API: ${response.status}`);
  }
  const data = (await response.json()) as ChatApiResponse;

  return {
    chatId: data.chatId,
    prompt: data.prompt,
    model: data.model,
    message: createAssistantReplyFromText(data.answer),
  };
}
