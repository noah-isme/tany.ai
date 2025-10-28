"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { ArrowUpRight, Sparkles } from "lucide-react";

import type { KnowledgeProject } from "@/lib/knowledge";

type PortfolioShowcaseProps = {
  projects: KnowledgeProject[];
};

const containerVariants = {
  visible: {
    transition: {
      staggerChildren: 0.1,
    },
  },
};

const cardVariants = {
  hidden: { opacity: 0, y: 24 },
  visible: { opacity: 1, y: 0 },
};

export function PortfolioShowcase({ projects }: PortfolioShowcaseProps) {
  if (!projects.length) {
    return null;
  }

  return (
    <motion.div
      className="grid gap-6 md:grid-cols-2 xl:gap-8"
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, margin: "-120px" }}
      variants={containerVariants}
    >
      {projects.map((project, index) => {
        const isFeatured = project.isFeatured && index < 2;
        const metadata = [
          project.category ? { label: project.category, tone: "neutral" } : null,
          project.durationLabel ? { label: `Durasi ${project.durationLabel}`, tone: "neutral" } : null,
          project.priceLabel ? { label: project.priceLabel, tone: "accent" } : null,
          project.budgetLabel ? { label: project.budgetLabel, tone: "accent" } : null,
        ].filter(Boolean) as { label: string; tone: "neutral" | "accent" }[];

        return (
          <motion.article
            key={project.id}
            variants={cardVariants}
            transition={{ duration: 0.3, ease: "easeOut" }}
            className={`group relative flex h-full flex-col justify-between overflow-hidden rounded-3xl border border-white/10 bg-gradient-to-br from-white/12 via-white/5 to-white/4 p-8 shadow-[0_20px_40px_rgba(8,12,24,0.35)] backdrop-blur-xl transition duration-300 hover:-translate-y-2 hover:shadow-[0_28px_66px_rgba(124,58,237,0.25)] focus-within:border-white/20 focus-within:shadow-[0_28px_66px_rgba(124,58,237,0.25)] ${
              isFeatured ? "md:col-span-2 md:p-10" : ""
            }`}
          >
            <div className="relative z-10 space-y-6">
              {isFeatured ? (
                <span className="inline-flex items-center gap-2 rounded-full bg-white/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.32em] text-white/70">
                  <Sparkles className="h-4 w-4" /> Featured
                </span>
              ) : null}
              <div className="space-y-3">
                <h3 className="font-display text-[1.75rem] leading-tight text-white sm:text-[1.9rem]">
                  {project.title}
                </h3>
                {project.description ? (
                  <p className="text-base leading-relaxed text-white/70">
                    {project.description}
                  </p>
                ) : null}
              </div>
              {metadata.length ? (
                <div className="flex flex-wrap gap-2 text-xs font-semibold uppercase tracking-[0.24em]">
                  {metadata.map((item) => (
                    <span
                      key={item.label}
                      className={`rounded-full px-3 py-1 ${
                        item.tone === "accent"
                          ? "bg-cyan-400/15 text-cyan-100"
                          : "border border-white/15 text-white/60"
                      }`}
                    >
                      {item.label}
                    </span>
                  ))}
                </div>
              ) : null}
              {project.techStack.length ? (
                <div className="flex flex-wrap gap-2 text-xs text-white/45">
                  {project.techStack.map((tech) => (
                    <span
                      key={tech}
                      className="rounded-full border border-white/15 bg-white/5 px-3 py-1 uppercase tracking-wide"
                    >
                      {tech}
                    </span>
                  ))}
                </div>
              ) : null}
            </div>
            {project.projectUrl ? (
              <Link
                href={project.projectUrl}
                className="relative z-10 mt-6 inline-flex w-fit items-center gap-2 rounded-full border border-white/20 bg-white/5 px-4 py-2 text-xs font-semibold uppercase tracking-wide text-white/85 transition hover:border-white/40 hover:bg-white/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
              >
                Lihat detail proyek
                <ArrowUpRight className="h-4 w-4" />
              </Link>
            ) : null}
            <div className="pointer-events-none absolute inset-0 opacity-0 transition duration-300 group-hover:opacity-100">
              <div className="absolute inset-0 bg-gradient-to-r from-violet-500/20 via-transparent to-cyan-400/20" />
            </div>
          </motion.article>
        );
      })}
    </motion.div>
  );
}
