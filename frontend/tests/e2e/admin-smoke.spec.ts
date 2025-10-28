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

  await page.route("**/api/admin/external/items/*/visibility", async (route) => {
    if (route.request().method() === "PATCH") {
      let visible = false;
      try {
        const body = route.request().postDataJSON?.();
        if (body && typeof body.visible === "boolean") {
          visible = body.visible;
        }
      } catch {
        // ignore JSON parsing errors and fall back to false
      }

      const url = new URL(route.request().url());
      const segments = url.pathname.split("/").filter(Boolean);
      const itemId = segments.at(-2) ?? "external-item";

      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          data: {
            id: itemId,
            sourceName: "noahis.me",
            kind: "project",
            title: "Mock External Item",
            summary: "Konten mock Playwright",
            url: "https://www.noahis.me/mock-item",
            visible,
            publishedAt: new Date().toISOString(),
            metadata: {},
          },
        }),
      });
      return;
    }

    await route.continue();
  });

  await page.goto("/admin/integrations");
  await expect(page.getByRole("heading", { name: "Integrasi Konten Eksternal" })).toBeVisible();
  const syncButton = page.getByRole("button", { name: /Sinkron sekarang/i }).first();
  await syncButton.click();
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

  const toggleInput = firstRow.getByRole("checkbox", { name: /^Atur visibilitas/i });
  await expect(toggleInput).toBeVisible({ timeout: 15000 });

  const initialState = await toggleInput.isChecked();

  const waitForToggleResponse = () =>
    page.waitForResponse((response) => {
      if (response.request().method() !== "PATCH") {
        return false;
      }
      const url = new URL(response.url());
      return /\/api\/admin\/external\/items\/.+\/visibility$/.test(url.pathname) && response.ok();
    });

  const ensureState = async (checked: boolean) => {
    if ((await toggleInput.isChecked()) !== checked) {
      await toggleInput.focus();
      await Promise.all([waitForToggleResponse(), page.keyboard.press(" ")]);
    }
    await expect(toggleInput)[checked ? "toBeChecked" : "not.toBeChecked"]({
      timeout: 15000,
    });
  };

  if (initialState) {
    await ensureState(false);
    await ensureState(true);
  } else {
    await ensureState(true);
    await ensureState(false);
  }

  await page.getByLabel("Keluar").click();
});
