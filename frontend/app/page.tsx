import Image from "next/image";
import Link from "next/link";
import { ArrowUpRight, BadgeCheck, MessageCircle, Sparkles } from "lucide-react";

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
  const heroProject = featuredProjects[0] ?? knowledgeBase.projects[0];
  const heroServices = fallbackServices.slice(0, 2);
  const heroProjectBadges = heroProject
    ? ([
        heroProject.category ? { label: heroProject.category, tone: "neutral" as const } : null,
        heroProject.durationLabel ? { label: `Durasi ${heroProject.durationLabel}`, tone: "neutral" as const } : null,
        heroProject.priceLabel ? { label: heroProject.priceLabel, tone: "accent" as const } : null,
        heroProject.budgetLabel ? { label: heroProject.budgetLabel, tone: "accent" as const } : null,
      ].filter(Boolean) as { label: string; tone: "neutral" | "accent" }[])
    : [];
  const heroTech = heroProject ? heroProject.techStack.slice(0, 3) : [];

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
        <section className="grid gap-12 lg:grid-cols-[minmax(0,1.1fr)_minmax(0,0.9fr)] lg:items-start lg:gap-16">
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
            <div className="flex flex-wrap gap-2">
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
            <div className="flex flex-wrap items-center gap-4">
              <Link
                href="#chat"
                className="btn-accent flex items-center gap-2 rounded-xl px-6 py-4 text-base font-semibold shadow-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
              >
                Mulai chat
                <ArrowUpRight className="h-4 w-4" />
              </Link>
              <Link
                href="#services"
                className="inline-flex items-center gap-2 text-base font-semibold text-white/70 underline underline-offset-4 transition hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
                prefetch
              >
                Lihat layanan
              </Link>
            </div>
            {heroServices.length ? (
              <ul className="grid gap-4 sm:max-w-xl sm:grid-cols-2 lg:max-w-none">
                {heroServices.map((service) => (
                  <li
                    key={service.id}
                    className="flex items-start gap-4 rounded-2xl border border-white/10 bg-white/5 p-4 backdrop-blur-xl"
                  >
                    <span className="mt-0.5 flex h-8 w-8 items-center justify-center rounded-2xl bg-white/10 text-cyan-300">
                      <Sparkles className="h-4 w-4" />
                    </span>
                    <div className="space-y-2">
                      <p className="text-sm font-semibold uppercase tracking-[0.28em] text-white/60">
                        {service.name}
                      </p>
                      {service.description ? (
                        <p className="text-sm leading-relaxed text-white/70">
                          {service.description}
                        </p>
                      ) : null}
                      <div className="flex flex-wrap gap-2 text-xs text-white/60">
                        {service.durationLabel ? (
                          <span className="rounded-full border border-white/15 px-3 py-1">
                            Durasi {service.durationLabel}
                          </span>
                        ) : null}
                        {service.priceRange?.length ? (
                          <span className="rounded-full border border-white/15 px-3 py-1">
                            {service.currency ?? "IDR"} {service.priceRange.join(" – ")}
                          </span>
                        ) : null}
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            ) : null}
          </div>
          <div className="space-y-6">
            {heroProject ? (
              <article className="flex h-full flex-col justify-between gap-6 rounded-3xl border border-white/10 bg-gradient-to-br from-white/10 to-white/5 p-8 shadow-[0_18px_50px_rgba(16,24,48,0.45)] backdrop-blur-xl">
                <div className="space-y-4">
                  <span className="inline-flex items-center gap-2 text-xs uppercase tracking-[0.32em] text-white/60">
                    Highlight proyek
                  </span>
                  <h3 className="font-display text-3xl leading-tight text-white">
                    {heroProject.title}
                  </h3>
                  {heroProject.description ? (
                    <p className="text-base leading-relaxed text-white/70">{heroProject.description}</p>
                  ) : null}
                </div>
                <div className="space-y-4">
                  {heroProjectBadges.length ? (
                    <div className="flex flex-wrap gap-2 text-xs">
                      {heroProjectBadges.map((badge) => (
                        <span
                          key={badge.label}
                          className={
                            badge.tone === "accent"
                              ? "rounded-full bg-emerald-400/15 px-3 py-1 text-emerald-100"
                              : "rounded-full border border-white/15 px-3 py-1 text-white/70"
                          }
                        >
                          {badge.label}
                        </span>
                      ))}
                    </div>
                  ) : null}
                  {heroTech.length ? (
                    <div className="flex flex-wrap gap-2 text-xs text-white/60">
                      {heroTech.map((tech) => (
                        <span key={tech} className="rounded-full border border-white/15 px-3 py-1">
                          {tech}
                        </span>
                      ))}
                    </div>
                  ) : null}
                  {heroProject.projectUrl ? (
                    <Link
                      href={heroProject.projectUrl}
                      className="inline-flex items-center gap-2 text-sm font-semibold text-cyan-300 transition hover:text-cyan-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
                    >
                      Lihat studi kasus
                      <ArrowUpRight className="h-4 w-4" />
                    </Link>
                  ) : null}
                </div>
              </article>
            ) : null}
          </div>
        </section>

        <section id="chat" className="space-y-6">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <span className="text-xs uppercase tracking-[0.4em] text-white/50">Chat assist</span>
              <h2 className="mt-2 font-display text-3xl text-white">Eksplor layanan lewat chat interaktif</h2>
              <p className="mt-2 max-w-2xl text-base text-white/60">
                Tanya paket layanan, estimasi harga, hingga contoh proyek yang relevan. Semua jawaban bersumber
                dari knowledge base Anda.
              </p>
            </div>
            <span className="rounded-full border border-white/15 bg-white/5 px-3 py-1 text-xs text-white/70">
              {updatedLabel}
            </span>
          </div>
          <div className="surface-card relative min-h-[520px] overflow-hidden rounded-3xl border border-white/10 bg-white/5 p-6 shadow-[0_18px_60px_rgba(16,24,48,0.45)] backdrop-blur-xl">
            <ChatWindow initialKnowledge={knowledgeBase} />
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
