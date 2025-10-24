# tany.ai Frontend

Antarmuka pengguna untuk prototipe **tany.ai**. Dibangun dengan Next.js 16, Tailwind CSS 4, dan difokuskan pada pengalaman chat
modern yang mengikuti blueprint pada dokumen utama proyek.

## ✨ Fitur Utama
- **Chat Simulation** – Komponen chat interaktif dengan mock response yang menggambarkan alur integrasi OpenAI GPT.
- **Knowledge Base Snapshot** – Ringkasan layanan, portofolio, serta kontak yang dipakai untuk membangun context AI.
- **System Prompt Viewer** – Memperlihatkan prompt dasar yang akan digunakan oleh model bahasa di sisi backend.
- **Admin Panel** – Shell admin lengkap untuk mengelola profil, skills, layanan, proyek, statistik placeholder, dan pengaturan tema/API key.

## 🛡️ Admin Panel

Panel admin dapat diakses melalui `/login` dan terlindungi oleh middleware JWT + role admin. Setelah login berhasil, pengguna diarahkan ke `/admin` dan mendapatkan navigasi sidebar dengan tab:

- **Dashboard** – Ringkasan profil, layanan, dan proyek.
- **Profil** – Formulir profil lengkap dengan pratinjau avatar.
- **Skills** – CRUD + drag & drop reorder.
- **Layanan** – CRUD, rentang harga, toggle visibilitas, dan reorder.
- **Proyek** – CRUD, tech stack dinamis, set featured, reorder, dan pratinjau gambar.
- **Statistik** – Placeholder insight sesuai blueprint README.
- **Settings** – Preferensi tema dan penyimpanan placeholder API keys (di-hash, tidak diekspose ke klien).

> **Catatan:** Mock backend sederhana tersedia untuk kebutuhan Playwright test. Untuk integrasi penuh gunakan backend Golang (PR-2/PR-3) dan set `API_BASE_URL` ke alamat server tersebut.

## 🧱 Struktur Direktori
```
app/                # Halaman Next.js (App Router)
components/chat/    # Komponen UI chat (bubble, window, input)
data/               # Knowledge base statis sesuai profil Tanya A.I.
lib/                # Helper chat dan utilitas mock API
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
Mock response pada frontend mengikuti struktur knowledge base yang sama dengan layanan Golang (Gin) sehingga integrasi API di
fase selanjutnya dapat dilakukan secara mulus.
