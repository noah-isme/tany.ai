import { ServicesManager } from "@/components/admin/ServicesManager";
import { fetchServices } from "@/lib/admin-api";

import {
  createServiceAction,
  deleteServiceAction,
  reorderServiceAction,
  toggleServiceAction,
  updateServiceAction,
} from "./actions";

export const dynamic = "force-dynamic";

export default async function ServicesPage() {
  const services = await fetchServices();

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-100">Layanan</h1>
        <p className="text-sm text-slate-600 dark:text-slate-400">
          Atur layanan, harga, dan status agar AI memahami penawaran Anda.
        </p>
      </div>
      <ServicesManager
        initialServices={services.items}
        createService={createServiceAction}
        updateService={updateServiceAction}
        deleteService={deleteServiceAction}
        reorderService={reorderServiceAction}
        toggleService={toggleServiceAction}
      />
    </div>
  );
}
