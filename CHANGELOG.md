# Changelog

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

