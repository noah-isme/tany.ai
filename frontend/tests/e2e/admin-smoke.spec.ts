import { test, expect } from "@playwright/test";

test("admin skill management flow", async ({ page }) => {
  const loginResponse = await page.goto("/login");
  if (loginResponse?.status() === 404) {
    await page.waitForTimeout(2000);
    await page.reload();
  }
  await page.waitForSelector('input[name="email"]');

  await page.fill('input[name="email"]', "admin@example.com");
  await page.fill('input[name="password"]', "Password123!");
  await page.click('button:has-text("Masuk")');

  await page.waitForURL("**/admin");

  await page.goto("/admin/skills");
  const main = page.locator("#admin-main");
  await expect(main.getByRole("heading", { name: "Skills" })).toBeVisible();

  await page.fill('input[placeholder="Contoh: Prompt Engineering"]', "Testing Skill");
  await page.click('button:has-text("Tambah")');
  await expect(page.getByText("Testing Skill")).toBeVisible();

  const newSkillRow = page.locator("tr", { hasText: "Testing Skill" });
  const moveUpButton = newSkillRow.getByRole("button", { name: "Pindah ke atas" });
  await moveUpButton.click();
  await expect(page.locator("tbody tr").nth(1)).toContainText("Testing Skill");
  await newSkillRow.getByRole("button", { name: "Pindah ke atas" }).click();

  await expect(page.locator("tbody tr").first()).toContainText("Testing Skill");

  page.on("dialog", (dialog) => dialog.accept());
  await newSkillRow.getByRole("button", { name: "Hapus" }).click();
  await expect(page.getByText("Testing Skill")).not.toBeVisible();

  await page.goto("/admin/services");
  const serviceRow = page.locator("tr", { hasText: "AI Discovery Workshop" }).first();
  const toggle = serviceRow.getByRole("checkbox", { name: /toggle status layanan ai discovery workshop/i });
  await expect(toggle).toBeChecked();
  await toggle.click({ force: true });
  await expect(toggle).not.toBeChecked();
  await toggle.click({ force: true });
  await expect(toggle).toBeChecked();

  await page.goto("/admin/projects");
  const projectRow = page.locator("tr", { hasText: "Sales Enablement Chatbot" }).first();
  const featureButton = projectRow.getByRole("button", { name: "Jadikan featured" });
  await featureButton.click();
  await expect(projectRow.getByRole("button", { name: "Hapus featured" })).toBeVisible();

  await page.getByLabel("Keluar").click();
  await page.waitForURL("**/login");
});
