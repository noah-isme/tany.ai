import { ProjectsManager } from "@/components/admin/ProjectsManager";
import { fetchProjects } from "@/lib/admin-api";

import {
  createProjectAction,
  deleteProjectAction,
  featureProjectAction,
  reorderProjectAction,
  updateProjectAction,
} from "./actions";

export const dynamic = "force-dynamic";

export default async function ProjectsPage() {
  const projects = await fetchProjects();

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-100">Proyek</h1>
        <p className="text-sm text-slate-600 dark:text-slate-400">
          Kelola studi kasus yang menjadi referensi utama AI.
        </p>
      </div>
      <ProjectsManager
        initialProjects={projects.items}
        createProject={createProjectAction}
        updateProject={updateProjectAction}
        deleteProject={deleteProjectAction}
        reorderProject={reorderProjectAction}
        featureProject={featureProjectAction}
      />
    </div>
  );
}
