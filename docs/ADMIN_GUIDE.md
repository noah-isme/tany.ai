# Panduan Admin tany.ai

Panduan ini ditujukan untuk operator non-developer yang mengelola knowledge base tany.ai sebelum rilis v1.0.0. Semua langkah mengacu pada panel admin bawaan (`/admin`).

## 1. Masuk ke Panel Admin
1. Buka `https://tanyai.app/login` (atau URL staging Anda).
2. Masukkan email & password yang diberikan saat provisioning. Default dari seed lokal adalah `admin@example.com / Password123!`.
3. Setelah berhasil login, Anda akan diarahkan ke dashboard admin. Middleware akan mengarahkan ulang ke halaman login jika token hilang atau kadaluarsa, dan ke halaman 403 bila akun tidak memiliki role `admin`.【F:frontend/middleware.ts†L16-L60】

## 2. Ringkasan Dashboard
- Halaman `/admin` menampilkan status singkat knowledge base (profil terkahir diperbarui, jumlah layanan aktif, proyek featured, dsb.). Komponen sidebar mengingatkan bahwa seluruh data di panel menjadi sumber jawaban AI.【F:frontend/components/admin/AdminSidebar.tsx†L76-L81】
- Setiap modul (Profil, Skills, Services, Projects) dapat diakses melalui menu di sisi kiri.

## 3. Mengelola Profil & Kontak
1. Masuk ke `/admin/profile`.
2. Lengkapi nama, jabatan, bio singkat, email, telepon, lokasi, dan URL avatar.
3. Klik **Simpan perubahan**. Form memvalidasi input menggunakan skema Zod dan menampilkan pesan sukses atau error per-field jika ada kendala.【F:frontend/components/admin/ProfileForm.tsx†L1-L158】
4. Avatar dapat diisi dengan menempelkan URL yang valid atau melalui uploader bawaan (lihat bagian Upload di bawah).

## 4. Mengelola Skills
1. Masuk ke `/admin/skills`.
2. Tambahkan skill baru via field input lalu tekan **Tambah**.
3. Gunakan tombol panah/drag & drop untuk mengatur urutan prioritas. Urutan paling atas akan muncul lebih dulu di prompt AI.【F:frontend/components/admin/SkillsManager.tsx†L1-L210】
4. Tekan ikon hapus untuk menghilangkan skill. Konfirmasi dialog akan muncul sebelum data dihapus.【F:frontend/components/admin/SkillsManager.tsx†L120-L171】

## 5. Mengelola Services
1. Masuk ke `/admin/services` lalu klik **Tambah Layanan**.
2. Isi nama, durasi, deskripsi singkat, harga minimum/maksimum, dan mata uang. Toggle **Aktif** menentukan apakah layanan tampil di knowledge base publik.【F:frontend/components/admin/ServicesManager.tsx†L209-L420】
3. Setelah tersimpan, layanan baru muncul di tabel. Gunakan tombol drag untuk mengatur urutan; perubahan akan tersimpan otomatis ke backend dan memicu invalidasi cache knowledge base.【F:frontend/components/admin/ServicesManager.tsx†L48-L103】
4. Toggle status layanan untuk menyembunyikan/menampilkan di chat publik. Endpoint backend hanya mengembalikan layanan `is_active = TRUE`, jadi pastikan minimal satu layanan aktif agar snippet chat tetap informatif.【F:backend/internal/services/kb/aggregator.go†L86-L128】

## 6. Mengelola Projects
1. Masuk ke `/admin/projects`.
2. Klik **Tambah Proyek** untuk membuka editor baru, lengkapi judul, deskripsi, kategori, URL, serta tech stack (bisa menambah banyak item).【F:frontend/components/admin/ProjectsManager.tsx†L1-L480】
3. Gunakan toggle **Featured** agar proyek muncul di highlight chat & landing page terlebih dahulu.【F:backend/internal/services/kb/aggregator.go†L130-L181】
4. Seret baris tabel untuk mengatur urutan tampilan; gunakan tombol **Edit** atau **Hapus** sesuai kebutuhan.

## 7. Upload Avatar & Gambar Proyek
- Komponen **ImageUploader** menerima file drag & drop atau klik untuk memilih file. Setelah upload berhasil, field otomatis terisi dengan URL publik storage.【F:frontend/components/admin/ImageUploader.tsx†L20-L108】
- Backend membatasi ukuran file (`UPLOAD_MAX_MB`, default 5MB), memverifikasi MIME dan mendukung JPEG/PNG/WebP/SVG (SVG harus diaktifkan eksplisit). Payload disanitasi agar bebas script sebelum disimpan.【F:backend/internal/handlers/admin/uploads.go†L45-L168】【F:backend/internal/config/env.go†L85-L206】
- Jika upload gagal, periksa pesan error pada panel. Log API juga menuliskan metadata (mime, size, key, latency) untuk debugging.【F:backend/internal/handlers/admin/uploads.go†L169-L206】

## 8. Best Practice Operasional
- **Jangan menonaktifkan semua layanan**: chat snippet hanya menampilkan layanan aktif. Minimal satu layanan harus aktif agar calon klien mendapat konteks.【F:backend/internal/services/kb/aggregator.go†L86-L128】
- **Gunakan urutan untuk prioritas**: layanan & proyek dengan urutan lebih rendah muncul lebih awal di prompt. Manfaatkan drag & drop setelah mengisi data baru.【F:frontend/components/admin/ServicesManager.tsx†L48-L103】
- **Perhatikan rate limit**: login dan upload dibatasi per menit. Jika menemukan pesan error `too many requests`, tunggu beberapa detik sebelum mencoba lagi.【F:backend/internal/server/server.go†L58-L116】【F:backend/internal/auth/ratelimit.go†L12-L68】
- **Gunakan bahasa profesional**: semua teks yang diisikan akan muncul dalam jawaban AI tanpa modifikasi tambahan.

## 9. Troubleshooting Umum
| Gejala | Penyebab | Solusi |
| --- | --- | --- |
| Upload gagal dengan pesan `unsupported media type` | Tipe file tidak ada di whitelist atau SVG belum diizinkan | Pastikan file JPEG/PNG/WebP, atau aktifkan `ALLOW_SVG=true` di environment backend jika perlu.【F:backend/internal/handlers/admin/uploads.go†L103-L140】【F:backend/internal/config/env.go†L172-L201】 |
| Mendapat `Sesi berakhir` saat berpindah halaman admin | JWT access token kadaluarsa atau tidak valid | Login ulang. Middleware otomatis menghapus cookie dan mengarahkan ke halaman login ketika token invalid.【F:frontend/middleware.ts†L24-L60】 |
| Chat publik belum menampilkan perubahan terbaru | Cache knowledge base masih dalam TTL | Tunggu ±1 menit atau lakukan perubahan tambahan (misal update layanan) untuk memicu invalidasi manual.【F:backend/internal/services/kb/aggregator.go†L21-L68】 |
| Login ditolak meski kredensial benar | Rate limit login tercapai (5/min default) | Tunggu 1-2 menit sebelum mencoba kembali atau tingkatkan nilai `LOGIN_RATE_LIMIT_PER_MIN` jika diperlukan di lingkungan staging.【F:backend/internal/config/env.go†L133-L214】 |

Jika kendala berlanjut, hubungi tim developer dengan menyertakan timestamp, endpoint, dan payload yang digunakan agar dapat ditelusuri melalui log JSON backend.【F:backend/internal/middleware/logger.go†L1-L23】

