# Panduan Deployment tany.ai

Panduan ini merangkum konfigurasi environment, proses build, dan pipeline deployment untuk staging & production v1.0.0.

## 1. Konfigurasi Environment
Siapkan variabel berikut sebelum build:

| Kategori | Variabel | Keterangan |
| --- | --- | --- |
| Aplikasi | `APP_ENV` | `staging` atau `production` untuk menandai lingkungan.【F:backend/internal/config/env.go†L37-L72】 |
| Database | `POSTGRES_URL` | URL koneksi Postgres (wajib).【F:backend/internal/config/env.go†L65-L72】 |
| Database (opsional) | `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME` | Tuning koneksi sesuai kebutuhan produksi.【F:backend/internal/config/env.go†L92-L123】 |
| Autentikasi | `JWT_SECRET` | Panjang minimal 32 karakter (wajib).【F:backend/internal/config/env.go†L73-L85】 |
| Token TTL | `ACCESS_TOKEN_TTL_MIN`, `REFRESH_TOKEN_TTL_DAY`, `REFRESH_COOKIE_NAME` | Menyesuaikan masa berlaku token & nama cookie refresh.【F:backend/internal/config/env.go†L107-L165】 |
| Rate Limit | `LOGIN_RATE_LIMIT_PER_MIN`, `LOGIN_RATE_LIMIT_BURST`, `KB_RATE_LIMIT_PER_5MIN`, `CHAT_RATE_LIMIT_PER_5MIN`, `UPLOAD_RATE_LIMIT_PER_MIN`, `UPLOAD_RATE_LIMIT_BURST` | Sesuaikan kapasitas trafik untuk login, knowledge base, chat, dan upload.【F:backend/internal/config/env.go†L125-L214】 |
| Knowledge Base | `KB_CACHE_TTL_SECONDS` | TTL cache knowledge base (detik).【F:backend/internal/config/env.go†L189-L206】 |
| Storage | `STORAGE_DRIVER` (`supabase`/`s3`), `SUPABASE_URL`, `SUPABASE_BUCKET`, `SUPABASE_SERVICE_ROLE`, `SUPABASE_PUBLIC_URL` | Wajib bila memakai Supabase Storage.【F:backend/internal/config/env.go†L215-L254】 |
| Storage (S3) | `S3_REGION`, `S3_BUCKET`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `S3_ENDPOINT`, `S3_PUBLIC_BASE_URL`, `S3_FORCE_PATH_STYLE` | Wajib bila memakai penyimpanan kompatibel S3.【F:backend/internal/config/env.go†L254-L314】 |
| Upload Policy | `UPLOAD_MAX_MB`, `UPLOAD_ALLOWED_MIME`, `ALLOW_SVG` | Atur kebijakan upload sesuai kebutuhan brand.【F:backend/internal/config/env.go†L147-L206】 |
| AI Model | `AI_MODEL` | Nama model yang dicatat di log & respons chat.【F:backend/internal/config/env.go†L206-L214】 |
| Frontend → API | `API_BASE_URL` atau `NEXT_PUBLIC_API_BASE_URL` | URL API yang digunakan Next.js (produk & staging).【F:frontend/lib/env.ts†L1-L9】 |
| Frontend Origin | `FRONTEND_ORIGIN` | (Opsional) domain publik frontend untuk konfigurasi reverse proxy & CSP tambahan. |

**Catatan:** simpan kredensial rahasia (JWT, Supabase, AWS) di secret manager platform Anda (GitHub Actions, Vercel, Railway, Fly.io, dsb.).

## 2. Database Migration & Seed
1. Pastikan database siap pakai.
2. Jalankan migrasi: `make migrate` (meneruskan ke `backend/cmd/migrate`).【F:backend/Makefile†L7-L9】
3. Jalankan seed dasar: `make seed` untuk mengisi profil, layanan, dan proyek contoh.【F:backend/Makefile†L11-L13】【F:backend/cmd/seed/main.go†L22-L115】

