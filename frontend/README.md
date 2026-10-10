# Frontend

React 19 + Vite + TypeScript + Tailwind v4. Renders everything from the
backend's `GET /api/v1/portfolio`; it has no content of its own.

```sh
npm install
npm run dev      # http://localhost:5173, proxies /api and /uploads to http://localhost:8080
npm run build    # → dist/, deploy to any static host
```

## Production build

The backend does not serve the site. Point the build at the backend and host
`dist/` anywhere static:

```sh
cp .env.example .env.production     # VITE_API_BASE=https://api.your-portfolio.example
npm run build
# upload dist/ to Vercel / Netlify / Cloudflare Pages / GitHub Pages / nginx
```

The backend must list the site's origin in `CORS_ORIGINS`. Upload URLs from
the API (`/uploads/…`) are prefixed with `VITE_API_BASE` automatically.

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

`src/api/client.ts` reads `VITE_API_BASE`; leave it unset in development.
