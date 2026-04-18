# CLAUDE.md

## Project Overview

Personal portfolio & blog website untuk **Syarif Hidayatullah** sebagai software engineer branding.
Monolith Go dengan server-side rendering — tidak ada frontend framework terpisah.

## Tech Stack

| Layer | Tech |
|---|---|
| Backend | Go 1.23 + Gin framework |
| ORM | GORM v2 |
| Database | PostgreSQL |
| Frontend | Alpine.js (CDN) + Tailwind CSS (CDN) |
| Templates | Go `html/template` |
| Markdown | goldmark |
| Auth | Session-based (gin-contrib/sessions + bcrypt) |
| Email | SMTP atau Resend API |

## Project Structure

```
cmd/server/main.go          # Entry point, router setup, DB init, admin seed
config/config.go            # Load env vars, build PostgreSQL DSN
internal/
  model/                    # GORM models (User, Post, Project, ContactMessage)
  repository/               # DB queries, satu file per entity
  service/                  # Business logic (post, project, message, email)
  handler/                  # Public HTTP handlers (home, blog, project, contact)
  handler/admin/            # Admin HTTP handlers (auth, dashboard, post, project, message)
  middleware/auth.go         # Session auth middleware
migrations/001_init.sql     # PostgreSQL DDL (referensi, AutoMigrate yang dipakai)
web/
  templates/
    partials/               # head.html, navbar.html, footer.html, admin_sidebar.html
    *.html                  # Public pages: home, blog_list, blog_detail, projects, error
    admin/*.html            # Admin pages: login, dashboard, posts, projects, messages
  static/
    css/input.css           # Tailwind source (belum di-build, sekarang pakai CDN)
```

## Database Config

PostgreSQL, menggunakan variabel terpisah (bukan `DATABASE_URL`):

```env
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=portfolio
DB_SSLMODE=disable
```

DSN dibangun otomatis di `config/config.go` via fungsi `buildDSN()`.

## Menjalankan

```bash
# 1. Buat database
createdb portfolio

# 2. Copy dan isi env
cp .env.example .env

# 3. Jalankan server
make dev

# Atau langsung
CGO_ENABLED=0 go run ./cmd/server/main.go

# Build Tailwind CSS (opsional, butuh Node.js)
npm install && npm run css:build
```

> Saat ini semua template menggunakan **Tailwind CDN**.
> Untuk production, build CSS lokal dan ganti CDN dengan `/static/css/app.css`.

## Routes

### Public
| Method | Path | Handler |
|---|---|---|
| GET | `/` | Home — featured projects + recent posts + contact form |
| GET | `/blog` | Daftar semua post yang published |
| GET | `/blog/:slug` | Detail post, render markdown → HTML |
| GET | `/projects` | Semua projects |
| POST | `/contact` | Submit contact form, returns JSON |

### Admin (semua butuh session auth)
| Method | Path | Fungsi |
|---|---|---|
| GET | `/admin` | Dashboard — stats + recent data |
| GET/POST | `/admin/posts` | List + Create post |
| GET/POST | `/admin/posts/new` | Form post baru |
| GET/POST | `/admin/posts/:id/edit` | Edit post |
| POST | `/admin/posts/:id/delete` | Hapus post |
| POST | `/admin/posts/:id/toggle-publish` | Publish/unpublish post |
| GET/POST | `/admin/projects` | List + Create project |
| GET/POST | `/admin/projects/new` | Form project baru |
| GET/POST | `/admin/projects/:id/edit` | Edit project |
| POST | `/admin/projects/:id/delete` | Hapus project |
| GET | `/admin/messages` | Inbox contact messages |
| POST | `/admin/messages/:id/read` | Mark as read |
| POST | `/admin/messages/:id/delete` | Hapus message |
| GET | `/admin/login` | Halaman login |
| POST | `/admin/login` | Submit login |
| POST | `/admin/logout` | Logout |

## Penting: Template System

- Template di-load dengan `html/template` + `ParseGlob`, bukan Gin default
- Setiap file template mendefinisikan nama-nya sendiri: `{{define "home.html"}}...{{end}}`
- Handler memanggil: `c.HTML(200, "home.html", gin.H{...})`
- Partials diinclude via: `{{template "navbar" .}}`
- Template function custom: `safeHTML` (untuk render HTML dari markdown) dan `joinStrings`
- **Jangan gunakan `slice` untuk membuat array di template** — Go template tidak support. Hardcode langsung atau pass dari handler.

## Model UUID

UUID di-generate di Go via `BeforeCreate` hook:

```go
func (m *Model) BeforeCreate(tx *gorm.DB) error {
    if m.ID == uuid.Nil {
        m.ID = uuid.New()
    }
    return nil
}
```

## Admin User

Di-seed otomatis saat pertama kali server jalan (jika tabel `users` kosong).
Credentials dari env:
```env
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=changeme123
```

## Email Notification

Contact form submission mengirim email ke `NOTIFY_EMAIL` (fallback ke `ADMIN_EMAIL`).
Driver dipilih via `EMAIL_DRIVER=smtp` atau `EMAIL_DRIVER=resend`.
Email dikirim secara goroutine (non-blocking) — gagal kirim email tidak gagalkan request.
