# Task Management API

REST API untuk manajemen task berbasis user. Aplikasi ini menyediakan fitur registrasi, login, autentikasi JWT, dan CRUD task milik masing-masing user.

## Teknologi

- Go
- PostgreSQL
- `database/sql` dengan driver `pgx`
- `httprouter` untuk routing HTTP
- `go-playground/validator` untuk validasi request
- `golang-jwt/jwt/v5` untuk JWT
- `bcrypt` untuk hashing password
- `godotenv` untuk membaca file `.env`
- `testify` untuk integration test

## Fitur API

Base URL default:

```text
http://localhost:8080/api/v1
```

Endpoint publik:

- `POST /register` - registrasi user
- `POST /login` - login user dan set cookie `access_token`

Endpoint task membutuhkan JWT melalui cookie `access_token` atau header:

```text
Authorization: Bearer <token>
```

Endpoint task:

- `POST /tasks` - membuat task
- `GET /tasks` - mengambil list task milik user login
- `GET /tasks/:taskId` - mengambil detail task
- `PUT /tasks/:taskId` - mengubah task
- `DELETE /tasks/:taskId` - menghapus task

Query list task:

- `status`: `pending`, `in_progress`, `completed`
- `title`: pencarian berdasarkan prefix title
- `page`: nomor halaman, default `1`
- `limit`: jumlah data, default `10`, maksimum `100`

## Environment

Buat file `.env` di root project:

```env
DATABASE_URL=postgres://postgres:password@localhost:5432/task_management_api?sslmode=disable
TEST_DATABASE_URL=postgres://postgres:password@localhost:5432/task_management_api_test?sslmode=disable
JWT_SECRET=change-this-secret
```

Keterangan:

- `DATABASE_URL` dipakai saat menjalankan aplikasi.
- `TEST_DATABASE_URL` dipakai integration test. Test akan skip jika env ini tidak tersedia.
- `JWT_SECRET` dipakai untuk sign dan validasi token.

## Database

Project ini mengasumsikan tabel `users` dan `tasks` sudah tersedia di PostgreSQL.

Minimal schema yang digunakan aplikasi:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    priority VARCHAR(50) NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);
```

Gunakan database berbeda untuk development dan test karena integration test menjalankan:

```sql
TRUNCATE TABLE tasks, users RESTART IDENTITY CASCADE;
```

## Menjalankan Aplikasi

Install dependency:

```powershell
go mod download
```

Jalankan server:

```powershell
go run .
```

Server berjalan di:

```text
http://localhost:8080
```

## Contoh Request

Register:

```json
{
  "name": "User Test",
  "email": "user@example.com",
  "password": "password123"
}
```

Login:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Create task:

```json
{
  "title": "Memasak Mie Goreng",
  "description": "langkah langkah benar dalam membuat mie goreng adalah sebagai berikut....",
  "priority": "high",
  "status": "pending",
  "due_date": "07-05-2026 12:30:00"
}
```

Format `due_date` utama:

```text
DD-MM-YYYY HH:mm:ss
```

Contoh:

```text
07-05-2026 12:30:00
```

Input `due_date` disimpan dan dikembalikan sebagai waktu UTC. Format RFC3339 juga diterima, misalnya:

```text
2026-05-07T12:30:00+07:00
```

## Menjalankan Test

Unit test dan integration test:

```powershell
go test ./...
```

Test khusus user controller:

```powershell
go test ./tests -run TestRegister -v
go test ./tests -run TestLogin -v
```

Test khusus task controller:

```powershell
go test ./tests -run TestCreateTask -v
go test ./tests -run TestGetAllTasks -v
go test ./tests -run TestUpdateTask -v
go test ./tests -run TestDeleteTask -v
```

Static check:

```powershell
go vet ./...
```

Catatan test:

- Pastikan `TEST_DATABASE_URL` mengarah ke database test.
- Jangan arahkan `TEST_DATABASE_URL` ke database development/production.
- Integration test akan melakukan truncate tabel `tasks` dan `users`.

## Struktur Project

```text
app/          konfigurasi database
controller/   HTTP handler
exception/    tipe error aplikasi
helper/       helper JSON, mapper response, transaksi
middleware/   validasi JWT
model/        domain model dan web DTO
repository/   akses database
service/      business logic
tests/        integration test
```

## Catatan Implementasi

- Password user disimpan sebagai bcrypt hash.
- Login mengirim JWT melalui cookie `access_token`.
- Task selalu dibatasi berdasarkan `user_id` dari JWT.
- Response `GET /tasks` hanya mengembalikan field ringkas sesuai contract: `id`, `title`, `status`, `priority`, `due_date`.
- Response detail task mengembalikan field lengkap termasuk `description`, `created_at`, dan `updated_at`.
