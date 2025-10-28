# Integrasi Konten noahis.me

Integrasi ini memungkinkan Tany.AI menarik proyek, layanan, dan artikel terbaru dari situs [noahis.me](https://noahis.me) secara otomatis. Konten yang berhasil disinkronkan akan muncul di basis pengetahuan serta dapat digunakan builder prompt saat merespons percakapan.

## Alur singkat
1. **Sumber eksternal** disimpan di tabel `external_sources`.
2. **Ingestion service** mengambil sitemap dan halaman HTML, menghormati `robots.txt`, ETag, dan header `Last-Modified`.
3. **Konten ter-normalisasi** disimpan di `external_items` dan disatukan dengan proyek/layanan internal saat membangun blok prompt.
4. Admin dapat memicu sinkronisasi manual melalui UI maupun CLI.

## Variabel lingkungan
| Nama | Deskripsi | Contoh |
| ---- | --------- | ------ |
| `EXTERNAL_SOURCES_DEFAULT` | Daftar sumber bawaan dalam format JSON. Akan dibuat/diupdate ketika API dijalankan atau perintah sinkronisasi dipanggil. | `[{"name":"noahis.me","base_url":"https://noahis.me","type":"auto","enabled":true}]` |
| `HTTP_TIMEOUT_MS` | Timeout HTTP ingestion (ms). Default `8000`. | `12000` |
| `EXTERNAL_RATE_LIMIT_RPM` | Batas request per menit saat crawling. Default `30`. | `45` |
| `EXTERNAL_DOMAIN_ALLOWLIST` | Daftar domain yang diizinkan. Gunakan koma sebagai pemisah. | `noahis.me,www.noahis.me` |

> **Catatan:** Jika `EXTERNAL_SOURCES_DEFAULT` tidak di-set, aplikasi otomatis memakai konfigurasi default (`noahis.me`). Pastikan string JSON valid (gunakan kutip ganda) agar parsing tidak gagal saat start-up.

## Menjalankan sinkronisasi manual
### Via CLI
1. Pastikan database sudah dimigrasi dan berisi kredensial admin.
2. Eksekusi perintah berikut di akar repo:
   ```bash
   make external-sync
   ```
3. Output JSON berisi ringkasan status setiap sumber:
   ```json
   {
     "completedAt": "2024-12-12T12:00:00Z",
     "results": [
       {"id": "...", "name": "noahis.me", "status": "ok", "items": 12}
     ]
   }
   ```

Perintah ini menggunakan konfigurasi yang sama dengan aplikasi utama, termasuk env `POSTGRES_URL` dan `JWT_SECRET`. Pastikan variabel tersebut tersedia (gunakan `.env` backend bila perlu).

### Via UI Admin
1. Masuk ke `/admin/integrations` dengan akun admin.
2. Klik **"Sinkron sekarang"** pada sumber yang diinginkan.
3. Setelah sukses, tabel konten akan menampilkan item terbaru. Anda dapat mematikan visibilitas tiap item lewat toggle tanpa menghapus data.

## Penjadwalan otomatis
Workflow GitHub Actions `external-sync.yml` menjalankan `make external-sync` terjadwal. Set `POSTGRES_URL` dan `JWT_SECRET` sebagai secret repository (`PROD_POSTGRES_URL`, `PROD_JWT_SECRET` misalnya) lalu mapping ke environment workflow. Workflow akan gagal bila sinkronisasi error sehingga dapat dipantau lewat notifikasi GitHub.

## Troubleshooting
- **Sinkronisasi melewatkan konten baru:** pastikan server remote mengirim ETag/Last-Modified. Gunakan `make external-sync` dengan `LOG_LEVEL=debug` (opsional) untuk melihat log terperinci.
- **Permintaan diblokir:** cek `EXTERNAL_DOMAIN_ALLOWLIST` serta `robots.txt`. Ingestion tidak akan mengakses domain yang tidak ada di allowlist atau dilarang oleh robots.
- **HTML berformat kompleks:** Sanitizer berbasis [bluemonday](https://github.com/microcosm-cc/bluemonday) menghapus script/iframe. Jika informasi penting hilang, pertimbangkan untuk memperkaya metadata melalui `metadata` JSON.

## Referensi tabel
```sql
-- external_sources
id UUID PRIMARY KEY
name TEXT NOT NULL
base_url TEXT NOT NULL UNIQUE
source_type TEXT NOT NULL DEFAULT 'auto'
enabled BOOLEAN NOT NULL DEFAULT true
etag TEXT
last_modified TIMESTAMPTZ
last_synced_at TIMESTAMPTZ
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

-- external_items
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
source_id UUID NOT NULL REFERENCES external_sources(id)
kind TEXT NOT NULL
title TEXT NOT NULL
url TEXT NOT NULL
summary TEXT
content TEXT
metadata JSONB NOT NULL DEFAULT '{}'
published_at TIMESTAMPTZ
hash TEXT NOT NULL
visible BOOLEAN NOT NULL DEFAULT true
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
```
