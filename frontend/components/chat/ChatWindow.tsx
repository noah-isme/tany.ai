"use client";

import { useEffect, useMemo, useState } from "react";

import {
  askAssistant,
  ChatMessage,
  createAssistantReplyFromText,
  createUserMessage,
} from "@/lib/chat";
import { fetchKnowledgeBase, type KnowledgeBase } from "@/lib/knowledge";
import { ChatBubble } from "./ChatBubble";
import { ChatInput } from "./ChatInput";

function buildAssistSnippet(base: KnowledgeBase): string {
  const services = base.services.slice(0, 3).map((service) => service.name);
  const projects = base.projects
    .filter((project) => project.isFeatured)
    .concat(base.projects)
    .slice(0, 1);
  const intro = base.profile.name
    ? `Hai! Saya ${base.profile.name}.`
    : "Hai!";
  const serviceText = services.length
    ? `Saya bisa membantu ${services.join(", ")}.`
    : "Saya siap membantu kebutuhanmu.";
  const projectText = projects.length
    ? `Contoh proyek terbaru: ${projects[0].title}.`
    : "Senang bisa berbagi portofolio saya.";

  return `${intro} ${serviceText} ${projectText}`;
}

type ChatWindowProps = {
  initialKnowledge?: KnowledgeBase | null;
};

export function ChatWindow({ initialKnowledge = null }: ChatWindowProps) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [chatId, setChatId] = useState<string | null>(null);
  const [knowledge, setKnowledge] = useState<KnowledgeBase | null>(initialKnowledge);
  const [knowledgeError, setKnowledgeError] = useState<string | null>(null);

  useEffect(() => {
    if (initialKnowledge) {
      return;
    }
    let active = true;
    fetchKnowledgeBase()
      .then((data) => {
        if (!active) return;
        setKnowledge(data);
      })
      .catch(() => {
        if (!active) return;
        setKnowledgeError("Gagal memuat knowledge base.");
      });
    return () => {
      active = false;
    };
  }, [initialKnowledge]);

  const handleSend = async (question: string) => {
    const userMessage = createUserMessage(question);
    setMessages((previous) => [...previous, userMessage]);
    setIsLoading(true);

    try {
      const reply = await askAssistant(question, chatId);
      setChatId(reply.chatId);
      setMessages((previous) => [...previous, reply.message]);
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

  const assistSnippet = knowledge ? buildAssistSnippet(knowledge) : null;

  return (
    <div className="flex h-full flex-col gap-4">
      <div className="flex-1 space-y-4 overflow-y-auto rounded-3xl border border-slate-200 bg-slate-50/60 p-6">
        {assistSnippet ? (
          <div className="rounded-2xl bg-white p-4 text-sm text-slate-700 shadow-sm">
            {assistSnippet}
          </div>
        ) : knowledgeError ? (
          <div className="rounded-2xl bg-red-50 p-4 text-sm text-red-600">
            {knowledgeError}
          </div>
        ) : null}
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
