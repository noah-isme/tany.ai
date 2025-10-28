import { ExternalIntegrationView } from "@/components/admin/ExternalIntegrationView";
import { fetchExternalItems, fetchExternalSources } from "@/lib/admin-api";

import {
  syncExternalSourceAction,
  toggleExternalItemVisibilityAction,
} from "./actions";

export const dynamic = "force-dynamic";

export default async function IntegrationsPage() {
  const [sources, items] = await Promise.all([
    fetchExternalSources(),
    fetchExternalItems(),
  ]);

  return (
    <ExternalIntegrationView
      initialSources={sources.items}
      initialItems={items.items}
      syncSource={syncExternalSourceAction}
      toggleItem={toggleExternalItemVisibilityAction}
    />
  );
}
