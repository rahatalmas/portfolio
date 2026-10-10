# Architecture

A single-page React portfolio whose content is managed in a Go admin and
stored in MySQL. The site and the backend are two deployments that talk over
the API.

```
  static host (Vercel, nginx, …)          Go server (backend/cmd/server)
  ┌───────────────────────────┐           ┌──────────────────────────────────────────┐
  │  React app (frontend/)    │  GET      │  /api/v1/*   JSON API ─────┐             │
  │  built with VITE_API_BASE ├──────────▶│  /admin/*    HTML admin ───┼─▶ store ─▶ MySQL
  │  = backend origin         │  CORS     │  /uploads/*  media on disk │             │
  └───────────────────────────┘           │  /           → /admin      │             │
                                          └────────────────────────────┴─────────────┘
```

* The React app (`frontend/`) makes **one** request, `GET /api/v1/portfolio`,
  and renders from it. It has no content of its own. In development the Vite
  proxy makes that request same-origin; in production the build is given
  `VITE_API_BASE` and the backend lists the site's origin in `CORS_ORIGINS`.
* The admin (`/admin`) is server-rendered Go templates. Editing there changes
  MySQL; the site reflects it on the next load (API responses are cached for
  30 s by `Cache-Control`).
* Uploads (photo, CV, images) are files in `UPLOAD_DIR`, served at `/uploads/`.
  The API returns them as backend-relative paths; the client prefixes them
  with `VITE_API_BASE` (`absolutizeUploads` in `client.ts`).

## Repository layout

```
portfolio/
├── frontend/                 React + TypeScript app (Vite, Tailwind v4)
│   ├── src/
│   │   ├── api/              types.ts (API shapes), client.ts (fetch), PortfolioContext.tsx (loader)
│   │   ├── components/       One file per UI piece; all read from usePortfolio()
│   │   ├── hooks/            useTheme, useScrollSpy, useReveal
│   │   └── pages/Home.tsx    Section order and visibility
│   ├── public/               Static files (placeholder photo)
│   ├── index.html, vite.config.ts (dev proxy → :8080), tsconfig.json
│   └── dist/                 Build output — deploy this to a static host
├── backend/
│   ├── Dockerfile            Backend image (API + admin + uploads)
│   ├── cmd/server/main.go    Wiring: config → db → migrate → seed → routes → serve
│   └── internal/
│       ├── schema/           THE content model (fields, resources, settings)
│       ├── db/               MySQL pool + embedded migrations
│       ├── store/            Data access (typed reads, generic CRUD, messages, users)
│       ├── models/           JSON shapes returned by the API
│       ├── api/              /api/v1 handlers + contact-form rate limit
│       ├── admin/            /admin handlers, sessions, templates (web/)
│       ├── media/            Upload directory
│       ├── httpx/            JSON helpers + middleware
│       └── config/           Environment / .env
├── docs/                     This file, API.md
├── docker-compose.yml        MySQL + backend
└── PLAN.md                   Original design notes
```

## The content model lives in one place

`backend/internal/schema` is the heart of the backend:

* **`fields.go`** – `Field` (key, label, kind, options, help…) and `Kind`
  (`text`, `textarea`, `lines`, `bool`, `number`, `select`, `tags`, `url`,
  `email`, `slug`, `media`). `Field.Clean` validates/normalises a form value;
  `Field.Typed` converts a stored string to its JSON type.
* **`resources.go`** – `Resources`: every list table (education, skills,
  profiles, experience, projects, papers) with its editable columns.
* **`settings.go`** – `SettingsGroups`: every single-value setting with its
  default, grouped as the admin shows them.

Everything else is generic over these specs:

| Consumer | What it derives from the spec |
| --- | --- |
| `store/records.go` | `SELECT`/`INSERT`/`UPDATE` column lists, reordering |
| `store/settings.go` | Seeding defaults, upserts |
| `admin/handlers.go` | List columns, form fields, validation, defaults |
| `admin/web/templates/partials/field.html` | Input widget per kind |
| `api` (`GET /site`) | JSON typing of settings |
| `store/seed.go` | Demo content validated against the spec |

