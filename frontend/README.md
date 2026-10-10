# Frontend

React 19 + Vite + TypeScript + Tailwind v4. Renders everything from the
backend's `GET /api/v1/portfolio`; it has no content of its own.

```sh
npm install
npm run dev      # http://localhost:5173, proxies /api and /uploads to http://localhost:8080
npm run build    # → dist/, which the Go server serves at /
```

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

To point at a backend on another origin, set `VITE_API_BASE=https://api.example.com`
at build time and allow that origin with `CORS_ORIGINS` on the server.
