import Link from "next/link";

import { ChatWindow } from "@/components/chat/ChatWindow";
import { fetchKnowledgeBase } from "@/lib/knowledge";

export const revalidate = 0;

export default async function HomePage() {
  const knowledgeBase = await fetchKnowledgeBase();
  const featuredServices = knowledgeBase.services.slice(0, 3);
  const featuredProjects = knowledgeBase.projects
    .filter((project) => project.isFeatured)
    .concat(knowledgeBase.projects)
    .slice(0, 3);

  return (
    <main className="mx-auto flex min-h-screen max-w-6xl flex-col gap-16 px-6 py-12">
      <section className="grid gap-12 lg:grid-cols-2 lg:items-start">
        <div className="space-y-6">
          <h1 className="text-4xl font-bold tracking-tight text-slate-900 lg:text-5xl">
            {knowledgeBase.profile.name || "tany.ai"}
          </h1>
          <p className="text-lg leading-relaxed text-slate-600">
            tany.ai menghadirkan antarmuka chat modern yang terhubung ke knowledge
            base personal sehingga setiap jawaban terasa relevan dengan layanan
            Anda.
          </p>
          <ul className="grid gap-3 text-sm text-slate-600 sm:grid-cols-2">
            {featuredServices.map((service) => (
              <li
                key={service.id}
                className="rounded-2xl border border-slate-200 bg-white/80 p-4 shadow-sm"
              >
                <h3 className="font-semibold text-slate-900">{service.name}</h3>
                <p className="mt-2 text-xs leading-relaxed text-slate-500">
                  {service.description}
                </p>
              </li>
            ))}
          </ul>
          <div className="rounded-2xl border border-slate-200 bg-slate-900 p-6 text-slate-100">
            <h2 className="text-sm font-semibold uppercase tracking-wide text-indigo-200">
              Kontak
            </h2>
            <ul className="mt-3 space-y-2 text-sm">
              {knowledgeBase.profile.email ? (
                <li>
                  Email: <Link href={`mailto:${knowledgeBase.profile.email}`} className="text-indigo-200 underline">
                    {knowledgeBase.profile.email}
                  </Link>
                </li>
              ) : null}
              {knowledgeBase.profile.phone ? <li>WhatsApp: {knowledgeBase.profile.phone}</li> : null}
            </ul>
          </div>
        </div>
        <div className="min-h-[480px] rounded-3xl border border-slate-200 bg-white/80 p-6 shadow-xl">
          <ChatWindow />
        </div>
      </section>
      <section className="space-y-6">
        <h2 className="text-2xl font-semibold text-slate-900">Portofolio</h2>
        <div className="grid gap-6 md:grid-cols-2">
          {featuredProjects.map((project) => (
            <article
              key={project.id}
              className="rounded-3xl border border-slate-200 bg-white/80 p-6 shadow-sm"
            >
              <h3 className="text-lg font-semibold text-slate-900">
                {project.title}
              </h3>
              <p className="mt-2 text-sm text-slate-600">{project.description}</p>
              {project.techStack.length ? (
                <p className="mt-2 text-xs uppercase tracking-wide text-slate-400">
                  {project.techStack.join(" • ")}
                </p>
              ) : null}
              {project.projectUrl ? (
                <Link
                  href={project.projectUrl}
                  className="mt-3 inline-block text-sm font-medium text-indigo-600"
                >
                  Lihat proyek
                </Link>
              ) : null}
            </article>
          ))}
        </div>
      </section>
    </main>
  );
}
