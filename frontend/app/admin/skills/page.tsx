import { SkillsManager } from "@/components/admin/SkillsManager";
import { fetchSkills } from "@/lib/admin-api";

import {
  createSkillAction,
  deleteSkillAction,
  reorderSkillAction,
  updateSkillAction,
} from "./actions";

export const dynamic = "force-dynamic";

export default async function SkillsPage() {
  const skills = await fetchSkills();

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-100">Skills</h1>
        <p className="text-sm text-slate-600 dark:text-slate-400">
          Urutkan skill sesuai prioritas agar jawaban AI tetap konsisten.
        </p>
      </div>
      <SkillsManager
        initialSkills={skills.items}
        createSkill={createSkillAction}
        updateSkill={updateSkillAction}
        deleteSkill={deleteSkillAction}
        reorderSkill={reorderSkillAction}
      />
    </div>
  );
}
