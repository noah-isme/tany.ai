# Changelog

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

