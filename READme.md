# MONALISA – RBAC Management System

MONALISA adalah sistem manajemen pengguna berbasis **Role-Based Access Control (RBAC)** yang dibangun dengan **Go (Gin)** sebagai backend dan **React (Vite)** sebagai frontend.  
Sistem ini dirancang untuk kebutuhan instansi (Balai Pemasyarakatan) dengan fokus pada **keamanan, auditabilitas, dan kontrol akses yang ketat**.

---

## 🎯 Fitur Utama

- ✅ Login berbasis **NIP**
- ✅ **JWT Authentication**
- ✅ **Role-Based Access Control (RBAC)**
- ✅ Assign & Remove Role (Admin)
- ✅ Proteksi **admin-self** (tidak bisa menghapus role admin dari dirinya sendiri)
- ✅ **Audit Log** (siapa melakukan apa, kapan)
- ✅ Frontend RBAC UI (assign/remove role)
- ✅ Backend siap audit & production-ready (struktur bersih)

---

## 🧱 Arsitektur

```

Frontend (React)
↓ JWT
Backend API (Gin)
↓
RBAC Middleware
↓
Service Layer
↓
Repository Layer
↓
PostgreSQL

```

### Prinsip Arsitektur
- Handler → Service → Repository (satu arah)
- Tidak ada import cycle
- Interface dikunci di service layer
- Audit log **WAJIB** di backend (tidak tergantung frontend)

---

## 📁 Struktur Folder

### Backend (`monalisa-be`)
```

monalisa-be/
├── main.go
├── go.mod
├── .env
└── internal/
├── model/
├── repository/
├── service/
├── handler/
└── middleware/

```

### Frontend (`monalisa-fe`)
```

monalisa-fe/
├── src/
│   ├── pages/
│   │   └── AdminUsers.jsx
│   ├── components/
│   │   └── RoleBadge.jsx
│   ├── services/
│   │   └── api.js
│   └── main.jsx
└── .env

````

---

## 🔐 Environment Variable

### Backend `.env`
```env
DATABASE_URL=postgres://postgres:password@localhost:5432/monalisa?sslmode=disable
JWT_SECRET=monalisa_dev_secret_123
````

> ⚠️ `JWT_SECRET` adalah **kunci server**, bukan token user.
> Jangan pernah dikirim ke frontend atau di-commit ke repository publik.

### Frontend `.env`

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

---

## 🗄️ Database Schema (Wajib)

### Tabel `audit_logs`

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID NOT NULL,
    action TEXT NOT NULL,
    target TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
```

### Tabel inti lain (ringkas)

* `users (id, nip, is_active)`
* `employees (nip, nama, jabatan)`
* `roles (id, code)`
* `permissions (id, code)`
* `user_roles (user_id, role_id)`
* `role_permissions (role_id, permission_id)`

---

## ▶️ Cara Menjalankan Backend

```bash
cd monalisa-be
go mod tidy
go run .
```

Server akan berjalan di:

```
http://localhost:8080
```

---

## ▶️ Cara Menjalankan Frontend

```bash
cd monalisa-fe
npm install
npm run dev
```

Frontend berjalan di:

```
http://localhost:5173
```

---

## 🔑 Authentication Flow

1. User login dengan NIP:

```http
POST /api/v1/auth/login
```

Request:

```json
{
  "nip": "196807241991032001"
}
```

Response:

```json
{
  "token": "JWT_TOKEN"
}
```

2. Token disimpan di `localStorage`
3. Semua request admin menggunakan:

```
Authorization: Bearer <token>
```

---

## 🛡️ RBAC Rules

* Semua endpoint `/api/v1/admin/*`:

  * Wajib JWT valid
  * Wajib permission `user.manage`

* **Hard Rule (Backend)**:

  * Admin **tidak bisa menghapus role `admin` dari dirinya sendiri**
  * Rule ini **tidak bisa dibypass oleh frontend**

---

## 👥 Admin User Management Endpoint

| Method | Endpoint                       | Keterangan        |
| ------ | ------------------------------ | ----------------- |
| GET    | `/admin/users`                 | List user + roles |
| POST   | `/admin/users/:id/roles`       | Assign role       |
| DELETE | `/admin/users/:id/roles/:role` | Remove role       |
| GET    | `/admin/roles`                 | List role         |

---

## 🧪 Contoh Test (PowerShell)

```powershell
$res = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/auth/login" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{ "nip": "196807241991032001" }'

$token = $res.token

Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/admin/users" `
  -Headers @{ Authorization = "Bearer $token" }
```

---

## 🧾 Audit Log

Setiap:

* assign role
* remove role

akan menghasilkan record di `audit_logs`:

```sql
SELECT * FROM audit_logs ORDER BY created_at DESC;
```

Audit log **tidak bergantung frontend** dan **tidak bisa dimatikan**.

---

## ⚠️ Catatan Penting

* Jangan menghapus audit log di production
* Jangan hardcode `JWT_SECRET`
* Jangan mengandalkan frontend untuk security
* Semua rule kritis **HARUS DI BACKEND**

---

## 🚀 Status Project

* ✅ Backend compile bersih
* ✅ JWT valid
* ✅ RBAC aktif
* ✅ Admin self-protection aktif
* ✅ Audit log konsisten
* ✅ Siap dikembangkan lebih lanjut

---

## 📌 Next Development (Opsional)

* Audit Log Viewer (UI)
* Permission-based menu (frontend)
* Refresh token & rate limiting
* Deployment (Docker / CI)
