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
