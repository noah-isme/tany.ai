# tany.ai · v1.3.0 Stable Release

[![CI](https://github.com/tanydotai/tanyai/actions/workflows/ci.yml/badge.svg)](https://github.com/tanydotai/tanyai/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-automated-green.svg)](./docs/ARCHITECTURE.md#ci--cicd-pipeline)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](#license)

**tany.ai** adalah platform chat assistant yang menggunakan data pribadi Anda (profil, layanan, proyek) sebagai knowledge base untuk menjawab calon klien secara kontekstual. Rangkaian PR-1 s.d. PR-12 melengkapi analitik, integrasi eksternal, hingga personalisasi AI; dokumen ini menandai rilis stabil *v1.3.0*.

---

## 📦 Struktur Monorepo

```
/apps/api   → modul Go (folder `backend/`) berisi server Gin, repositori, dan middleware.
/apps/web   → aplikasi Next.js (folder `frontend/`) untuk chat publik & panel admin.
/internal/* → paket Go bersama: config, auth, services (knowledge base, prompt, storage).
/docs       → dokumentasi arsitektur, panduan admin, dan panduan deployment.
```

---

## 🚀 Fitur Utama

- **Admin CRUD lengkap** – Panel Next.js untuk mengelola profil, skill, layanan, dan proyek dengan aksi drag & drop, toggle, reorder, dan fitur featured yang langsung menginvalasi cache knowledge base melalui API Gin.【F:backend/internal/server/server.go†L30-L109】【F:frontend/components/admin/ServicesManager.tsx†L1-L207】
- **Chat berbasis knowledge base dinamis** – Endpoint `/api/v1/chat` merakit prompt grounded dari database, menyimpan riwayat percakapan, dan antarmuka chat web langsung menampilkan snippet layanan terbaru.【F:backend/internal/handlers/chat_handler.go†L22-L205】【F:backend/internal/services/prompt/builder.go†L1-L205】【F:frontend/components/chat/ChatWindow.tsx†L1-L196】
- **Integrasi portofolio eksternal noahis.me** – CLI dan UI admin dapat menyinkronkan proyek, layanan, dan artikel terbaru dari noahis.me secara otomatis maupun manual, lengkap dengan toggle visibilitas real-time dan sanitasi konten.【F:backend/cmd/external-sync/main.go†L73-L170】【F:backend/internal/services/ingest/service.go†L96-L211】【F:frontend/components/admin/ExternalIntegrationView.tsx†L1-L200】
- **Prompt builder dengan konteks eksternal** – Builder menyusun blok "Portofolio unggulan" dan "Update terbaru" dari gabungan data internal + eksternal sehingga AI merespons dengan referensi yang relevan.【F:backend/internal/services/prompt/builder.go†L71-L205】【F:backend/internal/services/kb/aggregator.go†L204-L303】
- **AI personalization & semantic memory** – pgvector menyimpan embedding profil, layanan, proyek, dan tulisan eksternal; chat handler otomatis menambahkan instruksi persona sehingga gaya jawab konsisten dengan tone penulis.【F:backend/migrations/202503150000_add_embeddings.up.sql†L1-L44】【F:backend/internal/embedding/service.go†L1-L280】【F:backend/internal/handlers/chat_handler.go†L27-L206】【F:backend/internal/services/prompt/personalization.go†L1-L49】
- **Upload & Storage terproteksi** – Admin dapat mengunggah avatar/gambar proyek ke Supabase/S3 dengan validasi MIME, sanitasi SVG, rate limit, dan logging rinci.【F:backend/internal/handlers/admin/uploads.go†L21-L214】【F:backend/internal/storage/factory.go†L9-L63】
- **Autentikasi & otorisasi** – JWT access/refresh token dengan cookie aman, middleware role-guard admin, rate limiter login/upload/chat, serta middleware Next.js untuk proteksi route.【F:backend/internal/auth/jwt.go†L16-L137】【F:frontend/middleware.ts†L1-L60】【F:backend/internal/server/server.go†L58-L108】
- **Keamanan & observabilitas produksi** – Header keamanan default, Content Security Policy, redirect HTTPS di middleware Next.js, structured JSON logging, dan panic recovery pada API.【F:backend/internal/middleware/security.go†L1-L15】【F:frontend/next.config.ts†L3-L52】【F:backend/internal/middleware/logger.go†L1-L23】【F:backend/internal/middleware/recover.go†L1-L19】
- **CI/CD otomatis** – Workflow GitHub Actions menjalankan lint, unit test, build, dan Playwright e2e untuk backend dan frontend di setiap push/PR, plus workflow terjadwal untuk sinkronisasi sumber eksternal.【F:.github/workflows/ci.yml†L1-L88】【F:.github/workflows/external-sync.yml†L1-L83】

---

## 🛠️ Quickstart Developer

1. Pastikan dependensi Go (≥1.24) dan Node.js (≥20) terpasang, kemudian install paket frontend: `npm install --prefix frontend`.
2. Salin variabel contoh lalu sesuaikan kredensial Postgres & storage: `cp backend/.env.example backend/.env`.
3. Jalankan migrasi & seed: `make migrate` lalu `make seed`.
4. Mulai seluruh stack lokal:
   ```bash
   make dev
   # buka http://localhost:3000 untuk frontend & http://localhost:8080 untuk API
   ```
5. Gunakan kredensial seeding default `admin@example.com / Admin#12345` untuk masuk ke panel admin.

## 🚢 Deployment ke Server

### Persiapan

1. Siapkan server dengan:
   - PostgreSQL (≥15)
   - Node.js (≥20)
   - Go (≥1.24)
   - Nginx sebagai reverse proxy

2. Setup environment variables:
   ```bash
   # Clone repository
   git clone https://github.com/tanydotai/tanyai.git
   cd tanyai

   # Setup backend environment
   cp backend/.env.example backend/.env
   # Edit backend/.env sesuai konfigurasi production
   ```

3. Konfigurasi AI Provider:
   ```bash
   # Untuk Google Gemini
   AI_PROVIDER=gemini
   GEMINI_MODEL=gemini-2.5-flash  # atau gemini-2.5-pro
   GOOGLE_GENAI_API_KEY=your-key-here

   # Untuk Leapcell
   AI_PROVIDER=leapcell
   LEAPCELL_API_KEY=your-key-here
   LEAPCELL_PROJECT_ID=your-project-id
   LEAPCELL_TABLE_ID=your-table-id
   ```

### Build & Run

1. Frontend:
   ```bash
   cd frontend
   npm install
   npm run build
   npm run start      # Atau gunakan PM2: pm2 start npm --name "tanyai-web" -- start
   ```

2. Backend:
   ```bash
   cd backend
   make migrate       # Jalankan migrasi database
   make seed         # Opsional: seed data awal
   make build        # Build binary
   ./tmp/main       # Atau gunakan systemd untuk menjalankan service
   ```

3. Setup Nginx:
   ```nginx
   # /etc/nginx/sites-available/tanyai
   server {
      listen 80;
      server_name your-domain.com;

      # Frontend
      location / {
          proxy_pass http://localhost:3000;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
      }

      # Backend API
      location /api {
          proxy_pass http://localhost:8080;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
      }
   }
   ```

4. Enable HTTPS dengan Certbot:
   ```bash
   certbot --nginx -d your-domain.com
   ```

### Monitoring & Maintenance

- Setup logging dengan systemd untuk backend
- Gunakan PM2 untuk monitoring frontend
- Backup database secara berkala
- Monitor rate limits dan penggunaan storage

Panduan lebih rinci tersedia di dokumen berikut:

- [Arsitektur end-to-end](./docs/ARCHITECTURE.md)
- [Panduan operasional Admin](./docs/ADMIN_GUIDE.md)
- [Panduan AI personalization](./docs/AI_PERSONALIZATION_GUIDE.md)
- [Panduan deployment & konfigurasi environment](./docs/DEPLOYMENT.md)
- [Checklist Release v1.1.0](./RELEASE_CHECKLIST.md)
- [Laporan QA rilis v1.1.0](./docs/QA_RELEASE_V1.1.0.md)

---

## 🧱 Stack Teknologi

- **Backend**: Go + Gin, sqlx, PostgreSQL, Supabase/S3 object storage, JWT auth, rate limiter.
- **Frontend**: Next.js 16 (App Router), React 19, Tailwind utility, React Hook Form + Zod, Playwright untuk e2e.
- **Testing & QA**: Go test, Vitest, React Testing Library, Playwright smoke & release flow.
- **CI/CD**: GitHub Actions, artefak coverage backend, otomatisasi e2e.

---

## 📄 License

MIT License © tany.ai. Silakan gunakan atau fork untuk kebutuhan Anda sendiri.

