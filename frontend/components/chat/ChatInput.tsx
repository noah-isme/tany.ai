"use client";

import {
  FormEvent,
  KeyboardEvent,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";

type ChatInputProps = {
  onSend: (message: string) => Promise<void> | void;
  placeholder?: string;
  disabled?: boolean;
};

type Command = {
  id: string;
  label: string;
  description: string;
  value: string;
};

const commands: Command[] = [
  {
    id: "pricing",
    label: "Pricing",
    description: "Lihat struktur paket dan kisaran biaya",
    value: "Bisakah jelaskan paket harga layanan Anda?",
  },
  {
    id: "lead",
    label: "Lead magnet",
    description: "Tanya proses konsultasi awal",
    value: "Bagaimana proses konsultasi awal jika saya ingin mulai?",
  },
  {
    id: "timeline",
    label: "Timeline",
    description: "Estimasi durasi pengerjaan proyek",
    value: "Berapa estimasi durasi untuk proyek standar?",
  },
];

export function ChatInput({ onSend, placeholder, disabled }: ChatInputProps) {
  const [message, setMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showCommands, setShowCommands] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleSubmit = useCallback(
    async (event: FormEvent<HTMLFormElement>) => {
      event.preventDefault();
      const trimmed = message.trim();
      if (!trimmed) {
        return;
      }

      try {
        setIsSubmitting(true);
        await onSend(trimmed);
        setMessage("");
      } finally {
        setIsSubmitting(false);
      }
    },
    [message, onSend],
  );

  const handleKeyDown = (event: KeyboardEvent<HTMLTextAreaElement>) => {
    if (event.key === "/" && message === "") {
      event.preventDefault();
      setShowCommands(true);
    }
    if (event.key === "Escape") {
      setShowCommands(false);
    }
  };

  useEffect(() => {
    if (!showCommands) {
      return;
    }
    textareaRef.current?.focus();
  }, [showCommands]);

  const filteredCommands = useMemo(() => commands, []);

  const applyCommand = (command: Command) => {
    setMessage(command.value);
    setShowCommands(false);
    requestAnimationFrame(() => textareaRef.current?.focus());
  };

  return (
    <div className="relative">
      {showCommands ? (
        <div className="absolute bottom-[72px] left-0 right-0 z-30 space-y-2 rounded-2xl border border-white/10 bg-[var(--surface)]/95 p-4 shadow-[0_20px_40px_rgba(8,12,24,0.45)] backdrop-blur-xl">
          <p className="text-xs uppercase tracking-[0.32em] text-white/40">Command</p>
          <ul className="space-y-2 text-sm text-white/80">
            {filteredCommands.map((command) => (
              <li key={command.id}>
                <button
                  type="button"
                  onMouseDown={(event) => {
                    event.preventDefault();
                    applyCommand(command);
                  }}
                  className="flex w-full flex-col gap-1 rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-left transition hover:border-white/20 hover:bg-white/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
                >
                  <span className="text-sm font-semibold">/{command.label.toLowerCase()}</span>
                  <span className="text-xs text-white/60">{command.description}</span>
                </button>
              </li>
            ))}
          </ul>
        </div>
      ) : null}
      <form
        onSubmit={handleSubmit}
        className="sticky bottom-0 flex w-full items-end gap-3 rounded-3xl border border-white/10 bg-white/10 p-3 shadow-[0_12px_40px_rgba(8,12,24,0.35)] backdrop-blur-xl"
      >
        <textarea
          ref={textareaRef}
          value={message}
          onChange={(event) => setMessage(event.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={
            placeholder ??
            "Tanyakan layanan, harga, atau proses kerja tany.ai (ketik / untuk command)"
          }
          rows={1}
          className="max-h-40 flex-1 resize-none bg-transparent text-sm leading-relaxed text-white outline-none placeholder:text-white/40"
          disabled={disabled || isSubmitting}
        />
        <button
          type="submit"
          className="btn-accent inline-flex items-center justify-center rounded-2xl px-5 py-2 text-sm font-semibold uppercase tracking-wide focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400 disabled:cursor-not-allowed disabled:opacity-60"
          disabled={disabled || isSubmitting}
        >
          {isSubmitting ? "Mengirim…" : "Kirim"}
        </button>
      </form>
    </div>
  );
}
