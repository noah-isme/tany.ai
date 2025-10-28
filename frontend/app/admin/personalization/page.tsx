import { PersonalizationPanel } from "@/components/admin/PersonalizationPanel";
import { fetchPersonalizationSummary } from "@/lib/admin-api";

import {
  reindexPersonalizationAction,
  resetPersonalizationAction,
  updatePersonalizationWeightAction,
} from "./actions";

export const dynamic = "force-dynamic";

export default async function PersonalizationPage() {
  const summary = await fetchPersonalizationSummary();

  return (
    <div className="space-y-6">
      <header className="space-y-2">
        <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-100">AI Personalization</h1>
        <p className="text-sm text-slate-600 dark:text-slate-400">
          Atur embedding persona, bobot gaya tulisan, dan lakukan reindex untuk menjaga relevansi respons Tany.AI.
        </p>
      </header>

      <PersonalizationPanel
        summary={summary}
        updateWeight={updatePersonalizationWeightAction}
        reindex={reindexPersonalizationAction}
        reset={resetPersonalizationAction}
      />
    </div>
  );
}
