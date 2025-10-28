# Analytics & Insight Dashboard

Modul analytics menyediakan observabilitas real-time untuk performa chat Tany.AI. Dokumen ini menjelaskan alur data, konfigurasi environment, dan cara pengujian end-to-end.

## Arsitektur

```
/chat request
    │
    ├─ Gin middleware + ChatHandler
    │   ├─ Generate response via provider (Gemini, Leapcell, Mock)
    │   ├─ Persist ke `chat_history`
    │   └─ Record event -> analytics service
    │
    ├─ analytics.Service
    │   ├─ Simpan granular event ke `analytics_events`
    │   └─ Upsert snapshot harian ke `analytics_summary`
    │
    └─ Admin dashboard (Next.js 16 + Recharts)
        ├─ /api/admin/analytics/summary
        ├─ /api/admin/analytics/events
        └─ /api/admin/analytics/leads
```

## Skema Database

### `analytics_events`

| Kolom        | Tipe        | Deskripsi                                     |
|--------------|-------------|-----------------------------------------------|
| `id`         | `uuid`      | Primary key                                   |
| `timestamp`  | `timestamptz` | Waktu event                                  |
| `event_type` | `text`      | Jenis event (`chat`, `lead`, dll)             |
| `source`     | `text`      | Sumber request (header `X-Chat-Source`)       |
| `provider`   | `text`      | Provider AI yang digunakan                    |
| `duration_ms`| `int`       | Latensi dalam milidetik                       |
| `success`    | `bool`      | Status sukses provider                        |
| `user_agent` | `text`      | User-Agent klien                              |
| `metadata`   | `jsonb`     | Payload tambahan (cache hit, ip, dsb)         |

### `analytics_summary`

Snapshot harian dengan breakdown provider. Kolom `provider_breakdown` berupa objek JSON `{ provider: { totalChats, avgResponseTime, successRate } }`.

## Environment Variables

| Variabel                     | Default | Deskripsi                                           |
|------------------------------|---------|-----------------------------------------------------|
| `ENABLE_ANALYTICS`           | `false` | Aktifkan/Nonaktifkan pencatatan analytics           |
| `ANALYTICS_RETENTION_DAYS`   | `90`    | Menyimpan ringkasan selama N hari                   |
| `PROMETHEUS_PORT`            | `9090`  | Tersedia untuk integrasi metrik lanjutan (opsional) |

## API Endpoints

Semua endpoint berada di bawah `/api/admin/analytics` dan membutuhkan autentikasi admin.

- `GET /summary` — ringkasan periode (filter: `from`, `to`, `source`, `provider`).
- `GET /events` — daftar event granular, mendukung pagination (`page`, `limit`) dan filter `type`.
- `GET /leads` — alias `events` dengan `event_type = lead`.

## Frontend Dashboard

Halaman `/admin/analytics` menampilkan:

1. **KPI Cards**: total chat, rata-rata latensi, success rate, conversion rate.
2. **Tren Harian**: line chart (Recharts) untuk volume vs success rate.
3. **Performa Provider**: pie chart + legenda.
4. **Lead & Conversion**: bar chart + tabel lead terbaru.
5. **Log Interaksi**: tabel event chat beserta status provider.

Filter rentang tanggal, provider, dan sumber tersedia di panel atas. Data auto-refresh setiap 5 menit.

## Testing

### Backend

```bash
go test ./internal/analytics/...
```

### Frontend

```bash
npm test -- Analytics
```

### E2E (opsional)

Gunakan Playwright tag `@analytics` setelah menambahkan skenario di suite E2E.

## Troubleshooting

- **Analytics tidak aktif**: pastikan `ENABLE_ANALYTICS=true` di environment server.
- **Grafik kosong**: cek apakah rentang tanggal tidak melebihi retensi atau belum ada event.
- **Latency selalu 0**: verifikasi `chat_handler` menyuplai `RecordChat` dengan durasi.

Untuk insight lanjutan (export ke Prometheus, dsb) gunakan data `analytics_summary` sebagai sumber ETL.
