# AI Personalization & Profile Embedding

Fase ini menambahkan memori semantik ke Tany.AI sehingga jawaban dapat meniru gaya tulisan dan prioritas layanan secara konsisten. Dokumen ini menjelaskan arsitektur, pipeline reindexing, serta prosedur QA untuk fitur personalisasi.

## Ringkasan Arsitektur

- **Penyimpanan vektor**: PostgreSQL dengan ekstensi `pgvector`. Tabel utama `embeddings` menyimpan konten, metadata, dan vektor `vector(1536)`.
- **Konfigurasi**: Tabel `embedding_config` menyimpan bobot personalisasi serta timestamp reindex/reset.
- **Service backend**: Paket `internal/embedding` mengelola repository, caching, similarity search, dan admin handler.
- **Provider embedding**: Implementasi `internal/ai/openai_embeddings.go` menggunakan model OpenAI (default `text-embedding-3-large`).
- **Prompt builder**: `prompt.BuildPersonalizedPrompt` menambahkan instruksi persona berdasarkan hasil pencarian vektor.
- **Panel Admin**: `/admin/personalization` untuk memantau status, mengatur bobot, memicu reindex/reset.

## Variabel Lingkungan

| Variabel | Default | Deskripsi |
| --- | --- | --- |
| `ENABLE_PERSONALIZATION` | `false` | Mengaktifkan pipeline personalisasi. |
| `EMBEDDING_PROVIDER` | `openai` | Provider embedding yang digunakan. |
| `EMBEDDING_MODEL` | `text-embedding-3-large` | Model embedding saat `EMBEDDING_PROVIDER=openai`. |
| `EMBEDDING_DIM` | `1536` | Dimensi vektor. Harus cocok dengan model provider. |
| `EMBEDDING_CACHE_TTL` | `24h` | TTL cache in-memory untuk hasil similarity. |
| `PERSONALIZATION_WEIGHT` | `0.65` | Bobot default instruksi persona di prompt. |
| `OPENAI_API_KEY` | — | Wajib diisi jika menggunakan provider OpenAI. |

## Pipeline Reindexing

1. Admin membuka `/admin/personalization` dan menekan **Reindex embedding**.
2. Handler backend memuat snapshot knowledge base terkini dan memanggil `embedding.Service.Reindex`.
3. Service menghapus embedding lama, menghasilkan vektor baru (profil, layanan, proyek, post), dan menyimpannya ke tabel `embeddings`.
4. Timestamp `lastReindexedAt` diperbarui di `embedding_config`. Panel admin otomatis menampilkan status terbaru setelah halaman direfresh.

Untuk menjalankan reindex manual dari terminal (misal pada skrip):

```bash
go run ./cmd/server \
  -enable-personalization=true \
  # jalankan server lalu trigger endpoint POST /api/admin/personalization/reindex
```

## Reset Embedding

Jika perlu menghapus seluruh embedding (misal pergantian model), gunakan tombol **Reset embedding** di panel admin. API akan:

1. Menghapus seluruh baris `embeddings` dengan `vector IS NOT NULL`.
2. Mengosongkan cache layanan.
3. Menyimpan `lastResetAt` di `embedding_config`.

Setelah reset, lakukan reindex untuk mengisi ulang vektor.

## Bobot Personalisasi

- Bobot 0 → Persona tidak digunakan; prompt kembali ke gaya umum.
- Bobot 1 → Instruksi persona dominan; cocok untuk percakapan sangat personal.
- UI slider menyimpan nilai dengan presisi dua angka desimal.
- Penyesuaian bobot disimpan di database dan langsung digunakan oleh request berikutnya.

## QA Checklist

1. **Reindex sukses**: Tombol reindex mengembalikan pesan sukses dan jumlah embedding bertambah (cek admin + log DB).
2. **Prompt**: Saat chat, inspeksi `prompt` pada response API `/api/v1/chat` dan pastikan terdapat blok *Instruksi personalisasi* ketika fitur aktif.
3. **Fallback**: Nonaktifkan `ENABLE_PERSONALIZATION` dan pastikan slider terkunci, prompt kembali ke format lama.
4. **Latency**: Pastikan waktu respon tidak melebihi 2x baseline pada 5 percobaan berurutan.
5. **Akurasi**: Bandingkan 10 sample percakapan sebelum/sesudah personalisasi, cari peningkatan tone & relevansi minimal 20% berdasarkan rubric QA.

## Troubleshooting

- **OpenAI API Key kosong**: Personalizer otomatis nonaktif, cek log `[warn] OPENAI_API_KEY missing`. Isi environment dan restart server.
- **Dimension mismatch**: Pastikan `EMBEDDING_DIM` sesuai dengan panjang vektor model provider.
- **Slow queries**: Pastikan index `embeddings_vector_idx` terbuat (`migrate` terakhir). Sesuaikan parameter `lists` pada indeks IVFFLAT jika dataset besar.

## Referensi

- [pgvector documentation](https://github.com/pgvector/pgvector)
- [OpenAI embeddings API](https://platform.openai.com/docs/api-reference/embeddings)

