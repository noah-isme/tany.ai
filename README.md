# tany.ai · v1.0.0 Release Candidate

[![CI](https://github.com/tanydotai/tanyai/actions/workflows/ci.yml/badge.svg)](https://github.com/tanydotai/tanyai/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-automated-green.svg)](./docs/ARCHITECTURE.md#ci--cicd-pipeline)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](#license)

**tany.ai** adalah platform chat assistant yang menggunakan data pribadi Anda (profil, layanan, proyek) sebagai knowledge base untuk menjawab calon klien secara kontekstual. Rangkaian PR-1 s.d. PR-7 telah menyelesaikan seluruh fungsionalitas inti dan dokumentasi ini menandai fase *v1.0.0-ready*.

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
- **Chat berbasis knowledge base dinamis** – Endpoint `/api/v1/chat` merakit prompt deterministik dari database, menyimpan riwayat percakapan, dan antarmuka chat web langsung menampilkan snippet layanan terbaru.【F:backend/internal/handlers/chat_handler.go†L22-L118】【F:backend/internal/services/prompt/builder.go†L9-L99】【F:frontend/components/chat/ChatWindow.tsx†L1-L96】
- **Upload & Storage terproteksi** – Admin dapat mengunggah avatar/gambar proyek ke Supabase/S3 dengan validasi MIME, sanitasi SVG, rate limit, dan logging rinci.【F:backend/internal/handlers/admin/uploads.go†L21-L214】【F:backend/internal/storage/factory.go†L9-L63】
- **Autentikasi & otorisasi** – JWT access/refresh token dengan cookie aman, middleware role-guard admin, rate limiter login/upload/chat, serta middleware Next.js untuk proteksi route.【F:backend/internal/auth/jwt.go†L16-L137】【F:frontend/middleware.ts†L1-L60】【F:backend/internal/server/server.go†L58-L108】
- **Keamanan & observabilitas produksi** – Header keamanan default, Content Security Policy, redirect HTTPS di middleware Next.js, structured JSON logging, dan panic recovery pada API.【F:backend/internal/middleware/security.go†L1-L15】【F:frontend/next.config.ts†L3-L52】【F:backend/internal/middleware/logger.go†L1-L23】【F:backend/internal/middleware/recover.go†L1-L19】
- **CI/CD otomatis** – Workflow GitHub Actions menjalankan lint, unit test, build, dan Playwright e2e untuk backend dan frontend di setiap push/PR.【F:.github/workflows/ci.yml†L1-L88】

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
5. Gunakan kredensial seeding default `admin@example.com / Password123!` untuk masuk ke panel admin.

Panduan lebih rinci tersedia di dokumen berikut:

- [Arsitektur end-to-end](./docs/ARCHITECTURE.md)
- [Panduan operasional Admin](./docs/ADMIN_GUIDE.md)
- [Panduan deployment & konfigurasi environment](./docs/DEPLOYMENT.md)
- [Checklist Release v1.0.0](./RELEASE_CHECKLIST.md)

---

## 🧱 Stack Teknologi

- **Backend**: Go + Gin, sqlx, PostgreSQL, Supabase/S3 object storage, JWT auth, rate limiter.
- **Frontend**: Next.js 16 (App Router), React 19, Tailwind utility, React Hook Form + Zod, Playwright untuk e2e.
- **Testing & QA**: Go test, Vitest, React Testing Library, Playwright smoke & release flow.
- **CI/CD**: GitHub Actions, artefak coverage backend, otomatisasi e2e.

---

## 📄 License

MIT License © tany.ai. Silakan gunakan atau fork untuk kebutuhan Anda sendiri.

