# Portfolio

Single-page portfolio site (React + Vite + TypeScript + Tailwind v4) served
by a small Go backend that also provides the content API and an admin panel
for editing every text, image, project, paper and section.

There is no database. All content is one JSON file, `backend/data/content.json`,
which you can edit in the admin panel or by hand. The admin account is one more
JSON file next to it.

```
frontend/   React app — the public site; renders GET /api/v1/portfolio
backend/    Go server — serves the built site at /, JSON API, admin UI (/admin), uploads
backend/data/content.json   The content. Tracked in git; edit via /admin or a text editor.
docs/       API.md (endpoints), ARCHITECTURE.md (how it fits together, how to extend)
```

One deployment: the Go binary serves the built site, the API and the admin
from a single origin, so there is no CORS and nothing to host separately.

## Run locally

You need Node 20+ and Go 1.26+. Nothing else.

```sh
cd frontend && npm install && npm run build   # once, and again after changing frontend code
cd ../backend && go run ./cmd/server          # http://localhost:8080 — admin password is printed on first run
```

That is the whole site: `http://localhost:8080` is the page, `/admin` is the
editor, `/api/v1/portfolio` is the JSON. Log in, upload your photo and CV
under *Media*, then edit *Settings* and the content sections. The site
updates on reload. Every save rewrites `backend/data/content.json`, so you
can commit your content like code.

Only when you are editing the React code itself, also run the dev server in
a second terminal for hot reload:

```sh
cd frontend && npm run dev   # http://localhost:5173 — proxies /api and /uploads to :8080
```

## Deploying

It is one binary plus three folders. On any Linux box:

```sh
cd frontend && npm run build                                  # → frontend/dist
cd ../backend && GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
```

Copy `server`, `frontend/dist`, and `backend/data` to the machine, then run
the binary with `FRONTEND_DIR` pointing at `dist`, `DATA_DIR` at `data`,
`UPLOAD_DIR` somewhere persistent, `SESSION_SECRET` set, and
`SECURE_COOKIES=true` behind an HTTPS reverse proxy (Caddy, nginx). Keep it
running with a systemd unit. The site, API and admin are all on one port.

Back up by copying the `data` and `uploads` directories.

## Editing content

Everything is in the admin, or in `backend/data/content.json`:

| Where | What you control |
| --- | --- |
| Settings → Identity | Name, handle, role, bio, photo, email, location, GitHub, CV |
| Settings → SEO | Browser title, meta description |
| Settings → Sections | Show/hide About, Work, Research, Connect and the experience column |
| Settings → Navbar / Hero / About / Work / Research / Connect / Footer | Every label and micro-text |
| Education, Skills, Profiles, Experience, Projects, Papers | Add, edit, delete, reorder |
| Media | Upload images and PDFs |

The site also bundles `content.json` at build time, so it renders even if the
API is unreachable.

## Theme

Colours and fonts are tokens in `frontend/src/index.css` (`@theme` for light,
`.dark` for dark). Change the accent by editing `--color-accent` in both
places. The admin uses the same palette
(`backend/internal/admin/web/static/admin.css`).

## Docs

* [`docs/API.md`](docs/API.md) — every endpoint with examples
* [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — design, data files, request flows, how to add fields/resources
* [`backend/README.md`](backend/README.md) — configuration, admin, deployment
* [`frontend/README.md`](frontend/README.md) — app structure, build
* [`PLAN.md`](PLAN.md) — the original design plan
