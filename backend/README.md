# Portfolio backend

Go server that runs the whole portfolio: it serves the built React site from
`../frontend/dist` at `/`, plus the JSON API, HTML admin and media uploads.
No database and no framework — the standard library's `net/http`, JSON files
on disk, and `bcrypt`.

## Quick start

```sh
# 1. Configure (optional; everything has a default)
cp .env.example .env            # set SESSION_SECRET so logins survive restarts

# 2. Build the site once (or skip this; / then redirects to /admin)
(cd ../frontend && npm install && npm run build)

# 3. Run
go run ./cmd/server
```

On first start the server:

1. Loads `data/content.json` (shipped with the repo). If it is missing, a
   fresh one with the schema defaults and empty lists is written.
2. Creates the admin user from `ADMIN_USERNAME` / `ADMIN_PASSWORD` when
   `data/admin.json` has none. If the password is empty it generates one and
   **prints it in the log** — copy it and change it under *Password* once
   logged in.

Then open:

| URL | What |
| --- | --- |
| `http://localhost:8080/` | The site (from `FRONTEND_DIR`); redirects to `/admin` if no build is found |
| `http://localhost:8080/admin` | Admin (log in with the credentials above) |
| `http://localhost:8080/api/v1/portfolio` | The JSON the site renders from |
| `http://localhost:8080/uploads/<file>` | Uploaded media |

## Data files

Both files live in `DATA_DIR` (default `./data`):

| File | Holds | In git? |
| --- | --- | --- |
| `content.json` | Settings and every list (education, skills, profiles, experience, projects, papers) | Yes — it *is* your content |
| `admin.json` | Admin accounts (bcrypt hashes) | No |

`content.json` is written by the admin and is safe to edit by hand, even
while the server is running: it notices the file changed and reloads it on
the next request, and a later admin save keeps your edit. Values
are strings as the admin stores them (`"1"`/`"0"` for checkboxes, one item
per line for lists, comma-separated tags), but plain JSON booleans and
numbers are accepted too and rows without an `id` get one on load. Every
save rewrites the file completely, pretty-printed in schema order, via a
temp file and rename so a crash cannot leave it half-written.

Back up by copying the directory. Reset to the shipped content with
`git checkout data/content.json`.

## Configuration

Environment variables, or a `.env` file in the working directory
(`.env.example` lists them all):

| Variable | Default | Meaning |
| --- | --- | --- |
| `ADDR` | `:8080` | Listen address |
| `DATA_DIR` | `./data` | Directory holding the three JSON files |
| `ADMIN_USERNAME` | `admin` | First admin, created only when none exists |
| `ADMIN_PASSWORD` | generated | First admin's password |
| `SESSION_SECRET` | random per start | ≥ 16 chars; set it so sessions survive restarts |
| `UPLOAD_DIR` | `./uploads` | Media storage, served at `/uploads/` |
| `FRONTEND_DIR` | `../frontend/dist` | Built React site served at `/`; skipped if it has no `index.html` |
| `CORS_ORIGINS` | — | Extra origins allowed to call the API; only needed if the site is hosted elsewhere |
| `SECURE_COOKIES` | `false` | Set `true` behind HTTPS |

## Admin

`/admin` is a server-rendered UI with a sidebar:

* **Dashboard** – counts and quick links.
* **Settings** – every text, label, image URL and section toggle, grouped
  (Identity, SEO, Sections, Navbar, Hero, About, Projects, Research, Footer).
  Media fields offer a dropdown of uploaded files.
* **Education / Skills / Profiles / Experience / Projects / Papers** – list,
  create, edit, delete, and ↑/↓ reorder. Forms validate per field and keep
  your input on error.
* **Media** – upload jpg/png/webp/gif/svg/pdf (≤ 15 MB); copy the URL into a
  setting (photo, CV) or entry.
* **Password** – change your password.

Security: bcrypt passwords, HMAC-signed `HttpOnly` session cookies (12 h),
CSRF tokens on every form, `SameSite=Lax`.

## Project layout

```
cmd/server/main.go         wiring and startup
data/content.json          the content (tracked); admin.json appears at runtime
internal/
  schema/                  content model: fields, resources, settings (edit here to add fields)
  store/                   JSON files in memory: generic CRUD, typed reads, users
  models/                  API response types
  api/                     /api/v1 handlers, rate limiter
  admin/                   /admin handlers, auth, renderer, web/templates, web/static
  site/                    serves the built React app with SPA fallback
  media/                   upload directory
  httpx/                   JSON helpers, middleware
  config/                  env + .env loading
```

See `../docs/ARCHITECTURE.md` for the design and "how to add a field" recipes,
and `../docs/API.md` for the endpoint reference.

## Development

```sh
go test ./...        # unit tests, including a check that data/content.json is valid
go vet ./...
go build -o server ./cmd/server
```

Adding a field is a one-line change in `internal/schema`; there are no
migrations. Existing `content.json` files simply gain the new key on the
next save (empty, or the setting's `Default`).

## Deploying

Build the site (`npm run build` in `../frontend`) and the binary
(`go build -o server ./cmd/server`), then run the binary with `FRONTEND_DIR`
pointing at the built `dist/`, `DATA_DIR` and `UPLOAD_DIR` on persistent
storage, `SESSION_SECRET` set, and `SECURE_COOKIES=true` behind an HTTPS
reverse proxy (Caddy, nginx). Keep it running with a systemd unit, for
example:

```ini
[Unit]
Description=Portfolio
After=network.target

[Service]
User=www-data
WorkingDirectory=/srv/portfolio
Environment=FRONTEND_DIR=/srv/portfolio/dist
Environment=DATA_DIR=/srv/portfolio/data
Environment=UPLOAD_DIR=/srv/portfolio/uploads
Environment=SESSION_SECRET=change-me-to-a-long-random-string
Environment=SECURE_COOKIES=true
ExecStart=/srv/portfolio/server
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
