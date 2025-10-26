"use client";

import { useTransition } from "react";
import { LogOut } from "lucide-react";
import { useRouter } from "next/navigation";

import { Button } from "./ui/button";
import { logoutAction } from "@/app/admin/logout-action";

export function LogoutButton() {
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const handleLogout = () => {
    startTransition(async () => {
      await logoutAction();
      router.replace("/login");
      router.refresh();
    });
  };

  return (
    <Button
      type="button"
      variant="ghost"
      size="sm"
      onClick={handleLogout}
      disabled={isPending}
      aria-label="Keluar"
      className="text-muted-foreground hover:text-destructive"
    >
      <LogOut className="h-4 w-4" />
      <span className="sr-only">Keluar</span>
    </Button>
  );
}
