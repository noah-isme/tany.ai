# Changelog

# [v1.3.0] - 2025-03-15
### Added
- Skema pgvector baru `embeddings` dan `embedding_config` untuk menyimpan memori semantik, lengkap dengan indeks IVFFLAT agar pencarian cosine similarity efisien.【F:backend/migrations/202503150000_add_embeddings.up.sql†L1-L44】
- Paket `internal/embedding` yang menangani repository, layanan caching, handler admin, serta integrasi dengan analytics metadata.【F:backend/internal/embedding/service.go†L1-L280】【F:backend/internal/embedding/handler.go†L1-L92】【F:backend/internal/handlers/chat_handler.go†L27-L206】
- Provider OpenAI embedding dan builder prompt baru yang menyisipkan instruksi persona (`BuildPersonalizedPrompt`).【F:backend/internal/ai/openai_embeddings.go†L1-L97】【F:backend/internal/services/prompt/personalization.go†L1-L49】
- Panel admin **AI Personalization** lengkap dengan slider bobot, statistik embedding, serta aksi reindex/reset + server actions Next.js.【F:frontend/app/admin/personalization/page.tsx†L1-L28】【F:frontend/components/admin/PersonalizationPanel.tsx†L1-L196】【F:frontend/lib/admin-api.ts†L1-L214】
- Dokumentasi operator `docs/AI_PERSONALIZATION_GUIDE.md` untuk arsitektur, variabel environment, dan QA checklist personalisasi.【F:docs/AI_PERSONALIZATION_GUIDE.md†L1-L92】

### Changed
- Konfigurasi runtime menambahkan `ENABLE_PERSONALIZATION`, `EMBEDDING_*`, dan `PERSONALIZATION_WEIGHT`, serta otomatis memuat provider embedding saat server diinisialisasi.【F:backend/internal/config/env.go†L24-L233】【F:backend/internal/server/server.go†L30-L169】
- Chat handler kini memanggil personalizer sebelum mem-build prompt dan mengirim metrik tambahan ke analytics.【F:backend/internal/handlers/chat_handler.go†L27-L206】
- README diperbarui ke status rilis v1.3.0 dengan highlight personalisasi AI dan tautan panduan baru.【F:README.md†L1-L124】

## [v1.2.0] - 2025-01-15
### Added
- Modul analytics real-time dengan tabel baru `analytics_events` dan `analytics_summary` untuk menyimpan event chat dan ringkasan harian.【F:backend/migrations/202501010000_add_analytics_tables.up.sql†L1-L40】
- Service & repository analytics di backend termasuk endpoint admin `/api/admin/analytics/{summary,events,leads}` dan integrasi otomatis dari `ChatHandler`.【F:backend/internal/analytics/service.go†L1-L188】【F:backend/internal/server/server.go†L36-L126】
- Dashboard Admin → Analytics menggunakan Next.js 16 + Recharts lengkap dengan filter rentang waktu, auto-refresh, dan widget metrik utama.【F:frontend/components/admin/analytics/AnalyticsDashboard.tsx†L1-L399】【F:frontend/app/admin/analytics/page.tsx†L1-L24】
- Dokumentasi operator baru `docs/ANALYTICS_GUIDE.md` mencakup arsitektur, konfigurasi environment, dan panduan pengujian modul analytics.【F:docs/ANALYTICS_GUIDE.md†L1-L89】

### Changed
- Konfigurasi environment menambahkan dukungan `ENABLE_ANALYTICS` dan `ANALYTICS_RETENTION_DAYS` untuk mengontrol retensi serta aktivasi modul analytics.【F:backend/internal/config/env.go†L24-L205】

## [v1.1.0] - 2024-12-12
### Added
- Integrasi penuh dengan noahis.me mencakup CLI `make external-sync`, UI admin, dan workflow terjadwal untuk sinkronisasi konten eksternal.【F:backend/cmd/external-sync/main.go†L73-L170】【F:frontend/components/admin/ExternalIntegrationView.tsx†L1-L200】【F:.github/workflows/external-sync.yml†L1-L83】
- Prompt builder serta chat API kini menyisipkan blok "Portofolio unggulan" dan "Update terbaru" dengan sumber eksternal lengkap dengan header `ETag` untuk caching respons.【F:backend/internal/services/prompt/builder.go†L71-L205】【F:backend/internal/handlers/chat_handler.go†L120-L205】
- Dokumentasi operator (`docs/INTEGRATIONS_NOAHISME.md`) dan laporan QA rilis (`docs/QA_RELEASE_V1.1.0.md`) untuk panduan produksi dan bukti validasi.【F:docs/INTEGRATIONS_NOAHISME.md†L1-L129】【F:docs/QA_RELEASE_V1.1.0.md†L1-L98】

### Changed
- README diperbarui menjadi status rilis stabil v1.1.0 dengan highlight integrasi eksternal dan cakupan QA terbaru.【F:README.md†L1-L86】

### Fixed
- Menjaga idempoten sinkronisasi dengan pemanfaatan hash, ETag, dan Last-Modified sehingga tidak terjadi duplikasi konten ketika workflow terjadwal berjalan berulang.【F:backend/internal/services/ingest/service.go†L119-L213】

## [v1.0.0] - 2024-11-20
- Rilis kandidat pertama dengan CRUD admin, chat GPT-style, autentikasi JWT, dan pipeline CI lengkap.

