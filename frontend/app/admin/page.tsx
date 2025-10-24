import Link from "next/link";

import {
  fetchProfile,
  fetchProjects,
  fetchServices,
  fetchSkills,
} from "@/lib/admin-api";

export default async function AdminDashboardPage() {
  const [profile, skills, services, projects] = await Promise.all([
    fetchProfile(),
    fetchSkills(),
    fetchServices(),
    fetchProjects(),
  ]);

  const profileUpdatedAt = profile.updated_at ? new Date(profile.updated_at) : new Date();
  const profileUpdatedAtLabel = profileUpdatedAt.toLocaleDateString("id-ID");

  return (
    <div className="space-y-8">
      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <DashboardCard title="Skill Aktif" value={skills.items.length} href="/admin/skills" />
        <DashboardCard title="Layanan" value={services.items.length} href="/admin/services" />
        <DashboardCard title="Proyek" value={projects.items.length} href="/admin/projects" />
        <DashboardCard title="Profil Diperbarui" value={profileUpdatedAtLabel} href="/admin/profile" />
      </section>

      <section className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
          <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Ringkasan Profil</h2>
          <dl className="mt-4 space-y-2 text-sm text-slate-600 dark:text-slate-300">
            <div className="flex justify-between">
              <dt>Nama</dt>
              <dd className="font-medium text-slate-900 dark:text-slate-100">{profile.name}</dd>
            </div>
            <div className="flex justify-between">
              <dt>Jabatan</dt>
              <dd>{profile.title}</dd>
            </div>
            <div className="flex justify-between">
              <dt>Email</dt>
              <dd>{profile.email || "-"}</dd>
            </div>
            <div className="flex justify-between">
              <dt>Lokasi</dt>
              <dd>{profile.location || "-"}</dd>
            </div>
          </dl>
          <Link
            href="/admin/profile"
            className="mt-4 inline-flex text-sm font-semibold text-indigo-500 hover:text-indigo-400"
          >
            Kelola profil →
          </Link>
        </div>

        <div className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
          <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Layanan Populer</h2>
          <ul className="mt-4 space-y-3 text-sm text-slate-600 dark:text-slate-300">
            {services.items.slice(0, 3).map((service) => (
              <li key={service.id} className="rounded-lg border border-slate-200/60 bg-white/60 p-3 dark:border-slate-800/60 dark:bg-slate-950/40">
                <p className="font-semibold text-slate-900 dark:text-slate-100">{service.name}</p>
                <p className="text-xs text-slate-500 dark:text-slate-400">
                  {service.duration_label || "Durasi menyesuaikan"}
                </p>
              </li>
            ))}
          </ul>
          <Link
            href="/admin/services"
            className="mt-4 inline-flex text-sm font-semibold text-indigo-500 hover:text-indigo-400"
          >
            Atur layanan →
          </Link>
        </div>
      </section>

      <section className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
        <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Proyek Terbaru</h2>
        <ul className="mt-4 space-y-3 text-sm text-slate-600 dark:text-slate-300">
          {projects.items.slice(0, 4).map((project) => (
            <li key={project.id} className="rounded-lg border border-slate-200/60 bg-white/60 p-3 dark:border-slate-800/60 dark:bg-slate-950/40">
              <p className="font-semibold text-slate-900 dark:text-slate-100">{project.title}</p>
              <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">{project.category || "Tanpa kategori"}</p>
            </li>
          ))}
        </ul>
        <Link
          href="/admin/projects"
          className="mt-4 inline-flex text-sm font-semibold text-indigo-500 hover:text-indigo-400"
        >
          Kelola proyek →
        </Link>
      </section>
    </div>
  );
}

type DashboardCardProps = {
  title: string;
  value: string | number;
  href: string;
};

function DashboardCard({ title, value, href }: DashboardCardProps) {
  return (
    <Link
      href={href}
      className="group rounded-2xl border border-slate-200 bg-white/80 p-5 text-left shadow-sm transition hover:-translate-y-1 hover:border-indigo-400 hover:shadow-lg dark:border-slate-800 dark:bg-slate-900/70"
    >
      <p className="text-xs uppercase tracking-[0.25em] text-slate-400 group-hover:text-indigo-300">{title}</p>
      <p className="mt-4 text-2xl font-semibold text-slate-900 dark:text-slate-100">{value}</p>
      <p className="mt-4 text-sm font-semibold text-indigo-500 group-hover:text-indigo-400">Lihat detail →</p>
    </Link>
  );
}