### Storage conventions

* `settings(key, value)` – one row per setting, values stored as strings and
  typed on the way out (`"1"`/`"0"` for bools, newline-separated for `lines`,
  comma-separated for `tags`).
* Each resource table has `id`, its columns, and `sort_order`. The admin's
  ↑/↓ renumbers `sort_order` 1..n inside a transaction.
* `messages` holds contact-form submissions; `admin_users` holds bcrypt hashes.
* `schema_migrations` records applied migration versions.

## Request flow

**Site visit**
Static host serves `index.html` → React boots → `GET {VITE_API_BASE}/api/v1/portfolio`
→ CORS check → `store.Portfolio` runs seven small queries → JSON → upload
paths made absolute → components render.

**Admin edit**
`GET /admin/projects/3` → `auth` middleware verifies the signed session cookie
→ `GetRecord` → `form.html` renders inputs from `Resource.Fields` →
`POST /admin/projects/3` → CSRF check → `Field.Clean` on every column → on
errors re-render with messages, else `UpdateRecord` → redirect with a flash.

**Contact form**
`POST /api/v1/messages` → per-IP rate limit (5/min) → JSON validation →
honeypot check → `CreateMessage` → appears in `/admin/messages` with an
unread badge.

## Security notes

* Sessions are stateless HMAC-signed cookies (`admin_session`, `HttpOnly`,
  `SameSite=Lax`, `Path=/admin`, 12 h). Set `SESSION_SECRET` so restarts don't
  log you out, and `SECURE_COOKIES=true` behind HTTPS.
* Every admin POST needs a `_csrf` token derived from the session cookie.
* Login uses bcrypt and runs the hash comparison even for unknown users.
* All SQL uses placeholders; identifiers come from the spec, never from input.
* Uploads are restricted by extension, renamed to safe lowercase names, capped
  at 15 MB, and served without directory listing.
* Admin pages send `Cache-Control: no-store` and `X-Frame-Options: DENY`.

## How to…

### Add a text or flag to the site (a setting)

1. Append a `Field` to the right group in `schema/settings.go` with a
   `Default`. It is seeded on next start and editable in the admin.
2. Add the key to `Site` in `frontend/src/api/types.ts`.
3. Use `site.<key>` in a component.

### Add a column to a resource (e.g. a `year` on projects)

1. New migration `backend/internal/db/migrations/0002_project_year.sql`:
   `ALTER TABLE projects ADD COLUMN year INT NOT NULL DEFAULT 0;`
2. Add `{Key: "year", Label: "Year", Kind: Number}` to the projects
   `Resource` in `schema/resources.go`. The admin list/form/validation update
   themselves.
3. Add the column to the typed read in `store/lists.go` and the struct in
   `models/models.go`; mirror it in `frontend/src/api/types.ts`; render it.

### Add a whole new resource (e.g. talks)

1. Migration creating the table with `id`, your columns, `sort_order`.
2. A new `Resource` in `schema/resources.go` — the admin CRUD is now live.
3. A typed read in `store/lists.go`, a model, a line in `store.Portfolio`,
   an API route if you want it standalone, a TypeScript type, a component.
4. Optionally a demo row in `store/seed.go`.

### Add an admin page that isn't a resource

Follow `messages` in `admin/handlers.go`: a handler, a template in
`web/templates/pages/`, routes in `Routes`, and a sidebar link in
`layout.html`. Add a case to `TestAllPagesRender`.

## Testing

```sh
cd backend  && go test ./...    # schema validation, migrations split, media, sessions, templates, routes, API validation
cd frontend && npx tsc -b       # frontend types
```

Unit tests don't need MySQL. For an end-to-end check, run MySQL (see
`docker-compose.yml`), start the server, and open `/admin` and
`/api/v1/portfolio`.

## Frontend conventions

* `usePortfolio()` is the only data source; components never fetch.
* Theme tokens live in `frontend/src/index.css` (`@theme` for light, `.dark` for dark).
* Sections share `Section.tsx` (max-width 960 px, consistent spacing).
* No animation beyond the first-reveal fade and the nav underline.
