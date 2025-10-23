export type Profile = {
  name: string;
  tagline: string;
  bio: string;
  expertise: string[];
  yearsActive: number;
  location: string;
};

export type Service = {
  name: string;
  description: string;
  deliverables: string[];
  turnaround: string;
  startingAt: string;
};

export type Project = {
  title: string;
  summary: string;
  techStack: string[];
  impact: string;
};

export type PricingTier = {
  name: string;
  price: string;
  whatYouGet: string[];
};

export type ContactInfo = {
  email: string;
  website: string;
  whatsApp: string;
  linkedIn: string;
};

export type KnowledgeBase = {
  profile: Profile;
  services: Service[];
  projects: Project[];
  pricing: PricingTier[];
  contact: ContactInfo;
};

export const knowledgeBase: KnowledgeBase = {
  profile: {
    name: "Tanya A.I.",
    tagline: "AI Client Assistant for Modern Freelancers",
    bio: "Saya membantu menjawab pertanyaan calon klien secara profesional dengan konteks data freelancer.",
    expertise: ["Next.js", "Golang", "AI Prompt Engineering", "UI/UX Strategy"],
    yearsActive: 7,
    location: "Jakarta, Indonesia",
  },
  services: [
    {
      name: "Pembuatan Website Portfolio",
      description:
        "Website responsif dengan fokus pada storytelling brand personal.",
      deliverables: [
        "Landing page",
        "Halaman layanan",
        "Integrasi formulir leads",
      ],
      turnaround: "2-3 minggu",
      startingAt: "IDR 12.000.000",
    },
    {
      name: "AI Chat Assistant Custom",
      description:
        "Implementasi chatbot berbasis profil profesional untuk otomatisasi komunikasi klien.",
      deliverables: [
        "Setup knowledge base",
        "Integrasi OpenAI GPT",
        "Dashboard monitoring",
      ],
      turnaround: "3-4 minggu",
      startingAt: "IDR 18.000.000",
    },
    {
      name: "Consulting & Strategy",
      description:
        "Sesi konsultasi 1:1 mengenai positioning layanan digital dan optimasi workflow AI.",
      deliverables: [
        "Audit profil",
        "Rencana peningkatan layanan",
        "Template komunikasi klien",
      ],
      turnaround: "1 minggu",
      startingAt: "IDR 2.500.000",
    },
  ],
  projects: [
    {
      title: "SaaS Landing Page untuk Startup Fintech",
      summary:
        "Membangun landing page interaktif dengan konversi tinggi untuk peluncuran produk fintech.",
      techStack: ["Next.js", "Tailwind CSS", "Vercel"],
      impact: "Meningkatkan sign-up calon pengguna sebesar 32% dalam 3 minggu pertama.",
    },
    {
      title: "Chat Assistant untuk Agency Kreatif",
      summary:
        "Membuat chatbot khusus untuk menjawab pertanyaan layanan desain dan motion graphics.",
      techStack: ["Golang", "MongoDB", "OpenAI GPT-4"],
      impact: "Mengurangi waktu respons tim sales dari 6 jam menjadi 15 menit rata-rata.",
    },
    {
      title: "Knowledge Hub Freelancer",
      summary:
        "Portal internal untuk mengelola portfolio, layanan, dan studi kasus.",
      techStack: ["Next.js", "Supabase", "Clerk"],
      impact: "Mempercepat proses update portfolio bulanan dari 3 hari menjadi 1 hari kerja.",
    },
  ],
  pricing: [
    {
      name: "Starter",
      price: "IDR 7.500.000",
      whatYouGet: [
        "Landing page 1 halaman",
        "Optimasi konten profil",
        "Integrasi formulir kontak",
      ],
    },
    {
      name: "Professional",
      price: "IDR 15.000.000",
      whatYouGet: [
        "Website multi-halaman",
        "Integrasi blog & CMS",
        "Chat assistant dasar",
      ],
    },
    {
      name: "Elite",
      price: "Mulai IDR 25.000.000",
      whatYouGet: [
        "AI chat assistant custom",
        "Dashboard analitik",
        "Onboarding & training",
      ],
    },
  ],
  contact: {
    email: "hello@tany.ai",
    website: "https://tany.ai",
    whatsApp: "+62-811-2345-678",
    linkedIn: "https://linkedin.com/in/tanyai",
  },
};

export function buildSystemPrompt(base: KnowledgeBase): string {
  const services = base.services
    .map(
      (service) =>
        `- ${service.name} (mulai ${service.startingAt}): ${service.description}. Deliverables: ${service.deliverables.join(", ")}`,
    )
    .join("\n");

  const projects = base.projects
    .map(
      (project) =>
        `- ${project.title} menggunakan ${project.techStack.join(", ")} → ${project.impact}`,
    )
    .join("\n");

  const pricing = base.pricing
    .map(
      (tier) =>
        `- Paket ${tier.name} (${tier.price}): ${tier.whatYouGet.join(", ")}`,
    )
    .join("\n");

  return `Kamu adalah asisten virtual untuk ${base.profile.name}.\n\nProfil:\n- Tagline: ${base.profile.tagline}\n- Bio: ${base.profile.bio}\n- Keahlian: ${base.profile.expertise.join(", ")}\n- Pengalaman: ${base.profile.yearsActive} tahun\n- Lokasi: ${base.profile.location}\n\nLayanan:\n${services}\n\nPortfolio Utama:\n${projects}\n\nStruktur Harga:\n${pricing}\n\nKontak:\n- Email: ${base.contact.email}\n- Website: ${base.contact.website}\n- WhatsApp: ${base.contact.whatsApp}\n- LinkedIn: ${base.contact.linkedIn}\n\nAturan:\n1. Hanya gunakan informasi di atas.\n2. Jawab dengan nada ramah dan profesional.\n3. Jika pertanyaan di luar scope, arahkan klien ke email.`;
}

export function createMockAnswer(question: string, base: KnowledgeBase): string {
  return `Pertanyaan diterima: ${question}\n\nHalo! Saya ${base.profile.name}. Saya menawarkan ${base.services
    .map((service) => service.name)
    .join(", ")}. Contoh proyek terbaru saya adalah ${base.projects[0].title}.\nSilakan hubungi saya di ${base.contact.email} untuk diskusi lebih lanjut.`;
}
