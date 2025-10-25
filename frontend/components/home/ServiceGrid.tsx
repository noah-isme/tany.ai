"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import {
  Code2,
  LucideIcon,
  Palette,
  PenTool,
  Rocket,
  Sparkles,
} from "lucide-react";

import type { KnowledgeService } from "@/lib/knowledge";

const icons: LucideIcon[] = [Sparkles, Code2, Rocket, PenTool, Palette];

const containerVariants = {
  visible: {
    transition: {
      staggerChildren: 0.08,
    },
  },
};

const cardVariants = {
  hidden: { opacity: 0, y: 18 },
  visible: { opacity: 1, y: 0 },
};

type ServiceGridProps = {
  services: KnowledgeService[];
  showActions?: boolean;
  variant?: "compact" | "expanded";
};

export function ServiceGrid({
  services,
  showActions = false,
  variant = "compact",
}: ServiceGridProps) {
  if (!services.length) {
    return null;
  }

  const gridClass =
    variant === "expanded"
      ? "grid gap-5 md:grid-cols-2 xl:grid-cols-3"
      : "grid gap-4 sm:grid-cols-2";

  return (
    <motion.div
      className={gridClass}
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, margin: "-80px" }}
      variants={containerVariants}
    >
      {services.map((service, index) => {
        const Icon = icons[index % icons.length];
        return (
          <motion.article
            key={service.id}
            variants={cardVariants}
            transition={{ duration: 0.28, ease: "easeOut" }}
            whileHover={{ y: -4 }}
            className="group rounded-2xl border border-white/10 bg-white/5 p-5 shadow-[0_10px_30px_rgba(18,23,42,0.25)] backdrop-blur-xl transition focus-within:border-white/20 focus-within:shadow-[0_16px_40px_rgba(124,58,237,0.25)]"
          >
            <div className="flex items-center gap-3">
              <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-white/10 text-white/80">
                <Icon className="h-5 w-5" />
              </span>
              <h3 className="font-semibold text-lg text-white/90">{service.name}</h3>
            </div>
            {service.description ? (
              <p className="mt-3 text-sm leading-relaxed text-white/65">
                {service.description}
              </p>
            ) : null}
            <div className="mt-3 flex flex-wrap gap-2 text-xs text-white/50">
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
            {showActions ? (
              <div className="mt-5 flex items-center justify-between">
                <span className="text-xs uppercase tracking-[0.28em] text-white/40">
                  Detail layanan
                </span>
                <Link
                  href="#chat"
                  className="btn-accent inline-flex items-center gap-2 rounded-lg px-3 py-2 text-xs font-semibold uppercase tracking-wide shadow-[0_10px_26px_rgba(124,58,237,0.2)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
                >
                  Detail
                </Link>
              </div>
            ) : null}
          </motion.article>
        );
      })}
    </motion.div>
  );
}
