# Architecture

A single-page React portfolio whose content is managed in a Go admin and
stored as JSON files on disk. One Go process serves everything from a single
origin. There is no database.

```
  browser                     Go server (backend/cmd/server)
  ┌───────────┐               ┌──────────────────────────────────────────────────────┐
  │ React app │  GET /        │  /            built site (FRONTEND_DIR, SPA)          │
  │           ├──────────────▶│  /api/v1/*    JSON API ─────┐                         │
  │           │  GET /api/…   │  /admin/*     HTML admin ───┼─▶ store ─▶ DATA_DIR/     │
  └───────────┘  same origin  │  /uploads/*   media on disk │      content.json        │
                              └─────────────────────────────┘      messages.json       │
                                                                   admin.json          │
```

* The React app (`frontend/`) is built to `dist/` and served by the Go server
  at `/` (`internal/site`: real files as-is, hashed `/assets/*` cached for a
  year, everything else falls back to `index.html`). It makes **one** request,
  `GET /api/v1/portfolio`, same-origin, and renders from it. It first paints
  content bundled at build time from `backend/data/content.json`
  (`src/api/fallback.ts`) and swaps in the API response when it arrives; if
  the API fails, the bundled content stays and no error is shown. In
  development `npm run dev` serves it from Vite with a proxy to `:8080`
  instead. If no build is present the server still runs and `/` redirects to
  `/admin`.
* The admin (`/admin`) is server-rendered Go templates. Editing there rewrites
  `content.json`; the site reflects it on the next load (API responses are
  cached for 30 s by `Cache-Control`).
* Uploads (photo, CV, images) are files in `UPLOAD_DIR`, served at `/uploads/`.
  The API returns them as server-relative paths, which work as-is same-origin.

## Repository layout

```
portfolio/
├── frontend/                 React + TypeScript app (Vite, Tailwind v4)
│   ├── src/
│   │   ├── api/              types.ts (API shapes), client.ts (fetch), fallback.ts (bundled content), PortfolioContext.tsx (loader)
│   │   ├── components/       One file per UI piece; all read from usePortfolio()
│   │   ├── hooks/            useTheme, useScrollSpy, useReveal
│   │   └── pages/Home.tsx    Section order and visibility
│   ├── public/               Static files (placeholder photo)
│   ├── index.html, vite.config.ts (dev proxy → :8080), tsconfig.json
│   └── dist/                 Build output — served by the Go server at /
├── Dockerfile                One image: builds frontend + backend
├── backend/
│   ├── cmd/server/main.go    Wiring: config → store.Open → routes → serve
│   ├── data/content.json     THE content, tracked in git (messages.json, admin.json appear at runtime)
│   └── internal/
│       ├── schema/           THE content model (fields, resources, settings)
│       ├── store/            JSON files held in memory; CRUD, typed reads, messages, users
│       ├── models/           JSON shapes returned by the API
│       ├── api/              /api/v1 handlers + contact-form rate limit
│       ├── admin/            /admin handlers, sessions, templates (web/)
│       ├── site/             Serves the built React app with SPA fallback
│       ├── media/            Upload directory
│       ├── httpx/            JSON helpers + middleware
│       └── config/           Environment / .env loading
├── docs/                     This file + API.md
├── docker-compose.yml        The server with data + uploads volumes
└── PLAN.md                   Original design notes
```

## The content model lives in one place

`backend/internal/schema` is the heart of the backend:

* **`fields.go`** – `Field` (key, label, kind, options, help…) and `Kind`
  (`text`, `textarea`, `lines`, `bool`, `number`, `select`, `tags`, `url`,
  `email`, `slug`, `media`). `Field.Clean` validates/normalises a form value;
  `Field.Typed` converts a stored string to its JSON type.
* **`resources.go`** – `Resources`: every list type (education, skills,
  profiles, experience, projects, papers) with its editable fields.
* **`settings.go`** – `SettingsGroups`: every single-value setting with its
  default, grouped as the admin shows them.

Everything else is generic over these specs:

| Consumer | What it derives from the spec |
| --- | --- |
| `store/store.go` | Which keys to read/write in `content.json`, field order, slug uniqueness |
| `admin/handlers.go` | List columns, form fields, validation, defaults |
| `admin/web/templates/partials/field.html` | Input widget per kind |
| `api` (`GET /site`) | JSON typing of settings |
| `frontend/src/api/fallback.ts` | Conversion of the bundled `content.json` (mirrors `Field.Typed`) |

### Storage conventions

`content.json` is one object:

