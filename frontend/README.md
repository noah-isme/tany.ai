# tany.ai Frontend

Antarmuka pengguna untuk prototipe **tany.ai**. Dibangun dengan Next.js 16, Tailwind CSS 4, dan difokuskan pada pengalaman chat
modern yang mengikuti blueprint pada dokumen utama proyek.

## ✨ Fitur Utama
- **Chat Assistant Terhubung Backend** – Komponen chat memanggil endpoint `POST /api/v1/chat`, menyertakan `chatId` lanjutan, dan menampilkan jawaban berbasis knowledge base aktual.
- **Dynamic Knowledge Snapshot** – Halaman utama dan bubble asistif menarik data dari `GET /api/v1/knowledge-base` sehingga perubahan di Admin langsung terlihat.
- **Admin Panel** – Shell admin lengkap untuk mengelola profil, skills, layanan, proyek, statistik placeholder, dan pengaturan tema/API key dengan invalidasi cache otomatis.

## 🛡️ Admin Panel

Panel admin dapat diakses melalui `/login` dan terlindungi oleh middleware JWT + role admin. Setelah login berhasil, pengguna diarahkan ke `/admin` dan mendapatkan navigasi sidebar dengan tab:

- **Dashboard** – Ringkasan profil, layanan, dan proyek.
- **Profil** – Formulir profil lengkap dengan pratinjau avatar.
- **Skills** – CRUD + drag & drop reorder.
- **Layanan** – CRUD, rentang harga, toggle visibilitas, dan reorder.
- **Proyek** – CRUD, tech stack dinamis, set featured, reorder, dan pratinjau gambar.
- **Statistik** – Placeholder insight sesuai blueprint README.
- **Settings** – Preferensi tema dan penyimpanan placeholder API keys (di-hash, tidak diekspose ke klien).

### ✅ Checklist Acceptance Criteria

- [x] Shell admin responsif lengkap (sidebar, header, skip-link, toggle tema tersimpan).
- [x] Dashboard menampilkan ringkasan profil, layanan, dan proyek.
- [x] Profil terisi dari API, validasi email/URL, dan pratinjau avatar.
- [x] Skills mendukung CRUD, drag & drop reorder, serta notifikasi status.
- [x] Layanan mendukung CRUD, validasi rentang harga, toggle visibilitas, dan reorder.
- [x] Proyek mendukung CRUD, tech stack dinamis, pratinjau gambar, dan set featured.
- [x] Statistik & Settings menampilkan placeholder sesuai blueprint dengan toggle tema persisten.

> **Catatan:** Mock backend sederhana tersedia untuk kebutuhan Playwright test. Untuk integrasi penuh gunakan backend Golang (PR-2/PR-3) dan set `API_BASE_URL` ke alamat server tersebut.

## 🧱 Struktur Direktori
```
app/                # Halaman Next.js (App Router)
components/chat/    # Komponen UI chat (bubble, window, input)
lib/                # Helper chat, knowledge fetcher, utilitas API
tests/              # Vitest unit test & Playwright e2e (mock backend)
```

## 🛠️ Perintah Penting
```bash
npm run dev     # Menjalankan server pengembangan di http://localhost:3000
npm run lint    # Menjalankan Next.js lint rules
npm test        # Menjalankan unit test Vitest + Testing Library
npm run test:e2e # Menjalankan Playwright (mock backend otomatis)
npm run build   # Build produksi Next.js
```

### Menjalankan Admin Panel

1. Set `API_BASE_URL` dan `JWT_SECRET` pada `.env` (contoh `http://localhost:8080` serta secret acak minimal 32 karakter) sesuai backend yang berjalan.
2. Jalankan `npm run dev` dan buka `http://localhost:3000/login`.
3. Masuk menggunakan kredensial admin backend (`admin@example.com / Admin#12345`).
4. Semua tab admin menggunakan server actions + token httpOnly sehingga data tidak terekspos di bundle klien.

## 🧪 Testing & Kualitas Kode
- Unit test menggunakan **Vitest** dan **Testing Library**.
- `setupTests.ts` mengaktifkan matcher `@testing-library/jest-dom`.
- Playwright e2e berada di `tests/e2e` (menggunakan mock backend otomatis).
- Konfigurasi linting mengikuti `eslint-config-next` (core web vitals + TypeScript).

## 🔗 Integrasi Backend
Frontend kini mengonsumsi endpoint backend `GET /api/v1/knowledge-base` dan `POST /api/v1/chat` secara langsung. Mock backend Playwright meniru struktur respons tersebut sehingga alur e2e (ubah data admin → knowledge base berubah → chat mengikuti) dapat diverifikasi tanpa backend Go asli.
