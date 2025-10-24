# Release Checklist v1.0.0

Gunakan daftar berikut sebelum menandai rilis production/staging.

- [ ] go vet ./... && go test ./... && go build ./...
- [ ] npm run lint && npm test && CI=1 npm run build && npm run test:e2e
- [ ] .env.production diverifikasi & terenkripsi di secret manager
- [ ] Supabase/S3 bucket publik siap & kredensial diuji
- [ ] Badge CI di README berstatus hijau
- [ ] README.md dan seluruh dokumentasi (`/docs`) sudah final
- [ ] Tag rilis `v1.0.0` dibuat & release notes dipublikasikan

