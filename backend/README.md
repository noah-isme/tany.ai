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
