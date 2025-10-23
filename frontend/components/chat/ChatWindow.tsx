"use client";

import { useMemo, useState } from "react";

import {
  askAssistant,
  ChatMessage,
  createAssistantReplyFromText,
  createUserMessage,
  initialMessages,
} from "@/lib/chat";
import { ChatBubble } from "./ChatBubble";
import { ChatInput } from "./ChatInput";

export function ChatWindow() {
  const [messages, setMessages] = useState<ChatMessage[]>(() => [
    ...initialMessages,
  ]);
  const [isLoading, setIsLoading] = useState(false);

  const handleSend = async (question: string) => {
    const userMessage = createUserMessage(question);
    setMessages((previous) => [...previous, userMessage]);
    setIsLoading(true);

    try {
      const assistantMessage = await askAssistant(question);
      setMessages((previous) => [...previous, assistantMessage]);
    } catch {
      const fallback = createAssistantReplyFromText(
        "Maaf, terjadi kendala saat memproses pesan. Silakan coba lagi.",
      );
      setMessages((previous) => [...previous, fallback]);
    } finally {
      setIsLoading(false);
    }
  };

  const orderedMessages = useMemo(() => {
    return [...messages].sort(
      (a, b) => a.createdAt.getTime() - b.createdAt.getTime(),
    );
  }, [messages]);

  return (
    <div className="flex h-full flex-col gap-4">
      <div className="flex-1 space-y-4 overflow-y-auto rounded-3xl border border-slate-200 bg-slate-50/60 p-6">
        {orderedMessages.map((message) => (
          <ChatBubble key={message.id} message={message} />
        ))}
        {isLoading ? (
          <div className="text-xs font-medium uppercase tracking-wide text-indigo-500">
            tany.ai sedang menyiapkan jawaban…
          </div>
        ) : null}
      </div>
      <ChatInput onSend={handleSend} disabled={isLoading} />
    </div>
  );
}
