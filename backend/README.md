# tany.ai Backend (Golang + Gin)

Layanan API awal untuk prototipe tany.ai. Server ini menyiapkan endpoint dasar dan knowledge base statis yang akan digunakan
sebagai konteks bagi AI assistant.

## 🚀 Fitur
- Endpoint `GET /healthz` untuk pemeriksaan kesehatan.
- Endpoint `GET /api/v1/knowledge-base` untuk melihat data profil, layanan, portofolio, dan paket harga.
- Endpoint `POST /api/v1/chat` yang mengembalikan mock jawaban berdasarkan knowledge base dan system prompt.

## 🧱 Struktur Direktori
```
cmd/api/main.go          # Entry point server
internal/server/         # Konfigurasi Gin dan routing
internal/handlers/       # Handler HTTP untuk chat & health
internal/knowledge/      # Data statis dan utilitas prompt
```

## 🧪 Testing
```
go test ./...
```

Unit test meliputi:
- Validasi konten system prompt.
- Pengujian handler chat untuk memastikan struktur respons konsisten.

## ▶️ Menjalankan Server
```
PORT=8080 go run ./cmd/api
```

Server akan berjalan di `http://localhost:8080`.

## 🗄️ Database Tooling

Backend kini menggunakan PostgreSQL (Supabase kompatibel) sebagai sumber knowledge base. Workflow dasar:

```bash
cp .env.example .env        # konfigurasi kredensial DB
make migrate                # jalankan migrasi SQL
make seed                   # isi contoh data profil/skills/services/projects
```

Endpoint `GET /healthz` akan mengembalikan status koneksi database melalui field `database`.

## 🛠️ Admin API (Preview)

Backend menyediakan sekumpulan endpoint `REST` di bawah `/api/admin/*` untuk mengelola knowledge base. Semua response sukses menggunakan bentuk `{ "data": ... }` (atau `{items, page, limit, total}` untuk list) dan format error tunggal `{ "error": { code, message, details } }`.

> **Catatan**
>
> * Middleware `AuthzAdminStub` dapat diaktifkan melalui `ENABLE_ADMIN_GUARD=true` (default aktif otomatis bila `APP_ENV=prod`). Saat aktif, stub akan merespons `401` atau `403` berdasarkan `ADMIN_GUARD_MODE`.
> * Operasi `DELETE` pada skills/services/projects bersifat hard delete.
> * Reorder dilakukan dalam transaksi untuk menjaga konsistensi urutan.

### Profil
- `GET /api/admin/profile`
- `PUT /api/admin/profile`

Contoh request:
```http
PUT /api/admin/profile
Content-Type: application/json

{
  "name": "Jane Doe",
  "title": "Product Designer",
  "bio": "Desainer dengan fokus pada SaaS dan AI products.",
  "email": "jane@example.com",
  "location": "Jakarta, Indonesia",
  "avatar_url": "https://cdn.example.com/avatar.png"
}
```

Contoh response:
```json
{
  "data": {
    "id": "11111111-1111-1111-1111-111111111111",
    "name": "Jane Doe",
    "title": "Product Designer",
    "email": "jane@example.com",
    "updated_at": "2024-10-10T09:30:00Z"
  }
}
```

### Skills
- `GET /api/admin/skills?page=1&limit=20&sort=order&dir=asc`
- `POST /api/admin/skills`
- `PUT /api/admin/skills/:id`
- `DELETE /api/admin/skills/:id`
- `PATCH /api/admin/skills/reorder`

Contoh reorder:
```http
PATCH /api/admin/skills/reorder
Content-Type: application/json

[
  {"id": "22222222-1111-1111-1111-111111111111", "order": 1},
  {"id": "22222222-3333-1111-1111-111111111111", "order": 2}
]
```

### Services
- `GET /api/admin/services`
- `POST /api/admin/services`
- `PUT /api/admin/services/:id`
- `DELETE /api/admin/services/:id`
- `PATCH /api/admin/services/reorder`
- `PATCH /api/admin/services/:id/toggle`

Contoh create:
```http
POST /api/admin/services
Content-Type: application/json

{
  "name": "Website Development",
  "description": "Pembuatan website custom dengan stack modern.",
  "price_min": 5000000,
  "price_max": 20000000,
  "currency": "IDR",
  "duration_label": "3-6 minggu",
  "is_active": true,
  "order": 1
}
```

### Projects
- `GET /api/admin/projects`
- `POST /api/admin/projects`
- `PUT /api/admin/projects/:id`
- `DELETE /api/admin/projects/:id`
- `PATCH /api/admin/projects/reorder`
- `PATCH /api/admin/projects/:id/feature`

Contoh set featured:
```http
PATCH /api/admin/projects/44444444-1111-1111-1111-111111111111/feature
Content-Type: application/json

{ "is_featured": true }
```

### Uploads (stub)
- `POST /api/admin/uploads`

Response saat ini:
```json
{ "data": { "message": "upload stub" } }
```

## 🔁 Workflow Pengembangan

```
make migrate   # jalankan migrasi database
make seed      # isi data contoh
go fmt ./...
go vet ./...
go test ./...
go build ./...
```

## ✅ Checklist & Catatan

- [x] CRUD Profile/Skills/Services/Projects + reorder/toggle/feature
- [x] Middleware guard stub untuk PR berikutnya
- [x] Unit test handler (200/400/401/403/404)
- [x] Upload stub siap integrasi storage

Risiko & tindak lanjut:
- Integrasi auth & role-based access akan dilanjutkan pada PR-3.
- Endpoint uploads masih stub; integrasi storage/S3 direncanakan di PR-5.
- Sinkronisasi knowledge base `/api/v1/knowledge-base` dengan data admin akan dilakukan pada PR-6.
