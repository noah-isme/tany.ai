import { defineConfig, devices } from "@playwright/test";

const PORT = Number(process.env.PORT ?? 3000);
const API_PORT = Number(process.env.MOCK_API_PORT ?? 4000);
const JWT_SECRET = process.env.JWT_SECRET ?? "test-secret";

export default defineConfig({
  testDir: "./tests/e2e",
  timeout: 60_000,
  retries: process.env.CI ? 1 : 0,
  use: {
    baseURL: `http://localhost:${PORT}`,
    trace: "on-first-retry",
  },
  webServer: [
    {
      command: `npm run dev`,
      port: PORT,
      reuseExistingServer: !process.env.CI,
      env: {
        ...process.env,
        MOCK_API_PORT: String(API_PORT),
        API_BASE_URL: `http://127.0.0.1:${API_PORT}`,
        JWT_SECRET,
      },
    },
    {
      command: `node ./tests/e2e/mock-backend.mjs`,
      port: API_PORT,
      reuseExistingServer: !process.env.CI,
      env: {
        ...process.env,
        MOCK_API_PORT: String(API_PORT),
        JWT_SECRET,
      },
    },
  ],
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],
});
