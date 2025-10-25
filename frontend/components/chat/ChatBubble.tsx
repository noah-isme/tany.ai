import type { ChatMessage } from "@/lib/chat";

function formatTime(date: Date): string {
  try {
    return date.toLocaleTimeString("id-ID", {
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch {
    return "";
  }
}

const roleLabel: Record<ChatMessage["role"], string> = {
  user: "Anda",
  assistant: "tany.ai",
  system: "Instruksi Sistem",
};

export function ChatBubble({ message }: { message: ChatMessage }) {
  const timestamp = formatTime(message.createdAt);
  const isUser = message.role === "user";
  const isAssistant = message.role === "assistant";
  const alignment = isUser ? "items-end text-right" : "items-start text-left";
  const bubbleClass = isUser
    ? "bg-gradient-to-r from-violet-500 to-cyan-400 text-white shadow"
    : isAssistant
      ? "bg-white/90 text-slate-900 shadow-sm"
      : "bg-slate-900/90 text-slate-50";

  return (
    <div className={`flex w-full flex-col gap-1 ${alignment}`}>
      <span className="text-[10px] uppercase tracking-[0.32em] text-white/40">
        {roleLabel[message.role]}
        {timestamp ? ` · ${timestamp}` : ""}
      </span>
      <div
        className={`max-w-xl whitespace-pre-line rounded-3xl px-4 py-3 text-sm leading-relaxed ${bubbleClass}`}
      >
        {message.content}
      </div>
    </div>
  );
}
