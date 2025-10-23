"use client";

import { FormEvent, useState } from "react";

type ChatInputProps = {
  onSend: (message: string) => Promise<void> | void;
  placeholder?: string;
  disabled?: boolean;
};

export function ChatInput({ onSend, placeholder, disabled }: ChatInputProps) {
  const [message, setMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!message.trim()) {
      return;
    }

    try {
      setIsSubmitting(true);
      await onSend(message.trim());
      setMessage("");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="flex w-full items-center gap-3 rounded-2xl border border-slate-200 bg-white p-3 shadow-sm"
    >
      <textarea
        value={message}
        onChange={(event) => setMessage(event.target.value)}
        placeholder={
          placeholder ??
          "Tanyakan apa saja tentang layanan, pengalaman, atau tarif Tanya A.I."
        }
        rows={1}
        className="max-h-40 flex-1 resize-none border-none bg-transparent text-sm leading-relaxed text-slate-900 outline-none placeholder:text-slate-400"
        disabled={disabled || isSubmitting}
      />
      <button
        type="submit"
        className="inline-flex items-center justify-center rounded-full bg-indigo-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-60"
        disabled={disabled || isSubmitting}
      >
        {isSubmitting ? "Mengirim…" : "Kirim"}
      </button>
    </form>
  );
}
