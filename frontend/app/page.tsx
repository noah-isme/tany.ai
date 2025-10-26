import Image from "next/image";
import Link from "next/link";
import { ArrowUpRight, BadgeCheck, MessageCircle } from "lucide-react";

import { ChatWindow } from "@/components/chat/ChatWindow";
import { ContactPanel } from "@/components/home/ContactPanel";
import { PortfolioShowcase } from "@/components/home/PortfolioShowcase";
import { ServiceGrid } from "@/components/home/ServiceGrid";
import { fetchKnowledgeBase } from "@/lib/knowledge";

export const revalidate = 0;

function initialsFromName(name: string): string {
  if (!name) return "TA";
  const tokens = name.trim().split(/\s+/);
  if (tokens.length === 1) {
    return tokens[0].slice(0, 2).toUpperCase();
  }
  return (tokens[0][0] + tokens[1][0]).toUpperCase();
}

export default async function HomePage() {
  const knowledgeBase = await fetchKnowledgeBase();
  const featuredServices = knowledgeBase.services
    .filter((service) => service.description)
    .slice(0, 4);
  const fallbackServices =
    featuredServices.length > 0 ? featuredServices : knowledgeBase.services.slice(0, 4);
  const featuredProjects = knowledgeBase.projects
    .filter((project) => project.isFeatured)
    .concat(knowledgeBase.projects)
    .filter((project, index, array) => array.findIndex((item) => item.id === project.id) === index)
    .slice(0, 4);

  const profile = knowledgeBase.profile;
  const subtitle = profile.title ?? "AI Client Partner";
  const updatedLabel = profile.updatedAt
    ? new Date(profile.updatedAt).toLocaleDateString("id-ID", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      })
    : "Realtime";

  return (
    <main className="bg-app">
      <div className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-20 px-6 py-16 lg:px-10">
        <section className="grid gap-14 lg:grid-cols-[1.05fr_0.95fr] lg:items-start">
          <div className="space-y-8">
            <div className="flex items-center gap-5">
              <div className="relative h-20 w-20 overflow-hidden rounded-3xl border border-white/10 bg-white/10">
                {profile.avatarUrl ? (
                  <Image
                    src={profile.avatarUrl}
                    alt={profile.name}
                    fill
                    className="object-cover"
                    sizes="80px"
                    priority
                  />
                ) : (
                  <div className="flex h-full w-full items-center justify-center text-xl font-semibold text-white/80">
                    {initialsFromName(profile.name)}
                  </div>
                )}
              </div>
              <div>
                <span className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.28em] text-white/70">
                  <BadgeCheck className="h-3.5 w-3.5" /> {subtitle}
                </span>
                <h1 className="mt-4 max-w-xl font-display text-4xl leading-tight tracking-tight text-white sm:text-5xl">
                  {profile.name || "tany.ai"}
                </h1>
              </div>
            </div>
            <p className="max-w-xl text-lg leading-relaxed text-white/70">
              {profile.bio ??
                "Antarmuka chat yang grounded pada knowledge base pribadi Anda. Jawaban terasa manusiawi, relevan, dan dapat mengonversi prospek lebih cepat."}
            </p>
            <div className="flex flex-wrap gap-3">
              {profile.location ? (
                <span className="inline-flex items-center gap-2 rounded-full bg-white/5 px-4 py-2 text-sm text-white/80">
                  <MessageCircle className="h-4 w-4" /> {profile.location}
                </span>
              ) : null}
              <span className="inline-flex items-center gap-2 rounded-full bg-white/5 px-4 py-2 text-sm text-white/80">
                Real-time AI chat
              </span>
              <span className="inline-flex items-center gap-2 rounded-full bg-white/5 px-4 py-2 text-sm text-white/80">
                Knowledge base terkurasi
              </span>
            </div>
            <div className="flex flex-wrap gap-4">
              <Link
                href="#chat"
                className="btn-accent flex items-center gap-2 rounded-xl px-5 py-3 text-sm font-semibold shadow-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
              >
                Mulai chat
                <ArrowUpRight className="h-4 w-4" />
              </Link>
              <Link
                href="#services"
                className="inline-flex items-center gap-2 rounded-xl border border-white/20 px-5 py-3 text-sm font-semibold text-white/80 transition hover:bg-white/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
                prefetch
              >
                Lihat layanan
              </Link>
            </div>
            <ServiceGrid services={fallbackServices} />
          </div>
          <div id="chat" className="relative flex h-full flex-col">
            <div className="surface-card relative flex min-h-[520px] flex-1 flex-col overflow-hidden rounded-3xl p-6 shadow-[0_18px_60px_rgba(16,24,48,0.45)]">
              <div className="mb-4 flex items-center justify-between">
                <div>
                  <p className="text-xs uppercase tracking-[0.32em] text-white/50">Chat Assist</p>
                  <p className="mt-1 text-sm text-white/70">Tanya layanan, harga, hingga contoh proyek.</p>
                </div>
                <span className="rounded-full border border-white/15 bg-white/5 px-3 py-1 text-xs text-white/70">
                  {updatedLabel}
                </span>
              </div>
              <ChatWindow initialKnowledge={knowledgeBase} />
            </div>
          </div>
        </section>

        <section id="services" className="space-y-8">
          <div className="flex flex-col gap-3">
            <span className="text-xs uppercase tracking-[0.4em] text-white/50">Layanan utama</span>
            <h2 className="font-display text-3xl text-white">Solusi paling dicari klien</h2>
            <p className="max-w-2xl text-base text-white/60">
              Paket layanan disusun agar ringkas dan mudah dimengerti calon klien. Setiap kartu memuat ringkasan manfaat dan ajakan bertindak ke detail lengkap.
            </p>
          </div>
          <ServiceGrid services={fallbackServices} showActions variant="expanded" />
        </section>

        <section className="space-y-8">
          <div className="flex flex-col gap-3">
            <span className="text-xs uppercase tracking-[0.4em] text-white/50">Portofolio</span>
            <h2 className="font-display text-3xl text-white">Jejak proyek yang relevan</h2>
            <p className="max-w-2xl text-base text-white/60">
              Dari produk SaaS hingga landing page kampanye, highlight proyek berikut menunjukkan kedalaman kolaborasi dan teknologi yang digunakan.
            </p>
          </div>
          <PortfolioShowcase projects={featuredProjects} />
        </section>
      </div>
      <ContactPanel profile={profile} />
    </main>
  );
}