```json
{
  "settings": { "name": "…", "show_work": "1", "interests": "a\nb", … },
  "education": [ { "id": 1, "degree": "…", … } ],
  "skills":    [ … ], "profiles": [ … ], "experience": [ … ],
  "projects":  [ … ], "papers": [ … ]
}
```

* Settings are strings, typed on the way out (`"1"`/`"0"` for bools,
  newline-separated for `lines`, comma-separated for `tags`). Keys missing
  from the file use the schema `Default`.
* Each list entry has an `id` and its fields; array order is display order.
  The admin's ↑/↓ swaps entries and rewrites the file.
* Loading is forgiving: JSON booleans and numbers are accepted and stored as
  strings, rows without an `id` (or with a duplicate) get a fresh one, and
  unknown fields are dropped on the next save. Saving is strict: the file is
  regenerated in schema order and replaced atomically (temp file + rename).
* `messages.json` is an array of submissions; `admin.json` an array of
  accounts with bcrypt hashes. Both are written the same way.
* The store keeps all three files in memory behind one `sync.RWMutex`;
  reads never touch the disk.

## Request flow

**Site visit**
`GET /` → `site.Handler` serves `index.html` → React boots and renders the
bundled content → `GET /api/v1/portfolio` (same origin) → `store.Portfolio`
assembles the typed shape from memory → JSON → components re-render.

**Admin edit**
`GET /admin/projects/3` → `auth` middleware verifies the signed session cookie
→ `GetRecord` → `form.html` renders inputs from `Resource.Fields` →
`POST /admin/projects/3` → CSRF check → `Field.Clean` on every field → on
errors re-render with messages, else `UpdateRecord` rewrites `content.json`
→ redirect with a flash.

**Contact form**
`POST /api/v1/messages` → per-IP rate limit (5/min) → JSON validation →
honeypot check → `CreateMessage` appends to `messages.json` → appears in
`/admin/messages` with an unread badge.

## Security notes

* Sessions are stateless HMAC-signed cookies (`admin_session`, `HttpOnly`,
  `SameSite=Lax`, `Path=/admin`, 12 h). Set `SESSION_SECRET` so restarts don't
  log you out, and `SECURE_COOKIES=true` behind HTTPS.
* Every admin POST needs a `_csrf` token derived from the session cookie.
* Login uses bcrypt and runs the hash comparison even for unknown users.
* Data files are written only by the store, with names fixed in code; no
  user input ever forms a path.
* Uploads are restricted by extension, renamed to safe lowercase names, capped
  at 15 MB, and served without directory listing.
* Admin pages send `Cache-Control: no-store` and `X-Frame-Options: DENY`.

## How to…

### Add a text or flag to the site (a setting)

1. Append a `Field` to the right group in `schema/settings.go` with a
   `Default`. It is editable in the admin at once; `content.json` gains it on
   the next save.
2. Add the key to `Site` in `frontend/src/api/types.ts` (and to `BOOL_KEYS`
   or `LINES_KEYS` in `fallback.ts` if it is a bool or a list).
3. Use `site.<key>` in a component.

### Add a field to a resource (e.g. a `year` on projects)

1. Add `{Key: "year", Label: "Year", Kind: Number}` to the projects
   `Resource` in `schema/resources.go`. The admin list/form/validation and
   the stored file update themselves (existing entries read as empty).
2. Add the field to the typed read in `store/lists.go` and the struct in
   `models/models.go`; mirror it in `frontend/src/api/types.ts` and
   `fallback.ts`; render it.

### Add a whole new resource (e.g. talks)

1. A new `Resource` in `schema/resources.go` — the admin CRUD is now live and
   the key appears in `content.json` on the next save.
2. A typed read in `store/lists.go`, a model, a line in `store.Portfolio`,
   an API route if you want it standalone, a TypeScript type, a converter in
   `fallback.ts`, a component.
3. Optionally starter entries in `backend/data/content.json`.

### Add an admin page that isn't a resource

Follow `messages` in `admin/handlers.go`: a handler, a template in
`web/templates/pages/`, routes in `Routes`, and a sidebar link in
`layout.html`. Add a case to `TestAllPagesRender`.

## Testing

```sh
cd backend  && go test ./...    # store round-trips, shipped content.json validity, schema, media, sessions, templates, routes, API validation
cd frontend && npx tsc -b       # frontend types
```

Nothing external is needed. For an end-to-end check, `go run ./cmd/server`
in `backend/` and open `/admin` and `/api/v1/portfolio`.

## Frontend conventions

* `usePortfolio()` is the only data source; components never fetch.
* Theme tokens live in `frontend/src/index.css` (`@theme` for light, `.dark` for dark).
* Sections share `Section.tsx` (max-width 960 px, consistent spacing).
