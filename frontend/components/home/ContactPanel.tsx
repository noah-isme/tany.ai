import Link from "next/link";
import { Mail, MessageCircle } from "lucide-react";

import type { KnowledgeProfile } from "@/lib/knowledge";

type ContactPanelProps = {
  profile: KnowledgeProfile;
};

function normalizePhoneForWhatsApp(phone?: string): string | null {
  if (!phone) return null;
  const digits = phone.replace(/[^\d+]/g, "");
  if (!digits) return null;
  const normalized = digits.startsWith("+") ? digits.slice(1) : digits;
  return `https://wa.me/${normalized.replace(/^0/, "62")}`;
}

export function ContactPanel({ profile }: ContactPanelProps) {
  const hasEmail = Boolean(profile.email);
  const whatsappLink = normalizePhoneForWhatsApp(profile.phone);

  if (!hasEmail && !whatsappLink) {
    return null;
  }

  return (
    <div className="fixed bottom-8 right-8 z-40 hidden max-w-sm flex-col gap-3 rounded-3xl border border-white/10 bg-white/10 p-4 shadow-[0_20px_45px_rgba(8,12,24,0.35)] backdrop-blur-xl lg:flex">
      <span className="text-xs uppercase tracking-[0.32em] text-white/60">Hubungi langsung</span>
      <div className="flex flex-wrap gap-3">
        {hasEmail ? (
          <Link
            href={`mailto:${profile.email}`}
            className="btn-accent inline-flex flex-1 items-center justify-center gap-2 rounded-xl px-3 py-2 text-xs font-semibold uppercase tracking-wide focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
          >
            <Mail className="h-4 w-4" /> Email
          </Link>
        ) : null}
        {whatsappLink ? (
          <Link
            href={whatsappLink}
            className="inline-flex flex-1 items-center justify-center gap-2 rounded-xl border border-white/20 px-3 py-2 text-xs font-semibold uppercase tracking-wide text-white/80 transition hover:bg-white/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
            target="_blank"
            rel="noopener noreferrer"
          >
            <MessageCircle className="h-4 w-4" /> WhatsApp
          </Link>
        ) : null}
      </div>
    </div>
  );
}
