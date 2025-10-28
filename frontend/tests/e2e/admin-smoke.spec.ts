import { test, expect } from "@playwright/test";

import {
  authenticateAdmin,
  createProject,
  createSkill,
  deleteProject,
  deleteSkill,
  fetchServices,
  fetchSkills,
  toggleService,
  reorderSkills,
} from "./utils/auth";

test("admin skill management flow", async ({ page }) => {
  const loginResponse = await page.goto("/login");
  if (loginResponse?.status() === 404) {
    await page.waitForTimeout(2000);
    await page.reload();
  }
  await page.waitForSelector('input[name="email"]');
  await expect(page.getByRole("button", { name: "Masuk" })).toBeVisible();

  const token = await authenticateAdmin(page, "admin@example.com", "Password123!");

  await page.goto("/admin/skills");
  const main = page.locator("#admin-main");
  await expect(main.getByRole("heading", { name: "Skills" })).toBeVisible();

  const skillName = `Testing Skill ${Date.now()}`;
  const createdSkill = await createSkill(page, token, skillName);
  const skills = await fetchSkills(page, token);

  const reorderedSkills = [
    { id: createdSkill.id, order: 0 },
    ...skills
      .filter((skill) => skill.id !== createdSkill.id)
      .map((skill, index) => ({ id: skill.id, order: index + 1 })),
  ];
  await reorderSkills(page, token, reorderedSkills);

  await page.reload();
  await expect(page.locator("tbody tr").first()).toContainText(skillName);

  await deleteSkill(page, token, createdSkill.id);
  await page.reload();
  await expect(page.getByText(skillName)).not.toBeVisible();

  await page.goto("/admin/services");
  const services = await fetchServices(page, token);
  const targetService = services.find((service) => service.name === "AI Discovery Workshop");
  if (!targetService) {
    throw new Error("AI Discovery Workshop service not found");
  }

  await toggleService(page, token, targetService.id, false);
  await page.reload();
  const toggleOff = page
    .locator("tr", { hasText: "AI Discovery Workshop" })
    .first()
    .getByRole("checkbox", { name: /toggle status layanan ai discovery workshop/i });
  await expect(toggleOff).not.toBeChecked({ timeout: 15000 });

  await page.goto("/");
  const snippet = page.getByText(/Saya bisa membantu/i).first();
  await expect(snippet).toBeVisible();
  await expect(snippet).not.toContainText("AI Discovery Workshop");

  await toggleService(page, token, targetService.id, true);
  await page.goto("/admin/services");
  const toggleOn = page
    .locator("tr", { hasText: "AI Discovery Workshop" })
    .first()
    .getByRole("checkbox", { name: /toggle status layanan ai discovery workshop/i });
  await expect(toggleOn).toBeChecked({ timeout: 15000 });

  await page.goto("/admin/projects");
  const projectTitle = `Project Upload Test ${Date.now()}`;
  const project = await createProject(page, token, {
    title: projectTitle,
    description: "Testing upload flow",
    tech_stack: ["Playwright", "Next.js"],
    image_url: `https://mock-storage.example.com/uploads/${Date.now()}.png`,
    project_url: "https://example.com/project",
    category: "Testing",
    duration_label: "4 minggu",
    price_label: "QA sprint",
    budget_label: "IDR 60Jt",
    is_featured: true,
  });

  await page.reload();
  const newProjectRow = page.locator("tr", { hasText: projectTitle }).first();
  await expect(newProjectRow).toBeVisible();

  await deleteProject(page, token, project.id);
  await page.reload();
  await expect(page.getByText(projectTitle)).not.toBeVisible();

  await page.goto("/admin/integrations");
  await expect(page.getByRole("heading", { name: "Integrasi Konten Eksternal" })).toBeVisible();
  const syncButton = page.getByRole("button", { name: /Sinkron sekarang/i }).first();
  await syncButton.click();
  await page.waitForLoadState("networkidle");
  await expect(page.getByRole("status")).toContainText(/Sinkronisasi mock berhasil/i, {
    timeout: 15000,
  });

  const externalRows = page
    .locator("section", { hasText: "Konten yang Tersedia" })
    .locator("tbody tr");

  await expect
    .poll(async () => externalRows.count(), { timeout: 15000 })
    .toBeGreaterThan(0);

  const firstRow = externalRows.first();
  await expect(firstRow).toBeVisible({ timeout: 15000 });

  const toggleLabel = firstRow.locator(
    'label:has(input[aria-label^="Atur visibilitas"])',
  );
  await expect(toggleLabel).toBeVisible({ timeout: 15000 });

  const toggleInput = toggleLabel.locator('input[type="checkbox"]');
  const initialState = await toggleInput.isChecked();

  await toggleLabel.click();
  await expect(toggleInput).toHaveJSProperty("checked", !initialState);

  await toggleLabel.click();
  await expect(toggleInput).toHaveJSProperty("checked", initialState);

  await page.getByLabel("Keluar").click();
});
