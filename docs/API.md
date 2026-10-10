# Portfolio API

Base path: `/api/v1`. All responses are JSON (`application/json; charset=utf-8`).
Every endpoint is public and read-only.

GET responses carry `Cache-Control: public, max-age=30`.

Errors are always `{"error": "<message>"}` with an appropriate status code.

| Endpoint | Purpose |
| --- | --- |
| `GET /health` | Liveness check |
| `GET /portfolio` | Everything the site renders, in one call |
| `GET /site` | Single-value settings (texts, labels, flags, image URLs) |
| `GET /education` | Degrees |
| `GET /skills` | Skill groups |
| `GET /profiles` | Profile links, grouped |
| `GET /experience` | Roles for the hero timeline card |
| `GET /projects` | Projects (`?featured=1` for featured only) |
| `GET /papers` | Publications |

---

## `GET /health`

```json
{ "status": "ok" }
```

## `GET /portfolio`

The frontend makes exactly one request, to this endpoint.

```json
{
  "site":       { ...see GET /site... },
  "education":  [ ...see GET /education... ],
  "skills":     [ ...see GET /skills... ],
  "profiles":   [ ...see GET /profiles... ],
  "experience": [ ...see GET /experience... ],
  "projects":   [ ...see GET /projects... ],
  "papers":     [ ...see GET /papers... ]
}
```

## `GET /site`

A flat object. Keys are the setting keys defined in
`backend/internal/schema/settings.go`; each is typed by its field kind:

| Kind | JSON type |
| --- | --- |
| text, textarea, url, email, media, select | `string` |
| bool | `boolean` |
| number | `number` |
| lines, tags | `string[]` |

```json
{
  "name": "Rahat Almas",
  "handle": "rahat",
  "initials": "RA",
  "role": "Software Engineer",
  "bio": "I'm a software engineer and researcher …",
  "photo": "/uploads/profile-2.jpg",
  "location": "City, Country",
  "email": "you@example.com",
  "cv_url": "/uploads/cv.pdf",

  "meta_title": "Rahat Almas — Software Engineer",
  "meta_description": "Portfolio of …",

  "show_experience": true,
  "show_about": true,
  "show_work": true,
  "show_research": true,

  "nav_about": "About", "nav_work": "Work", "nav_research": "Research",

  "experience_title": "Experience",

  "about_title": "About", "education_label": "Education", "skills_label": "Skills",

  "work_title": "Selected work",
  "work_all_label": "All projects",
  "work_all_url": "https://github.com/yourhandle?tab=repositories",

  "research_title": "Research",
  "research_statement": "",
  "interests_label": "Interests",
  "interests": ["LLM inference systems", "Retrieval-augmented generation"],

  "contact_email_label": "Email", "contact_location_label": "Location",

  "footer_updated": "October 2026",
  "footer_back_to_top": "Back to top"
}
```

The full key list, with labels and help text, is the single source of truth in
`backend/internal/schema/settings.go`; the TypeScript mirror is
`frontend/src/api/types.ts` (`Site`).

## `GET /education`

```json
[
  { "id": 1, "degree": "B.Sc. Computer Science", "school": "University Name",
    "years": "2020 – 2024", "note": "GPA 3.9/4.0 · Thesis: …" }
]
```

## `GET /skills`

```json
[
  { "id": 1, "name": "AI", "items": ["PyTorch", "Transformers", "vLLM"] },
  { "id": 2, "name": "Backend", "items": ["Python", "Go", "FastAPI"] }
]
```

## `GET /profiles`

Links sharing a group name are grouped; groups appear in the order of their
first link.

```json
[
  { "name": "Professional",
    "items": [
      { "id": 1, "label": "GitHub",   "href": "https://github.com/yourhandle",      "icon": "github" },
      { "id": 2, "label": "LinkedIn", "href": "https://linkedin.com/in/yourhandle", "icon": "linkedin" }
    ] },
  { "name": "Coding", "items": [ { "id": 7, "label": "LeetCode", "href": "…", "icon": "code" } ] }
]
```

`icon` is one of: `github`, `linkedin`, `x`, `scholar`, `orcid`, `research`,
`code`, `trophy`, `chef`, `terminal`, `chart`, `mail`, `globe`.

## `GET /experience`

Display order (the admin's order; newest first by convention).

```json
[
  { "id": 1, "title": "Software Engineer", "company": "Company A · ML Platform", "period": "Jun 2024 – Present" },
  { "id": 2, "title": "Research Assistant", "company": "Systems Lab", "period": "Sep 2023 – May 2024" }
]
```

## `GET /projects`

Query: `featured=1` returns only featured projects.

```json
[
  { "id": 1, "slug": "fast-rag", "title": "FastRAG",
    "summary": "Retrieval pipeline that cut p95 query latency 40% …",
    "tags": ["AI", "Backend"],
    "stack": ["Python", "FastAPI", "pgvector", "Redis"],
    "featured": true,
    "links": { "github": "https://github.com/…", "live": "https://…" } }
]
```

`links` keys are omitted when empty. `tags` values are `AI`, `Backend`, `Research`.

## `GET /papers`

Newest year first, then admin order.

```json
[
  { "id": 1, "title": "Length-Aware Scheduling for Memory-Efficient LLM Inference",
    "authors": ["Rahat Almas", "Co Author"],
    "venue": "Workshop on ML Systems, NeurIPS", "year": 2025,
    "status": "published",
    "links": { "pdf": "https://arxiv.org/abs/…", "code": "https://github.com/…" } }
]
```

`status` is `published`, `preprint` or `in progress`. `links` keys (`pdf`,
`code`, `poster`, `slides`) are omitted when empty.

## CORS

The site is served by this server, so browser requests are same-origin and no
CORS is involved (in development the Vite proxy achieves the same). If you
host a client on another origin, list it in `CORS_ORIGINS` (comma-separated,
or `*`).

Upload URLs (`photo`, `cv_url`, and any link under `/uploads/`) are returned
as server-relative paths; a client on another origin must prefix them with
this server's origin (the React client does this when `VITE_API_BASE` is set).

## Admin

The management UI is not part of the JSON API. It lives at `/admin`, uses
session cookies, and is documented in `backend/README.md`.
