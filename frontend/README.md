# Frontend

React 19 + Vite + TypeScript + Tailwind v4. Renders everything from the
backend's `GET /api/v1/portfolio`; it has no content of its own.

```sh
npm install
npm run dev      # http://localhost:5173, proxies /api and /uploads to http://localhost:8080
npm run build    # → dist/, which the Go server serves at /
```

## Production build

The Go server serves `dist/` itself (`FRONTEND_DIR`, default `../frontend/dist`),
so a production deploy is just `npm run build` followed by starting the
backend, or `docker compose up --build` from the repo root, which does both.
No `VITE_API_BASE` and no CORS: the site and the API share one origin.

If you ever want to host the site elsewhere instead, set `VITE_API_BASE` in
`.env.production` before building and list the site's origin in the backend's
`CORS_ORIGINS`. Upload URLs (`/uploads/…`) are prefixed automatically.

| Path | What |
| --- | --- |
| `src/api/types.ts` | TypeScript mirror of the API (see `../docs/API.md`) |
| `src/api/client.ts` | `fetchPortfolio()`, `sendMessage()` |
| `src/api/PortfolioContext.tsx` | Loads once, exposes `usePortfolio()` |
| `src/App.tsx` | Loading / error / page shell; sets title + meta description |
| `src/pages/Home.tsx` | Section order and visibility toggles |
| `src/components/` | Nav, Hero, ExperienceCard, About, Work, Research, Contact, ContactForm, Footer… |
| `src/index.css` | Theme tokens (light in `@theme`, dark in `.dark`) |
| `public/profile.jpg` | Placeholder photo used until you upload your own |

`src/api/client.ts` reads `VITE_API_BASE`; leave it unset unless hosting the site separately.
