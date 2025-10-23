# tany.ai Frontend

Antarmuka pengguna untuk prototipe **tany.ai**. Dibangun dengan Next.js 16, Tailwind CSS 4, dan difokuskan pada pengalaman chat
modern yang mengikuti blueprint pada dokumen utama proyek.

## ✨ Fitur Utama
- **Chat Simulation** – Komponen chat interaktif dengan mock response yang menggambarkan alur integrasi OpenAI GPT.
- **Knowledge Base Snapshot** – Ringkasan layanan, portofolio, serta kontak yang dipakai untuk membangun context AI.
- **System Prompt Viewer** – Memperlihatkan prompt dasar yang akan digunakan oleh model bahasa di sisi backend.

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
npm run build   # Build produksi Next.js
```

## 🧪 Testing & Kualitas Kode
- Unit test menggunakan **Vitest** dan **Testing Library**.
- `setupTests.ts` mengaktifkan matcher `@testing-library/jest-dom`.
- Konfigurasi linting mengikuti `eslint-config-next` (core web vitals + TypeScript).

## 🔗 Integrasi Backend
Mock response pada frontend mengikuti struktur knowledge base yang sama dengan layanan Golang (Gin) sehingga integrasi API di
fase selanjutnya dapat dilakukan secara mulus.
