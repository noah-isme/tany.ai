import { test, expect } from "@playwright/test";

import {
  authenticateAdmin,
  createService,
  deleteService,
  updateProfile,
  API_BASE_URL,
} from "./utils/auth";

const ADMIN_EMAIL = "admin@example.com";
const ADMIN_PASSWORD = "Password123!";

test("release readiness end-to-end flow", async ({ page }) => {
  const loginResponse = await page.goto("/login");
  expect(loginResponse?.status()).toBe(200);
  expect(loginResponse?.headers()["content-security-policy"] ?? "").toContain("default-src 'self'");
  expect(loginResponse?.headers()["strict-transport-security"] ?? "").toContain("max-age=");

  const redirectResponse = await page.request.fetch("/admin", {
    maxRedirects: 0,
    headers: { "x-forwarded-proto": "http" },
  });
  expect([301, 302, 307, 308]).toContain(redirectResponse.status());
  expect(redirectResponse.headers()["location"] ?? "").toMatch(/^https:/i);

  const token = await authenticateAdmin(page, ADMIN_EMAIL, ADMIN_PASSWORD);
  await page.goto("/admin");
  await expect(page.getByText(/Ringkasan Profil/i)).toBeVisible();

  const newTitle = `Head of Delivery ${Date.now()}`;
  await updateProfile(page, token, { title: newTitle, location: "Bandung, Indonesia" });

  await page.goto("/admin");
  await expect(page.locator("dd", { hasText: newTitle }).first()).toBeVisible();
  await expect(page.locator("dd", { hasText: "Bandung, Indonesia" }).first()).toBeVisible();

  const serviceName = `Release QA ${Date.now()}`;
  const service = await createService(page, token, {
    name: serviceName,
    description: "Smoke coverage untuk final release.",
    duration_label: "1 minggu",
    price_min: 5000000,
    price_max: 7500000,
    currency: "IDR",
    is_active: true,
  });

  await page.goto("/admin/services");
  const serviceRow = page.locator("tr", { hasText: serviceName }).first();
  await expect(serviceRow).toBeVisible();

  const knowledgeResponse = await page.request.get(`${API_BASE_URL}/api/v1/knowledge-base`);
  expect(knowledgeResponse.ok()).toBeTruthy();
  const knowledge = (await knowledgeResponse.json()) as {
    services: { name: string }[];
  };
  expect(knowledge.services.map((service) => service.name)).toContain(serviceName);

  await page.goto("/");
  const servicesSection = page.locator("#services");
  await expect(servicesSection.getByRole("heading", { name: serviceName })).toBeVisible();

  await deleteService(page, token, service.id);
});
