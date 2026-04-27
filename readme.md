# Real Time Chat API + WebSocket (Go)

## Overview

REST API real-time chat yang dibangun dengan Go (Fiber) + WebSocket, mendukung message berbasis room, JWT Authentication, dan pub/sub Redis untuk broadcast.

---

## 🚀 Technology Stack

- Go (Golang)
- Fiber v2
- GORM – ORM
- PostgreSQL
- Docker & Docker Compose
- OpenAPI 3.0 (Swagger)
- Redis 7
- JWT
- Bcrypt
- Google UUID

---

## 📋 Features

- User registration & login with JWT authentication
- Create public & private chat rooms
- Real-time messaging via WebSocket
- Message history with pagination
- Redis pub/sub for multi-instance broadcasting
- Swagger API documentation

---

## 🏗️ Project Structure

```
Real-Time-Chat-API/
├── cmd/
│   └── main.go                  # Entry point
├── internal/
│   ├── config/                  # App & DB configuration
│   ├── handler/                 # HTTP & WebSocket handlers
│   ├── middleware/              # JWT middleware
│   ├── models/
│   │   ├── domain/              # Database models
│   │   └── web/                 # Request/Response structs
│   ├── repository/              # Database layer
│   ├── service/                 # Business logic
│   ├── ws/                      # WebSocket hub & client
│   ├── helper/                  # Response helpers
│   └── exception/               # Error handler
├── pkg/
│   └── pubsub/                  # Redis pub/sub client
├── docs/                        # Auto-generated Swagger docs
├── .env.example
├── docker-compose.yml
└── Dockerfile
```

---

## ⚙️ Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- Go 1.21+ (jika ingin run tanpa Docker)

### 1. Clone Repository

```bash
git clone https://github.com/DhahikaR/real-time-chat-api.git
cd real-time-chat-api
```

### 2. Setup Environment

```bash
cp .env.example .env
```

Isi file `.env` sesuai kebutuhan:

```env
APP_PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=real_time_chat
DB_SSLMODE=disable

REDIS_ADDR=redis:6379
JWT_SECRET=your_super_secret_key
```

### 3. Jalankan dengan Docker

```bash
docker-compose up -d
```

Aplikasi akan berjalan di `http://localhost:8080`

### 4. Akses Swagger Documentation

```
http://localhost:8080/swagger/index.html
```

---

## 📡 API Endpoints

### Authentication

| Method | Endpoint         | Description             | Auth |
| ------ | ---------------- | ----------------------- | ---- |
| POST   | `/auth/register` | Daftar user baru        | ❌   |
| POST   | `/auth/login`    | Login & dapat JWT token | ❌   |

### Rooms

| Method | Endpoint             | Description                   | Auth |
| ------ | -------------------- | ----------------------------- | ---- |
| POST   | `/rooms`             | Buat chat room baru           | ✅   |
| GET    | `/rooms`             | List semua chat room          | ✅   |
| GET    | `/rooms/:id/message` | Ambil pesan dengan pagination | ✅   |

### WebSocket

| Method | Endpoint                     | Description                          | Auth |
| ------ | ---------------------------- | ------------------------------------ | ---- |
| WS     | `/ws/connect?room_id={uuid}` | Connect ke room untuk real-time chat | ✅   |

---

## 🔐 Authentication

API ini menggunakan **JWT Bearer Token**. Setelah login, tambahkan token ke setiap request:

```
Authorization: Bearer <token>
```

---

## 📨 Request & Response Examples

### Register

```bash
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

### Login

```bash
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Create Room

```bash
POST /rooms
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "General Chat",
  "type": "public"
}
```

### Get Messages (dengan Pagination)

```bash
GET /rooms/{id}/message?page=1&limit=20
Authorization: Bearer <token>
```

### Connect WebSocket

```
# Gunakan ws:// protocol, bukan http://
ws://localhost:8080/ws/connect?room_id={uuid}

# Header
Authorization: Bearer <token>

# Format pesan yang dikirim (JSON)
{
  "room_id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Hello!"
}
```

---

## 🔄 Alur Penggunaan

```
1. POST /auth/register    → daftar akun
2. POST /auth/login       → login, simpan token
3. POST /rooms            → buat room, simpan room_id
4. WS  /ws/connect        → connect WebSocket dengan room_id
5. Send message           → kirim pesan JSON via WebSocket
6. GET /rooms/:id/message → ambil riwayat pesan
```

---

## 🧪 Testing

```bash
go test ./... -v
```

Testing mencakup:

- Service (auth)
- Repository (user)
- Config (database)

---

## Author

**Dhahika Rahmadani**  
Backend Developer • Go Enthusiast  
📧 [dhahikardani@gmail.com](mailto:dhahikardani@gmail.com)