## 3. Build & Jalankan Backend
1. Install dependency Go: `go mod download` di folder `backend`.
2. Build binary: `go build -o bin/api ./cmd/api` (gunakan flag tambahan bila perlu).【F:backend/cmd/api/main.go†L16-L55】
3. Jalankan dengan konfigurasi production:
   ```bash
   ./bin/api --env=production
   ```
   Proses ini memuat config environment, membuka koneksi database, dan menjalankan server Gin pada port `PORT` (default 8080).【F:backend/cmd/api/main.go†L16-L55】【F:backend/internal/server/server.go†L104-L134】
4. Pastikan log JSON menampilkan `method`, `path`, `status`, `latency_ms`, `ip` untuk memudahkan monitoring.【F:backend/internal/middleware/logger.go†L1-L23】

## 4. Build & Jalankan Frontend
1. Masuk ke `frontend/`, install dependency: `npm install` atau `npm ci` di CI.
2. Build produksi: `npm run build` (Next.js).【F:frontend/package.json†L6-L18】
3. Jalankan server produksi: `npm run start` dengan variabel `PORT` sesuai kebutuhan (default 3000). Middleware akan memaksa HTTPS dan memvalidasi token admin di produksi.【F:frontend/middleware.ts†L1-L60】
4. Pastikan header keamanan (CSP, HSTS, Referrer-Policy) aktif di semua response.【F:frontend/next.config.ts†L15-L52】

## 5. Integrasi CI/CD
- Workflow GitHub Actions `CI` otomatis menjalankan lint, test, build, dan Playwright e2e. Gunakan badge CI di README untuk memantau status setiap branch.【F:.github/workflows/ci.yml†L9-L88】
- Tambahkan deploy job terpisah yang terpicu setelah pipeline `CI` sukses untuk mendorong artefak ke Vercel (frontend) dan Railway/Fly.io (backend). Pastikan rahasia deployment tersimpan di repository secrets.

## 6. Monitoring & Logging
- Gunakan log JSON backend untuk menganalisis latensi & status request. Tambahkan shipping log (misal ke Grafana Loki atau ELK) sesuai infrastruktur Anda.【F:backend/internal/middleware/logger.go†L1-L23】
- Gunakan header `request_id` (bisa ditambahkan di reverse proxy) untuk korelasi error lintas layanan.
- Pantau storage upload melalui log handler uploads (menyertakan mime, size, key, latency).【F:backend/internal/handlers/admin/uploads.go†L130-L206】

## 7. Security Checklist
| ✅ | Item | Sumber |
| --- | --- | --- |
| ✅ | Content Security Policy & header keamanan aktif | Next.js headers & middleware security pada backend.【F:frontend/next.config.ts†L15-L52】【F:backend/internal/middleware/security.go†L1-L15】 |
| ✅ | HTTPS redirect di frontend production | Middleware Next.js memaksa protokol `https`.【F:frontend/middleware.ts†L7-L34】 |
| ✅ | JWT tervalidasi & role guard admin | Middleware backend `Authn` + `AuthzAdmin` serta middleware Next.js.|【F:backend/internal/server/server.go†L58-L116】【F:frontend/middleware.ts†L16-L60】 |
| ✅ | Rate limit global untuk login/chat/knowledge/upload | Rate limiter dibangun via `auth.NewRateLimiter`.【F:backend/internal/auth/ratelimit.go†L12-L68】【F:backend/internal/server/server.go†L58-L116】 |
| ✅ | Panic recovery & logging untuk API | Middleware `RecoverWithLog` menangani panic dan mencatat stack trace.【F:backend/internal/middleware/recover.go†L1-L19】 |

## 8. Deploy Checklist Singkat
1. Pastikan pipeline `CI` hijau.
2. Verifikasi `.env.production` berisi seluruh variabel di atas.
3. Jalankan `go vet ./... && go test ./... && go build ./...` di backend sebelum tagging.
4. Jalankan `npm run lint && npm test && CI=1 npm run build && npm run test:e2e` di frontend.
5. Buat tag `v1.0.0` dan release notes setelah staging disetujui.

