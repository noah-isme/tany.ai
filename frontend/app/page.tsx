import Link from "next/link";

import { LogoMark } from "@/components/LogoMark";
import { ChatWindow } from "@/components/chat/ChatWindow";
import { buildSystemPrompt, knowledgeBase } from "@/data/knowledge";

const systemPrompt = buildSystemPrompt(knowledgeBase);

export default function Home() {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <div className="mx-auto grid max-w-6xl gap-12 px-6 py-16 lg:grid-cols-[1.15fr_minmax(0,1fr)] lg:gap-16 lg:px-8 lg:py-24">
        <section className="space-y-10">
          <div className="flex items-center gap-4">
            <LogoMark />
            <div className="space-y-1">
              <p className="text-lg font-semibold text-white">tany.ai</p>
              <span className="inline-flex items-center gap-2 rounded-full border border-indigo-400/40 bg-indigo-500/10 px-4 py-1 text-[0.65rem] font-semibold uppercase tracking-[0.35em] text-indigo-200">
                MVP Fokus · Chat Assistant Pribadi
              </span>
            </div>
          </div>
          <header className="space-y-6">
            <h1 className="text-4xl font-semibold leading-tight text-white sm:text-5xl">
              tany.ai membantu Anda menjawab calon klien secara instan, profesional, dan 24/7.
            </h1>
            <p className="max-w-2xl text-base leading-relaxed text-slate-300">
              Prototipe awal ini menyatukan pondasi frontend Next.js dan backend Golang sesuai panduan roadmap. Fokusnya adalah
              menghadirkan antarmuka chat modern yang terhubung ke knowledge base personal sehingga setiap jawaban terasa relevan
              dan kredibel.
            </p>
          </header>

          <div className="grid gap-6 sm:grid-cols-2">
            {knowledgeBase.services.map((service) => (
              <article
                key={service.name}
                className="rounded-3xl border border-white/10 bg-white/5 p-5 shadow-lg shadow-indigo-900/30 backdrop-blur"
              >
                <h2 className="text-lg font-semibold text-white">{service.name}</h2>
                <p className="mt-2 text-sm text-slate-300">{service.description}</p>
                <p className="mt-4 text-xs font-semibold uppercase tracking-wide text-indigo-200">
                  Mulai {service.startingAt} · {service.turnaround}
                </p>
              </article>
            ))}
          </div>

          <div className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-inner shadow-black/30">
            <h2 className="text-xl font-semibold text-white">Mengapa calon klien menyukai tany.ai?</h2>
            <ul className="mt-4 grid gap-3 text-sm text-slate-300">
              <li>• Jawaban konsisten berdasarkan profil, layanan, dan portofolio Anda.</li>
              <li>• Konteks percakapan disimpan agar diskusi terasa natural.</li>
              <li>• Prompt engineering memastikan AI tetap on-brand dan tidak berimprovisasi berlebihan.</li>
              <li>• Tersedia log percakapan untuk analitik dan follow-up lead.</li>
            </ul>
          </div>

          <details className="group rounded-3xl border border-white/10 bg-slate-900/60 p-6 text-sm text-slate-200">
            <summary className="cursor-pointer text-base font-semibold text-white">
              Lihat system prompt yang men-drive respons AI
            </summary>
            <pre className="mt-4 max-h-80 overflow-y-auto whitespace-pre-wrap rounded-2xl bg-black/60 p-4 text-xs text-slate-200">
              {systemPrompt}
            </pre>
          </details>
        </section>

        <section className="flex flex-col gap-6 lg:sticky lg:top-12">
          <div className="rounded-3xl border border-white/10 bg-slate-900/70 p-6 shadow-xl shadow-black/40 ring-1 ring-white/10 backdrop-blur">
            <div className="mb-5 space-y-2">
              <h2 className="text-xl font-semibold text-white">Simulasi percakapan tany.ai</h2>
              <p className="text-sm text-slate-300">
                Ketik pertanyaan tentang layanan, tarif, atau pengalaman. Saat ini jawaban masih mock response, namun struktur API
                sudah siap terhubung dengan OpenAI GPT sesuai roadmap.
              </p>
            </div>
            <ChatWindow />
          </div>

          <div className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-lg shadow-black/20">
            <h3 className="text-lg font-semibold text-white">Hubungi langsung</h3>
            <ul className="mt-4 space-y-2 text-sm text-slate-300">
              <li>
                Email: <Link href={`mailto:${knowledgeBase.contact.email}`} className="text-indigo-200 underline">
                  {knowledgeBase.contact.email}
                </Link>
              </li>
              <li>
                Website: <Link href={knowledgeBase.contact.website} className="text-indigo-200 underline">
                  {knowledgeBase.contact.website.replace("https://", "")}
                </Link>
              </li>
              <li>WhatsApp: {knowledgeBase.contact.whatsApp}</li>
              <li>
                LinkedIn: <Link href={knowledgeBase.contact.linkedIn} className="text-indigo-200 underline">
                  {knowledgeBase.contact.linkedIn.replace("https://", "")}
                </Link>
              </li>
            </ul>
            <p className="mt-4 text-xs uppercase tracking-wide text-indigo-200">
              Dibangun dengan Next.js 16 · Tailwind CSS 4 · Terhubung ke backend Golang + Gin
            </p>
          </div>
        </section>
      </div>
    </div>
  );
}
