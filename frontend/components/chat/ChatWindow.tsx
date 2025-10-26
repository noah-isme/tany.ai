"use client";

import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useMemo, useRef, useState } from "react";

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
  const [retryMessage, setRetryMessage] = useState<string | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

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
    setRetryMessage(null);

    try {
      const reply = await askAssistant(question, chatId);
      setChatId(reply.chatId);
      setMessages((previous) => [...previous, reply.message]);
    } catch {
      const fallback = createAssistantReplyFromText(
        "Maaf, terjadi kendala saat memproses pesan. Silakan coba lagi.",
      );
      setMessages((previous) => [...previous, fallback]);
      setRetryMessage(question);
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
  const quickReplies = useMemo(() => {
    const base = ["Harga layanan", "Durasi pengerjaan", "Contoh proyek"];
    const servicePrompts = knowledge
      ? knowledge.services.slice(0, 3).map((service) => `Ceritakan tentang ${service.name}`)
      : [];
    return Array.from(new Set([...base, ...servicePrompts]));
  }, [knowledge]);

  useEffect(() => {
    if (!containerRef.current) return;
    containerRef.current.scrollTo({
      top: containerRef.current.scrollHeight,
      behavior: "smooth",
    });
  }, [orderedMessages, isLoading]);

  const handleQuickReply = (text: string) => {
    if (isLoading) {
      return;
    }
    void handleSend(text);
  };

  const handleRetry = () => {
    if (!retryMessage) {
      return;
    }
    void handleSend(retryMessage);
  };

  const showSkeleton = !knowledge && !knowledgeError;

  return (
    <div className="flex h-full flex-col gap-4">
      <div
        ref={containerRef}
        className="flex-1 space-y-4 overflow-y-auto rounded-3xl border border-white/10 bg-white/5 p-5 shadow-inner shadow-black/20"
      >
        {assistSnippet ? (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            className="rounded-3xl bg-white/10 p-4 text-sm leading-relaxed text-white/75 shadow-sm"
          >
            {assistSnippet}
          </motion.div>
        ) : null}
        {knowledgeError ? (
          <div className="rounded-3xl border border-red-500/20 bg-red-500/10 p-4 text-sm text-red-200">
            {knowledgeError}
          </div>
        ) : null}
        {showSkeleton ? <KnowledgeSkeleton /> : null}
        <AnimatePresence initial={false}>
          {orderedMessages.map((message) => (
            <motion.div
              key={message.id}
              initial={{ opacity: 0, y: 24 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -12 }}
              transition={{ duration: 0.25, ease: "easeOut" }}
            >
              <ChatBubble message={message} />
            </motion.div>
          ))}
        </AnimatePresence>
        {isLoading ? <TypingIndicator /> : null}
        <div />
      </div>
      <div className="space-y-3">
        {quickReplies.length ? (
          <div className="flex flex-wrap gap-2">
            {quickReplies.map((reply) => (
              <button
                key={reply}
                type="button"
                onClick={() => handleQuickReply(reply)}
                disabled={isLoading}
                className="rounded-full border border-white/15 px-3 py-1.5 text-sm text-white/80 transition hover:bg-white/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400 disabled:opacity-50"
              >
                {reply}
              </button>
            ))}
          </div>
        ) : null}
        {retryMessage ? (
          <div className="flex items-center justify-between rounded-2xl border border-yellow-400/30 bg-yellow-400/10 px-4 py-3 text-xs text-yellow-100">
            <span>Kendala jaringan. Coba kirim ulang?</span>
            <button
              type="button"
              onClick={handleRetry}
              className="rounded-lg border border-yellow-200/30 px-3 py-1 font-semibold uppercase tracking-wide text-[10px] text-yellow-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
            >
              Retry
            </button>
          </div>
        ) : null}
        <ChatInput onSend={handleSend} disabled={isLoading} />
      </div>
    </div>
  );
}

function KnowledgeSkeleton() {
  return (
    <div className="space-y-4">
      {[0, 1, 2].map((index) => (
        <div
          key={index}
          className="animate-pulse rounded-3xl border border-white/5 bg-white/[0.08] p-4"
        >
          <div className="h-3 w-1/4 rounded-full bg-white/10" />
          <div className="mt-3 h-4 w-3/4 rounded-full bg-white/10" />
          <div className="mt-2 h-4 w-2/3 rounded-full bg-white/5" />
        </div>
      ))}
    </div>
  );
}

function TypingIndicator() {
  return (
    <div className="flex justify-start">
      <div className="rounded-3xl bg-white/10 px-4 py-2 text-xs text-white/80">
        <motion.div
          className="flex items-center gap-1"
          initial="animate"
          animate="animate"
          variants={{
            animate: {
              transition: { staggerChildren: 0.16, repeat: Infinity },
            },
          }}
        >
          {[0, 1, 2].map((index) => (
            <motion.span
              key={index}
              className="inline-block h-2 w-2 rounded-full bg-white/70"
              variants={{
                animate: {
                  opacity: [0.3, 1, 0.3],
                  scale: [0.9, 1.1, 0.9],
                  transition: { duration: 0.9, repeat: Infinity },
                },
              }}
            />
          ))}
        </motion.div>
      </div>
    </div>
  );
}
