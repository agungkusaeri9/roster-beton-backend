# 📦 Database Migrations

Folder ini berisi file-file migrasi database untuk PostgreSQL menggunakan [golang-migrate](https://github.com/golang-migrate/migrate).

## 🚀 Setup

### 1. Install golang-migrate

**Linux/Mac:**

```bash
# Install via Homebrew (Mac)
brew install golang-migrate

# Atau download binary dari https://github.com/golang-migrate/migrate/releases
```

**Atau via Go:**

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 2. Setup Environment Variables

Pastikan file `.env` sudah ada di root project dengan konfigurasi:

```env
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=go_arch
```

### 3. Buat Database (jika belum ada)

```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE go_arch;

# Keluar
\q
```

## 📝 Cara Menggunakan

### Menjalankan Migrasi (Up)

```bash
cd migrations
chmod +x migration.sh
./migration.sh
```

### Rollback Migrasi (Down)

```bash
cd migrations
chmod +x migrate-down.sh
./migrate-down.sh
```

### Membuat Migrasi Baru

```bash
cd migrations
chmod +x create-migration.sh
./create-migration.sh <nama_migrasi>
# Contoh: ./create-migration.sh add_email_to_users
```

### Force Migration (jika ada masalah)

```bash
cd migrations
chmod +x migrate-force.sh
./migrate-force.sh <version>
# Contoh: ./migrate-force.sh 1
```

## 📁 Struktur File

- `*.up.sql` - File untuk menjalankan migrasi (create table, add column, dll)
- `*.down.sql` - File untuk rollback migrasi (drop table, remove column, dll)
- `migration.sh` - Script untuk menjalankan migrasi up
- `migrate-down.sh` - Script untuk rollback migrasi
- `migrate-force.sh` - Script untuk force migration ke version tertentu
- `create-migration.sh` - Script untuk membuat file migrasi baru

## ⚠️ Catatan Penting

1. **Jangan edit file migrasi yang sudah dijalankan** - Buat migrasi baru untuk perubahan
2. **Selalu test migrasi down** - Pastikan rollback berfungsi dengan baik
3. **Backup database** - Sebelum menjalankan migrasi di production, backup dulu
4. **Gunakan transaksi** - Pastikan migrasi bisa di-rollback jika gagal

## 🔍 Check Migration Status

Untuk melihat status migrasi, gunakan command:

```bash
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" version
```
