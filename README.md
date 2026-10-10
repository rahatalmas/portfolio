# Portfolio

Single-page portfolio site (React + Vite + TypeScript + Tailwind v4) served
by a Go + MySQL backend that also provides the content API and an admin panel
for editing every text, image, project, paper and section.

```
frontend/   React app — the public site; renders GET /api/v1/portfolio
backend/    Go server — serves the built site at /, JSON API, admin UI (/admin), uploads
docs/       API.md (endpoints), ARCHITECTURE.md (how it fits together, how to extend)
Dockerfile  Builds the site and the server into one image
```

One deployment: the Go binary serves the built site, the API and the admin
from a single origin, so there is no CORS and nothing to host separately.

## Run locally

You need Node 20+, Go 1.24+, and a MySQL 8 / MariaDB 10.6+ database.

```sh
# Backend (terminal 1)
cd backend
cp .env.example .env         # set DB_* and SESSION_SECRET
go run ./cmd/server          # http://localhost:8080/admin — password is printed on first run

# Frontend (terminal 2)
cd frontend
npm install
npm run dev                  # http://localhost:5173 — proxies /api and /uploads to :8080
```

Log in at `http://localhost:8080/admin`, upload your photo and CV under
*Media*, then edit *Settings* and the content sections. The site updates on
reload.

## Deploying

Easiest: `docker compose up --build` in the repo root. It builds the React
site and the Go server into one image and starts it with MySQL on :8080.
Change the passwords and `SESSION_SECRET` in `docker-compose.yml` first, and
set `SECURE_COOKIES=true` once it sits behind an HTTPS proxy (Caddy, nginx).

Without Docker: `cd frontend && npm run build`, then
`cd ../backend && go build -o server ./cmd/server` and run the binary with
`FRONTEND_DIR` pointing at `frontend/dist`, `UPLOAD_DIR` on persistent
storage, and `SESSION_SECRET` set. The site, API and admin are all on the
same port.

## Editing content

Everything is in the admin:

| Where | What you control |
| --- | --- |
| Settings → Identity | Name, handle, role, bio, photo, email, location, GitHub, CV |
| Settings → SEO | Browser title, meta description |
| Settings → Sections | Show/hide About, Work, Research, Contact, the experience card, the form |
| Settings → Navbar / Hero / About / Work / Research / Contact / Footer | Every label and micro-text |
| Education, Skills, Profiles, Experience, Projects, Papers | Add, edit, delete, reorder |
| Messages | Contact-form inbox |
| Media | Upload images and PDFs |

## Theme

Colours and fonts are tokens in `frontend/src/index.css` (`@theme` for light,
`.dark` for dark). Change the accent by editing `--color-accent` in both
places. The admin uses the same palette
(`backend/internal/admin/web/static/admin.css`).

## Docs

* [`docs/API.md`](docs/API.md) — every endpoint with examples
* [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — design, data model, request flows, how to add fields/resources
* [`backend/README.md`](backend/README.md) — configuration, admin, deployment
* [`frontend/README.md`](frontend/README.md) — app structure, build
* [`PLAN.md`](PLAN.md) — the original design plan
