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
      className="grid gap-5 md:grid-cols-2"
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, margin: "-120px" }}
      variants={containerVariants}
    >
      {projects.map((project, index) => {
        const isFeatured = project.isFeatured && index < 2;
        return (
          <motion.article
            key={project.id}
            variants={cardVariants}
            transition={{ duration: 0.3, ease: "easeOut" }}
            className={`group relative overflow-hidden rounded-3xl border border-white/10 bg-gradient-to-br from-white/10 to-white/5 p-6 shadow-[0_20px_40px_rgba(8,12,24,0.35)] backdrop-blur-xl transition hover:shadow-[0_26px_60px_rgba(124,58,237,0.25)] focus-within:border-white/20 focus-within:shadow-[0_26px_60px_rgba(124,58,237,0.25)] ${
              isFeatured ? "md:col-span-2 md:grid md:grid-cols-[1.1fr_0.9fr] md:items-center md:gap-8" : ""
            }`}
          >
            <div className="relative z-10 space-y-4">
              {isFeatured ? (
                <span className="inline-flex items-center gap-2 rounded-full bg-white/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.32em] text-white/70">
                  <Sparkles className="h-4 w-4" /> Featured
                </span>
              ) : null}
              <div>
                <h3 className="font-display text-2xl text-white">{project.title}</h3>
                {project.description ? (
                  <p className="mt-2 text-sm leading-relaxed text-white/65">
                    {project.description}
                  </p>
                ) : null}
              </div>
              {project.techStack.length ? (
                <div className="flex flex-wrap gap-2 text-xs uppercase tracking-wide text-white/40">
                  {project.techStack.map((tech) => (
                    <span
                      key={tech}
                      className="rounded-full border border-white/15 bg-white/5 px-3 py-1"
                    >
                      {tech}
                    </span>
                  ))}
                </div>
              ) : null}
              {project.projectUrl ? (
                <Link
                  href={project.projectUrl}
                  className="btn-accent inline-flex items-center gap-2 rounded-lg px-3 py-2 text-xs font-semibold uppercase tracking-wide shadow-[0_10px_26px_rgba(124,58,237,0.2)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
                >
                  Lihat proyek
                  <ArrowUpRight className="h-4 w-4" />
                </Link>
              ) : null}
            </div>
            <div className="pointer-events-none absolute inset-0 opacity-0 transition duration-300 group-hover:opacity-100">
              <div className="absolute inset-0 bg-gradient-to-r from-violet-500/20 via-transparent to-cyan-400/20" />
            </div>
          </motion.article>
        );
      })}
    </motion.div>
  );
}
