# **Tau-Tau Run Registration System Constitution**

<!-- Spec Constitution -->

## Core Principles

### I. MVP-First, Event-Driven Design

Setiap keputusan teknis harus mendukung **kecepatan rilis dan operasional event nyata**.
Fitur hanya boleh ditambahkan jika:

* Mendukung langsung alur pendaftaran peserta
* Mendukung operasional admin event
* Mengurangi risiko kegagalan di hari-H

Tidak ada fitur spekulatif, tidak ada optimasi prematur, dan tidak ada generalisasi berlebihan.

---

### II. Backend as the Source of Truth

Backend (Golang + PostgreSQL) adalah **otoritas tunggal** untuk:

* Status pendaftaran
* Status pembayaran
* Trigger pengiriman email

Frontend dan dashboard admin **tidak boleh** menyimpan state bisnis kritikal di sisi klien.

---

### III. State-Driven Workflow (NON-NEGOTIABLE)

Semua proses bisnis dikendalikan oleh **state eksplisit di database**.

State utama:

* `registration_status`: PENDING | CONFIRMED
* `payment_status`: UNPAID | PAID

Aturan:

* Email konfirmasi **hanya** dikirim ketika `payment_status` berubah ke `PAID`
* Tidak ada email manual di luar sistem
* Tidak ada implicit behavior

---

### IV. Simple Admin, Secure by Default

Admin dashboard:

* Single admin role
* Autentikasi sederhana (email + password)
* Password **wajib hashed (bcrypt)**
* Akses admin **terisolasi dari public site**

Keamanan dasar tidak boleh dikompromikan demi kecepatan.

---

### V. Design Consistency & Color Discipline

UI harus:

* Menggunakan **color palette resmi** proyek
* Menghindari gradient, glassmorphism, dan efek berlebihan
* Flat, poster-like, modern, dan readable

Desain bertujuan **mendukung event**, bukan mendominasi konten.

---

## Technical Constraints

* **Backend:** Golang (REST API)
* **Database:** PostgreSQL (relational, strongly typed)
* **Frontend:** Next JS AI (public page + admin dashboard)
* **Email:** SMTP (manual configuration)
* **Architecture:** Monolith MVP (NO microservices)
* **Deployment:** Single environment (production-first)

Larangan eksplisit:

* Tidak menggunakan NoSQL
* Tidak menggunakan message queue
* Tidak menggunakan external email SaaS API
* Tidak menggunakan frontend state management kompleks

---

## Development Workflow

1. **Define State First**

   * Tentukan field & enum di PostgreSQL
   * Tentukan kapan state berubah

2. **Design API Contract**

   * Endpoint jelas
   * Request/response eksplisit
   * Error handling konsisten

3. **Implement Backend**

   * Validasi input
   * Logging minimal tapi bermakna

4. **Integrate Frontend (Next JS React)**

   * Frontend mengikuti API, bukan sebaliknya
   * Tidak ada business logic di UI

5. **Manual Test Flow**

   * Register user
   * Login admin
   * Update payment
   * Verify email terkirim

Tidak ada deployment sebelum flow end-to-end berhasil.

---

## Governance

* Constitution ini **mengalahkan** keputusan teknis ad-hoc
* Setiap penambahan kompleksitas harus:

  * Dijustifikasi secara tertulis
  * Disetujui oleh owner proyek
* Jika melanggar prinsip MVP atau State-Driven Workflow → **harus direvisi**

Dokumen ini adalah referensi utama untuk:

* System design
* Review fitur
* Diskusi teknis

---

**Version**: 1.0.0
**Ratified**: 2026-01-01
**Last Amended**: 2026-01-01