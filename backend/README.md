# Portfolio backend

Go + MySQL server that powers the portfolio: JSON API, HTML admin and media
uploads. The React site in `../frontend` is hosted separately and reads from
this API. No framework — the standard library's `net/http`, `database/sql`
with the MySQL driver, and `bcrypt`.

## Quick start

```sh
# 1. A MySQL 8 / MariaDB 10.6+ database
mysql -u root -p -e "CREATE DATABASE portfolio CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
                     CREATE USER 'portfolio'@'%' IDENTIFIED BY 'secret';
                     GRANT ALL ON portfolio.* TO 'portfolio'@'%';"

# 2. Configure
cp .env.example .env            # edit DB_* and set SESSION_SECRET

# 3. Run
go run ./cmd/server
```

On first start the server:

1. Applies the migrations in `internal/db/migrations/`.
2. Seeds every setting's default value.
3. Seeds demo content if all content tables are empty (`SEED_DEMO=true`).
4. Creates the admin user from `ADMIN_USERNAME` / `ADMIN_PASSWORD`. If the
   password is empty it generates one and **prints it in the log** — copy it
   and change it under *Password* once logged in.

Then open:

| URL | What |
| --- | --- |
| `http://localhost:8080/admin` | Admin (log in with the credentials above); `/` redirects here |
| `http://localhost:8080/api/v1/portfolio` | The JSON the site renders from |
| `http://localhost:8080/uploads/<file>` | Uploaded media |

## Configuration

Environment variables, or a `.env` file in the working directory
(`.env.example` lists them all):

| Variable | Default | Meaning |
| --- | --- | --- |
| `ADDR` | `:8080` | Listen address |
| `DB_DSN` | — | Full MySQL DSN; overrides the `DB_*` parts |
| `DB_HOST` / `DB_PORT` | `127.0.0.1` / `3306` | |
| `DB_USER` / `DB_PASSWORD` / `DB_NAME` | `portfolio` / — / `portfolio` | |
| `ADMIN_USERNAME` | `admin` | First admin, created only when none exists |
| `ADMIN_PASSWORD` | generated | First admin's password |
| `SESSION_SECRET` | random per start | ≥ 16 chars; set it so sessions survive restarts |
| `UPLOAD_DIR` | `./uploads` | Media storage, served at `/uploads/` |
| `CORS_ORIGINS` | — | Comma-separated origins allowed to call the API from a browser — set to the deployed site's origin |
| `SECURE_COOKIES` | `false` | Set `true` behind HTTPS |
| `SEED_DEMO` | `true` | Insert demo content into empty content tables |

## Admin

`/admin` is a server-rendered UI with a sidebar:

* **Dashboard** – counts, recent messages, quick links.
* **Settings** – every text, label, image URL and section toggle, grouped
  (Identity, SEO, Sections, Navbar, Hero, About, Work, Research, Contact,
  Footer). Media fields offer a dropdown of uploaded files.
* **Education / Skills / Profiles / Experience / Projects / Papers** – list,
  create, edit, delete, and ↑/↓ reorder. Forms validate per field and keep
  your input on error.
* **Messages** – contact-form inbox with unread badge; open to mark read,
  reply by email, delete.
* **Media** – upload jpg/png/webp/gif/svg/pdf (≤ 15 MB); copy the URL into a
  setting (photo, CV) or entry.
* **Password** – change your password.

Security: bcrypt passwords, HMAC-signed `HttpOnly` session cookies (12 h),
CSRF tokens on every form, `SameSite=Lax`.

## Project layout

```
cmd/server/main.go         wiring and startup
internal/
  schema/                  content model: fields, resources, settings (edit here to add fields)
  db/                      connection + embedded migrations/NNNN_name.sql
  store/                   SQL: typed reads, generic CRUD, messages, users, demo seed
  models/                  API response types
  api/                     /api/v1 handlers, rate limiter
  admin/                   /admin handlers, auth, renderer, web/templates, web/static
  media/                   upload directory
  httpx/                   JSON helpers, middleware
  config/                  env + .env loading
```

See `../docs/ARCHITECTURE.md` for the design and "how to add a field" recipes,
and `../docs/API.md` for the endpoint reference.

## Development

```sh
go test ./...        # unit tests (no database needed)
go vet ./...
go build -o server ./cmd/server
```

Migrations: add `internal/db/migrations/0002_description.sql` with plain SQL
statements ending in `;`. They are embedded in the binary and applied in
order on start; applied versions are tracked in `schema_migrations`.

## Deploying

`go build -o server ./cmd/server`, then run the binary with `UPLOAD_DIR` on
persistent storage, `SESSION_SECRET` set, `SECURE_COOKIES=true` behind an
HTTPS reverse proxy (Caddy, nginx), and `CORS_ORIGINS` set to the site's
origin. `Dockerfile` here and `docker-compose.yml` in the repo root do this.

The frontend is deployed on its own: build it with `VITE_API_BASE` pointing at
this server (see `../frontend/README.md`).
