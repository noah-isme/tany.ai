import { apiFetch } from "./api-client";

export type LoginResponse = {
  accessToken: string;
  user: {
    id: string;
    email: string;
    name?: string | null;
    roles: string[];
  };
};

export async function loginRequest(email: string, password: string): Promise<LoginResponse> {
  return apiFetch<LoginResponse>("/api/auth/login", {
    method: "POST",
    body: { email, password },
    withAuth: false,
  });
}

export async function logoutRequest(): Promise<void> {
  await apiFetch("/api/auth/logout", {
    method: "POST",
    withAuth: true,
  });
}
