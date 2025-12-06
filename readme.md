# 🧱 Go-Arch — Clean Architecture Boilerplate (Gin + PostgreSQL)

Project ini dibuat dengan pendekatan **Clean Architecture**, memisahkan tiap lapisan aplikasi agar lebih modular, scalable, dan mudah di-maintain.  
Framework utama yang digunakan adalah **Gin Gonic** untuk web server dan **SQLX** untuk koneksi database PostgreSQL.

---

## 🏗️ Project Structure

/cmd
└── main.go # Entry point aplikasi

/internal
├── entity/ # Domain entities (User, Product, dsb.)
├── usecase/ # Business logic (Auth, Register, Login)
├── handler/ # Controller / HTTP handler (Gin)
├── repository/ # Interface repository (contract DB)
├── infrastructure/
│ └── pgsql/ # Implementasi repository untuk PostgreSQL
├── errors/ # Custom error wrapper
└── presenter/ # (Opsional) response presenter / view formatter

/migrations # Folder migrasi database (SQL files)

---

## 🧩 Clean Architecture Layers

| Layer                       | Role                                                                            |
| --------------------------- | ------------------------------------------------------------------------------- |
| **Entities**                | Berisi struktur data inti (misal: `User`, `Product`, `Order`)                   |
| **Use Cases**               | Berisi logika bisnis utama (register, login, update user, dsb.)                 |
| **Controllers / Handlers**  | Menangani HTTP request dari user (via Gin)                                      |
| **Gateways / Repositories** | Abstraksi untuk data access (database, API, file, dsb.)                         |
| **Presenters**              | Mengatur format response API (status, message, data)                            |
| **External Interfaces**     | Mengatur koneksi ke sistem luar (DB, Redis, API lain)                           |
| **UI / Web**                | Layer antarmuka web (Gin HTTP routes)                                           |
| **DB / Devices**            | Layer paling bawah, berhubungan langsung dengan driver DB atau sistem eksternal |

---

## ⚙️ Dependency Installation

Install dependencies berikut:

```bash
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/joho/godotenv
go get github.com/jmoiron/sqlx
```

---

## 🔐 Environment Variables

Buat file `.env` di root project dengan konfigurasi berikut (lihat `.env.example`):

```env
# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_PORT=5432
DB_NAME=go_arch

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
TOKEN_EXPIRY=24h  # Format: 24h, 1h, 30m, etc.
```

**Catatan:**

- `JWT_SECRET` **WAJIB** diisi dengan secret key yang kuat (minimal 32 karakter)
- `TOKEN_EXPIRY` menggunakan format Go duration (contoh: `24h`, `1h`, `30m`, `720h` untuk 30 hari)
- Jika `TOKEN_EXPIRY` tidak di-set, default akan menggunakan `24h`

---

## 🔒 Security Best Practices

### ⚠️ File Sensitif - JANGAN COMMIT!

File-file berikut **TIDAK BOLEH** di-commit ke repository:

- ✅ `.env` - Environment variables (sudah di-ignore)
- ✅ File dengan kata `secret`, `key`, `credential`, `password`
- ✅ File certificate/keys (`.pem`, `.key`, `.crt`, `.cert`)
- ✅ Database dumps dan backups (`.sql`, `.dump`, `.backup`)
- ✅ Log files yang mungkin mengandung informasi sensitif

### ✅ Yang BOLEH di-commit:

- ✅ `.env.example` - Template environment variables (tanpa nilai sensitif)
- ✅ Migration files di folder `migrations/`
- ✅ Source code dan dokumentasi

### 🛡️ Checklist Sebelum Commit:

1. ✅ Pastikan file `.env` tidak ter-track: `git status | grep .env`
2. ✅ Pastikan tidak ada credential hardcoded di source code
3. ✅ Pastikan JWT_SECRET menggunakan environment variable
4. ✅ Review perubahan dengan: `git diff` sebelum commit

### 🔧 Jika File Sensitif Terlanjur Ter-commit:

Jika file sensitif sudah ter-commit, segera:

```bash
# Hapus dari git tracking (tapi tetap di local)
git rm --cached .env

# Commit perubahan
git commit -m "Remove sensitive .env file from tracking"

# Jika sudah ter-push, ubah semua secrets yang ter-expose!
```

**⚠️ PENTING:** Jika file sensitif sudah ter-push ke GitHub, anggap semua credentials sudah compromised dan **WAJIB diganti**!
