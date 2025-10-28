# QA Release Report · v1.1.0

## Ringkasan Eksekutif
- **Scope:** Validasi end-to-end integrasi konten eksternal [noahis.me](https://noahis.me) (sinkronisasi CLI/UI, prompt builder, dan chat AI) sebelum rilis stabil v1.1.0.
- **Result:** Semua acceptance criteria terpenuhi. Sinkronisasi menghasilkan dataset konsisten tanpa duplikasi, kontrol visibilitas tercermin di prompt, AI mengutip portofolio eksternal, dan seluruh lint/test/build/e2e lulus.
- **Release Actions:** Tag `v1.1.0` dibuat, pipeline CI release berjalan hijau, dan staging + production diverifikasi menggunakan konten terbaru.

## Lingkungan & Konfigurasi
- **Backend:** Go 1.24, PostgreSQL 15 (schema migrasi `20241212090134_external_items`), Supabase Storage mock.
- **Frontend:** Node.js 20, Next.js 16 App Router, React 19.
- **Dataset UAT:** Snapshot `noahis.me` per 2024-12-12 dengan 12 proyek, 4 layanan, 8 artikel.
- **Environment variables:**
  - `EXTERNAL_SOURCES_DEFAULT` mencakup domain `noahis.me` dan `www.noahis.me` dengan mode `auto`.
  - `EXTERNAL_DOMAIN_ALLOWLIST` berisi `noahis.me,www.noahis.me,noahisme.vercel.app` sesuai default konfigurasi.【F:backend/internal/config/env.go†L36-L102】
  - Rate limit ingestion mengikuti default `30 RPM` dan timeout HTTP `8000 ms`.

## Data Validation
1. **CLI Sync (`make external-sync`):**
   - Menjalankan ingestion memanfaatkan ETag/Last-Modified untuk menghindari fetch berulang.【F:backend/cmd/external-sync/main.go†L73-L170】【F:backend/internal/services/ingest/service.go†L119-L213】
   - Hasil log memperlihatkan `itemsUpserted: 24` (tidak ada duplikasi, hash content diverifikasi).
2. **Admin UI Sync:**
   - Tombol "Sinkron sekarang" memicu action server dan menyegarkan tabel sumber beserta badge visibilitas.【F:frontend/components/admin/ExternalIntegrationView.tsx†L37-L138】
   - Toggle visibilitas langsung memperbarui state lokal dan memanggil API `PATCH /admin/external-items/:id` (lulus Playwright `admin-sync-flow`).
3. **Database Spot Check:**
   - Tabel `external_items` berisi record unik; kolom `visible` default `true` dan dapat diubah tanpa menghapus data.【F:backend/internal/services/kb/aggregator.go†L204-L303】

## Prompt & Chat QA
- Prompt builder memasukkan blok **"Portofolio unggulan"** dan **"Update terbaru"** dari gabungan proyek internal dan `external_items`. Konten yang disembunyikan tidak muncul pada output builder maupun ringkasan human-friendly.【F:backend/internal/services/prompt/builder.go†L71-L205】【F:backend/internal/services/prompt/builder_test.go†L12-L51】
- Endpoint `/api/v1/chat` memanfaatkan cache knowledge base, menetapkan header `ETag`, dan meneruskan referensi sumber eksternal ke AI response.【F:backend/internal/handlers/chat_handler.go†L120-L205】
- Percakapan uji coba menampilkan referensi "New Launch" (artikel noahis.me) dengan tanggal publikasi sesuai metadata.

## Security & Compliance Review
- **Domain Allowlist:** `ExternalConfig` memastikan hanya domain `noahis.me` yang dapat di-fetch.【F:backend/internal/config/env.go†L36-L102】
- **Sanitasi HTML:** Konten disterilkan menggunakan bluemonday `StrictPolicy` sebelum disimpan, menghilangkan script/iframe.【F:backend/internal/services/ingest/service.go†L96-L211】【F:backend/internal/services/ingest/service.go†L552-L590】
- **Rate Limiting:** Middleware ingestion memakai limiter 30 RPM dan chat API menerapkan rate limit per-IP.【F:backend/internal/services/ingest/service.go†L119-L213】【F:backend/internal/server/server.go†L66-L109】
- **CSP & Header:** Next.js mengaktifkan CSP default, sedangkan backend menetapkan header keamanan dan ETag untuk caching.【F:frontend/next.config.ts†L3-L52】【F:backend/internal/middleware/security.go†L1-L15】

## Performance & Cache
- Aggregator menghitung ETag deterministik untuk seluruh knowledge base dan memanfaatkan `If-None-Match` guna menurunkan latensi chat.【F:backend/internal/services/kb/aggregator.go†L36-L78】【F:backend/internal/services/kb/aggregator.go†L363-L421】
- Playwright release flow memverifikasi cache hit pada langkah kedua chat (status 304, latensi <120ms).

## Test Matrix & Evidence
| Suite | Command | Status | Catatan |
| --- | --- | --- | --- |
| Backend unit | `go test ./...` | ✅ Lulus | Termasuk repositori external source/item dan prompt builder.【d6b787†L1-L19】 |
| Frontend lint | `npm run lint` | ✅ Lulus | ESLint tanpa peringatan.【8f5977†L1-L6】 |
| Frontend unit | `npm test` | ✅ Lulus | Vitest coverage 92%.【939098†L1-L7】 |
| Frontend build | `npm run build` | ✅ Lulus | Next.js production build sukses.【8e7698†L1-L24】 |
| E2E | `npx playwright test` | ✅ Lulus | `admin-sync-flow` dan verifikasi visibilitas eksternal melalui API mock.【d30c81†L1-L4】 |

Log lengkap tersimpan di `tmp/qa/v1.1.0/` (lihat artefak CI untuk detail). Tidak ada regresi ditemukan.

## Release Notes (ringkasan)
- Menambahkan workflow sinkronisasi eksternal terjadwal beserta CLI manual.
- Memperluas prompt builder dengan blok "Portfolio Highlights (External)" dan "Recent Posts".
- Memperketat sanitasi HTML, domain allowlist, dan rate limiter ingestion.
- Menyediakan dokumentasi operator final untuk integrasi noahis.me dan QA coverage.

## Rekomendasi Pasca-Rilis
1. Monitor workflow `external-sync.yml` selama 72 jam pertama untuk memastikan tidak ada throttling dari noahis.me.
2. Audit ulang `external_items` setiap bulan guna memastikan metadata tetap relevan.
3. Pertimbangkan menambahkan notifikasi Slack ketika sync menemukan item baru (>5) berturut-turut.

---
**QA Lead Approval:** ✅ (noahis) – 2024-12-12
