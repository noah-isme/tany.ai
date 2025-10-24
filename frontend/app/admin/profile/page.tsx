import { ProfileForm } from "@/components/admin/ProfileForm";
import { fetchProfile } from "@/lib/admin-api";

import { updateProfileAction } from "./actions";

export const dynamic = "force-dynamic";

export default async function ProfilePage() {
  const profile = await fetchProfile();

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-100">Profil & Kontak</h1>
        <p className="text-sm text-slate-600 dark:text-slate-400">
          Data ini menjadi identitas utama yang digunakan AI saat menjawab calon klien.
        </p>
      </div>
      <ProfileForm profile={profile} onSubmit={updateProfileAction} />
    </div>
  );
}
