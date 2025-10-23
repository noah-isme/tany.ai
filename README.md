# 🤖 tany.ai - AI Client Chat Assistant

> Asisten chat AI yang menjadi representasi digital Anda untuk menjawab calon klien secara otomatis, profesional, dan 24/7.

---

## 📋 Daftar Isi

- [Tentang Aplikasi](#-tentang-aplikasi)
- [Konsep Aplikasi](#-konsep-aplikasi)
- [Fitur Utama](#️-fitur-utama)
- [Arsitektur & Stack Teknologi](#-arsitektur--stack-teknologi)
- [Arsitektur Sistem](#-arsitektur-sistem)
- [Desain UI/UX](#-desain-uiux)
- [Roadmap Pengembangan](#️-roadmap-pengembangan)

---

## 🎯 Tentang Aplikasi

**tany.ai** adalah chatbot berbasis AI yang dirancang khusus untuk freelancer dan profesional yang ingin memberikan respons cepat dan profesional kepada calon klien. Aplikasi ini menggunakan data profil, proyek, dan layanan Anda sebagai knowledge base untuk menjawab pertanyaan seperti:

- "Kamu bisa bikin website pakai teknologi apa?"
- "Berapa tarif untuk desain logo sederhana?"
- "Bisa kasih contoh proyek yang pernah kamu buat?"

**Keunggulan utama**: AI ini menjawab berdasarkan data Anda sendiri, bukan dari internet, sehingga lebih personal dan akurat.


---

## 🧭 Konsep Aplikasi

### Definisi
**AI Client Chat Assistant** adalah chatbot berbasis AI yang menjadi representasi digital Anda — menjawab pertanyaan dari calon klien secara otomatis dan profesional.

### Cara Kerja
AI ini mengambil jawaban dari:
- ✅ Data profil pribadi Anda
- ✅ Portfolio proyek yang pernah dikerjakan
- ✅ Daftar layanan dan tarif
- ❌ **BUKAN** dari internet atau sumber eksternal

### Tujuan Utama
Membantu freelancer menjawab calon klien dengan:
- ⚡ **Cepat** - Respons instan tanpa menunggu
- 💼 **Profesional** - Jawaban terstruktur dan informatif
- 🤖 **Otomatis** - Tersedia 24 jam sehari, 7 hari seminggu

---

## ⚙️ Fitur Utama

### 🎯 MVP (Minimum Viable Product)
*Estimasi waktu pengembangan: 2-3 minggu*

#### 1. Chat Interface (Frontend)
- Tampilan chat modern seperti ChatGPT
- Input teks dan balasan AI real-time
- Desain responsif untuk mobile dan desktop

#### 2. Knowledge Base (Backend)
Menyimpan informasi tentang freelancer:
- 👤 Bio dan deskripsi singkat
- 🛠️ Skill dan keahlian
- 💼 Layanan yang ditawarkan
- 📁 Portfolio proyek
- 💰 Daftar harga/tarif
- 📞 Informasi kontak

Fitur: Dashboard sederhana untuk mengubah data

#### 3. AI Chat Engine
- Menggunakan **GPT API** / **Gemini** untuk respons AI
- Sistem **prompt injection**: AI diarahkan hanya menjawab sesuai profil
- Context-aware responses

#### 4. Context Memory (Short-term)
- Menyimpan percakapan terakhir agar AI paham konteks
- Contoh: "berapa harganya?" setelah topik project dapat dipahami dengan benar

---

### 🚀 Fitur Lanjutan (Advanced / Portfolio Upgrade)

#### 1. Admin Dashboard
- **CRUD Operations**: Update profil, proyek, tarif, layanan
- **Analytics Dashboard**: 
  - Pertanyaan yang paling sering ditanyakan klien
  - Grafik traffic pengunjung
  - Tingkat engagement
- **Export Data**: Laporan dalam format CSV/PDF

#### 2. Lead Capture System
- Form otomatis muncul saat user tertarik
- Template: *"Tinggalkan email Anda untuk penawaran khusus"*
- Simpan leads ke database
- Notifikasi email otomatis ke Anda

#### 3. Voice Chat Mode
- **Speech-to-Text**: Gunakan Whisper API
- **Text-to-Speech**: Implementasi ElevenLabs atau Web Speech API
- Interface tombol microphone

#### 4. Multilingual Support
- Deteksi bahasa otomatis (Indonesia/Inggris)
- Respons AI dalam bahasa yang sama dengan user
- Support bahasa tambahan: Mandarin, Spanyol, Arab

#### 5. Integrasi ke Website Pribadi
- Widget chat di pojok kanan bawah
- Floating button dengan animasi
- Customizable theme dan branding
- Embed code yang mudah dipasang

#### 6. Advanced Features
- **Sentiment Analysis**: Deteksi mood calon klien
- **Auto-scheduling**: Integrasi dengan Calendly
- **File Sharing**: Upload brief project langsung di chat
- **Payment Gateway**: Link pembayaran otomatis

---

## 🧱 Arsitektur & Stack Teknologi

| Komponen | Teknologi Rekomendasi | Alternatif |
|----------|----------------------|------------|
| **Frontend** | React.js / Next.js | Vue.js, Svelte |
| **Chat UI** | React Chat UI / Custom Tailwind | Chakra UI, Material-UI |
| **Backend** | Golang (Gin) | Node.js + Express.js, NestJS |
| **Database** | MongoDB / Supabase | PostgreSQL, Firebase |
| **AI Layer** | OpenAI GPT-4 API | Gemini, Claude API |
| **Auth (Admin)** | Firebase Auth / JWT | Auth0, Clerk |
| **Deployment** | Vercel (frontend) + Render (backend) | Railway, DigitalOcean |
| **Voice (Optional)** | OpenAI Whisper API + Web Speech API | ElevenLabs, Google TTS |
| **Analytics** | Google Analytics / Mixpanel | PostHog, Amplitude |

### Alasan Pemilihan Stack

#### Frontend: Next.js
- ✅ SEO-friendly dengan Server-Side Rendering
- ✅ Fast page loading
- ✅ Built-in API routes
- ✅ Easy deployment ke Vercel

#### Backend: Golang (Gin)
- ✅ Performa tinggi dan latensi rendah (compiled language)
- ✅ Concurrency model ringan (goroutines) cocok untuk banyak koneksi WebSocket
- ✅ Binary tunggal mudah di-deploy dan skalabel
- ✅ Libraries HTTP/JSON dan middleware mature untuk integrasi OpenAI

#### Database: MongoDB
- ✅ Flexible schema untuk knowledge base
- ✅ Easy to scale
- ✅ JSON-like documents cocok untuk chat history
- ✅ Free tier di MongoDB Atlas

#### AI: OpenAI GPT-4
- ✅ Kualitas respons terbaik
- ✅ Dokumentasi lengkap
- ✅ Support context memory
- ✅ Customizable dengan system prompts

---

## 🧩 Arsitektur Sistem

### Diagram Konseptual
```
User (Client)
     ↓
┌────────────────┐
│ React Chat UI  │ ←→  Web Socket (Real-time)
└────────────────┘
     ↓
┌─────────────────────────┐
│   API Server            │
│   (Golang + Gin)        │
└─────────────────────────┘
     ↓
┌─────────────────────────┐
│  AI Processor Layer     │
│  - Prompt Engineering   │
│  - Context Management   │
└─────────────────────────┘
     ↓                ↓
┌──────────────┐  ┌──────────────┐
│ Knowledge    │  │  OpenAI API  │
│ Base (DB)    │  │  (GPT-4)     │
└──────────────┘  └──────────────┘
```

### Flow Diagram Detail
```
1. User mengetik pertanyaan
   ↓
2. Frontend mengirim request ke API Server
   ↓
3. API Server mengambil data dari Knowledge Base
   ↓
4. AI Processor membuat prompt:
   "Kamu adalah asisten virtual dari [Nama Freelancer].
    Jawab pertanyaan hanya berdasarkan data di bawah ini:
    [Data profil, layanan, proyek, tarif]"
   ↓
5. Request dikirim ke OpenAI API
   ↓
6. OpenAI mengembalikan respons
   ↓
7. Backend menyimpan chat history ke database
   ↓
8. Frontend menampilkan respons ke user
```

### Knowledge Base Structure
```json
{
  "profile": {
    "name": "John Doe",
    "title": "Full Stack Developer",
    "bio": "Experienced developer with 5+ years...",
    "contact": {
      "email": "john@example.com",
      "phone": "+62812345678",
      "location": "Jakarta, Indonesia"
    }
  },
  "skills": [
    "React.js",
    "Golang",
    "MongoDB",
    "UI/UX Design"
  ],
  "services": [
    {
      "name": "Website Development",
      "description": "Custom website with modern tech",
      "price_range": "Rp 5.000.000 - Rp 20.000.000",
      "duration": "2-4 minggu"
    }
  ],
  "projects": [
    {
      "title": "E-commerce Platform",
      "description": "Full-featured online store",
  "tech_stack": ["React", "Golang", "Stripe"],
      "image_url": "https://...",
      "project_url": "https://..."
    }
  ]
}
```

---

## 🎨 Desain UI/UX

### 💬 Halaman Chat Utama

#### Header Section
```
┌─────────────────────────────────────────┐
│  🤖 Tanya Asisten [Nama Kamu]          │
│  ● Online - Biasanya membalas cepat    │
└─────────────────────────────────────────┘
```

#### Chat Area
```
┌─────────────────────────────────────────┐
│  Bot: Halo! Saya asisten virtual        │
│       [Nama]. Ada yang bisa saya        │
│       bantu? 😊                         │
│                                         │
│                User: Halo, layanan apa  │
│                      yang kamu tawarkan?│
│                                         │
│  Bot: Saya menawarkan beberapa         │
│       layanan:                          │
│       • Website Development            │
│       • Mobile App Development         │
│       • UI/UX Design                   │
│       Mau tahu detail salah satunya?   │
└─────────────────────────────────────────┘
```

#### Input Bar
```
┌─────────────────────────────────────────┐
│  💬 Ketik pertanyaan Anda...      [📎] [🎤]│
└─────────────────────────────────────────┘
```

#### Placeholder Examples
Contoh pertanyaan yang muncul sebagai suggestion:
- 💼 "Layanan apa yang kamu tawarkan?"
- 💰 "Berapa tarif pembuatan website?"
- 📁 "Bisa kasih contoh portofolio?"
- ⏱️ "Berapa lama waktu pengerjaannya?"
- 📞 "Bagaimana cara menghubungi kamu?"

---

### ⚙️ Admin Dashboard

#### Layout Structure
```
┌──────────────────────────────────────────┐
│  SIDEBAR        │  MAIN CONTENT          │
│                 │                        │
│  📊 Dashboard   │  [Content Area]       │
│  👤 Profil      │                        │
│  💼 Layanan     │                        │
│  📁 Proyek      │                        │
│  📈 Statistik   │                        │
│  ⚙️ Settings    │                        │
└──────────────────────────────────────────┘
```

#### Tab 1: Dashboard Overview
- 📊 Total conversations
- 👥 Unique visitors
- 🔥 Popular questions
- 📧 Leads captured

#### Tab 2: Profil
- Edit nama
- Edit deskripsi singkat
- Update skills
- Ubah informasi kontak
- Upload foto profil

#### Tab 3: Layanan
- CRUD operations untuk layanan
- Kolom: Nama layanan, deskripsi, harga, durasi
- Toggle visibility (aktif/nonaktif)

#### Tab 4: Proyek
- Upload gambar proyek
- Deskripsi singkat
- Tech stack yang digunakan
- Link ke demo/repository
- Kategori proyek

#### Tab 5: Statistik
Grafik dan data:
- 📊 Pertanyaan paling populer (bar chart)
- 📈 Traffic harian (line chart)
- ⏰ Peak hours pengunjung
- 🌍 Lokasi pengunjung (map)
- 📱 Device breakdown (mobile vs desktop)

---

### 🧭 Integrasi ke Website

#### Chat Widget (Floating Button)
```
                                    ┌─────┐
                                    │ 💬  │
                                    │ Chat│
                                    └─────┘
```

#### Pop-up Chat Window
```
      ┌────────────────────────────┐
      │ 🤖 Tanya Asisten          ×│
      ├────────────────────────────┤
      │                            │
      │  [Chat Area]               │
      │                            │
      ├────────────────────────────┤
      │ 💬 Ketik pesan...          │
      └────────────────────────────┘
```

#### Embed Code
```html
<!-- Copy paste code ini ke website Anda -->
<script src="https://tany.ai/widget.js"></script>
<script>
  TanyAI.init({
    userId: 'your-user-id',
    theme: 'light', // or 'dark'
    position: 'bottom-right',
    greeting: 'Halo! Ada yang bisa dibantu?'
  });
</script>
```

---

## 🗓️ Roadmap Pengembangan (6 Minggu)

### Minggu 1: Foundation & Design
**Target:**
- ✅ Riset kompetitor dan user flow
- ✅ Desain UI/UX di Figma (chat interface + admin dashboard)
- ✅ Setup project structure
  - Initialize Next.js project
  - Setup Golang backend (Gin) scaffold
  - Configure MongoDB connection (or preferred DB)
- ✅ Setup version control (Git)

**Deliverables:**
- Figma prototype
- Project boilerplate ready

---

### Minggu 2: Core Chat Functionality
**Target:**
- ✅ Buat UI chat interface
- ✅ Implementasi chat bubble components
- ✅ Koneksi OpenAI API
- ✅ Setup environment variables
- ✅ Buat API endpoint untuk chat

**Deliverables:**
- Working chat interface
- AI respons berfungsi (basic)

---

### Minggu 3: Knowledge Base Implementation
**Target:**
- ✅ Design database schema untuk knowledge base
- ✅ Implement static JSON untuk testing
- ✅ Buat prompt engineering system
- ✅ Test berbagai skenario pertanyaan
- ✅ Fine-tune AI responses

**Deliverables:**
- Knowledge base structure
- AI menjawab berdasarkan data profil

---

### Minggu 4: Admin Panel Development
**Target:**
- ✅ Buat authentication system
- ✅ CRUD untuk profil
- ✅ CRUD untuk layanan dan tarif
- ✅ CRUD untuk proyek/portfolio
- ✅ Image upload functionality

**Deliverables:**
- Functional admin dashboard
- Data management system

---

### Minggu 5: Advanced Features
**Target:**
- ✅ Implementasi analytics dashboard
- ✅ Lead capture system
- ✅ Email notification setup
- ✅ Chat history storage
- ✅ Context memory improvement

**Deliverables:**
- Analytics showing visitor data
- Lead capture form working

---

### Minggu 6: Polish & Deployment
**Target:**
- ✅ UI/UX improvements
- ✅ Bug fixing dan testing
- ✅ Performance optimization
- ✅ Deploy frontend ke Vercel
- ✅ Deploy backend ke Render
- ✅ Setup custom domain
- ✅ Buat demo video untuk portfolio

**Deliverables:**
- Fully functional production app
- Demo video
- Documentation

---

## 📊 Metrik Kesuksesan

### KPI (Key Performance Indicators)

1. **Response Time**
   - Target: < 3 detik per respons
   - Measure: Average API response time

2. **Accuracy Rate**
   - Target: 85%+ pertanyaan terjawab dengan benar
   - Measure: User feedback (thumbs up/down)

3. **Lead Conversion**
   - Target: 15%+ visitors meninggalkan email
   - Measure: Lead capture form submissions

4. **User Engagement**
   - Target: Average 4+ messages per session
   - Measure: Chat session analytics

5. **Availability**
   - Target: 99%+ uptime
   - Measure: Server monitoring tools

---

## 🚀 Next Steps Setelah MVP

### Phase 2: Enhancement (Bulan 2-3)
- [ ] Voice chat integration
- [ ] Multilingual support
- [ ] Advanced analytics
- [ ] A/B testing untuk respons AI
- [ ] Mobile app (React Native)

### Phase 3: Monetization (Bulan 4-6)
- [ ] Subscription tiers
- [ ] White-label solution untuk agency
- [ ] API access untuk developer
- [ ] Template marketplace

---

## 💡 Tips Pengembangan

### Best Practices
1. **Version Control**: Commit code setiap selesai feature
2. **Testing**: Unit test untuk setiap endpoint
3. **Documentation**: Tulis API documentation lengkap
4. **Security**: Jangan hardcode API keys
5. **Monitoring**: Setup error tracking (Sentry)

### Common Pitfalls to Avoid
- ❌ Over-engineering di awal
- ❌ Tidak test dengan user real
- ❌ Mengabaikan mobile responsive
- ❌ Lupa setup analytics dari awal
- ❌ Tidak membuat backup database

---

## 📚 Resources & References

### Documentation
- [OpenAI API Docs](https://platform.openai.com/docs)
- [Next.js Documentation](https://nextjs.org/docs)
- [MongoDB Manual](https://docs.mongodb.com)
- [Tailwind CSS](https://tailwindcss.com/docs)

### Inspirasi Design
- ChatGPT interface
- Intercom chat widget
- Drift conversational marketing
- Tidio live chat

### Tutorial Recommended
- Next.js full-stack tutorial
- OpenAI API integration
- WebSocket for real-time chat
- JWT authentication

---

## 📞 Support & Contact

Untuk pertanyaan atau kolaborasi:
- 📧 Email: [your-email@example.com]
- 💼 LinkedIn: [Your LinkedIn]
- 🌐 Portfolio: [your-portfolio.com]

---

## 📄 License

MIT License - Feel free to use this concept for your own project!

---

**Last Updated**: October 22, 2025

**Version**: 1.0.0

**Status**: 🚀 Ready for Development

---

## 🗄️ Running with Database (Local & CI)

1. Salin konfigurasi contoh lalu sesuaikan kredensial Supabase/Postgres Anda:
   ```bash
   cp backend/.env.example backend/.env
   ```
2. Pastikan variabel `POSTGRES_URL` mengarah ke database yang bisa diakses.
3. Jalankan migrasi dan seeder contoh:
   ```bash
   make -C backend migrate
   make -C backend seed
   ```
4. Start server backend dengan konfigurasi yang sama:
   ```bash
   go run ./backend/cmd/api
   ```
5. Cek status API dan database melalui endpoint kesehatan:
   ```bash
   curl http://localhost:8080/healthz
   ```

Di lingkungan CI, workflow `ci-be.yml` akan menjalankan rangkaian yang sama (migrate → seed → test) menggunakan service Postgres.
