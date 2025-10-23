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

const bubbleStyles: Record<ChatMessage["role"], string> = {
  user: "bg-indigo-600 text-white shadow-sm",
  assistant: "bg-white text-slate-900 border border-slate-200 shadow-sm",
  system: "bg-slate-900/90 text-slate-50",
};

const alignmentStyles: Record<ChatMessage["role"], string> = {
  user: "items-end text-right",
  assistant: "items-start text-left",
  system: "items-start text-left",
};

export function ChatBubble({ message }: { message: ChatMessage }) {
  const timestamp = formatTime(message.createdAt);
  const bubbleClass = bubbleStyles[message.role];
  const alignment = alignmentStyles[message.role];

  return (
    <div className={`flex w-full flex-col gap-1 ${alignment}`}>
      <span className="text-xs uppercase tracking-wide text-slate-400">
        {roleLabel[message.role]}
        {timestamp ? ` · ${timestamp}` : ""}
      </span>
      <div
        className={`max-w-xl whitespace-pre-line rounded-2xl px-4 py-3 text-sm leading-relaxed ${bubbleClass}`}
      >
        {message.content}
      </div>
    </div>
  );
}
