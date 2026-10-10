# Portfolio API

Base path: `/api/v1`. All responses are JSON (`application/json; charset=utf-8`).
Every endpoint is public and read-only except `POST /messages`.

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
| `POST /messages` | Contact-form submission |

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
  "github_url": "https://github.com/yourhandle",
  "cv_url": "/uploads/cv.pdf",

  "meta_title": "Rahat Almas — Software Engineer",
  "meta_description": "Portfolio of …",

  "show_experience": true,
  "show_about": true,
  "show_work": true,
  "show_research": true,
  "show_contact": true,
  "show_contact_form": true,

  "nav_about": "About", "nav_work": "Work", "nav_research": "Research",
  "nav_contact": "Contact", "nav_cv_label": "CV",

  "hero_cv_label": "CV", "hero_github_label": "GitHub", "hero_email_label": "Email",
  "experience_title": "Experience",
  "stat_years_label": "Years", "stat_papers_label": "Papers", "stat_projects_label": "Projects",
  "full_cv_label": "Full CV",

  "about_title": "About", "education_label": "Education", "skills_label": "Skills",

  "work_title": "Selected work",
  "work_all_label": "All projects",
  "work_all_url": "https://github.com/yourhandle?tab=repositories",

  "research_title": "Research",
  "research_statement": "",
  "interests_label": "Interests",
  "interests": ["LLM inference systems", "Retrieval-augmented generation"],

  "contact_title": "Contact",
  "contact_email_label": "Email", "contact_location_label": "Location",
  "form_title": "Send a message",
  "form_subtitle": "Questions, collaborations, or opportunities.",
  "form_name_placeholder": "Name", "form_email_placeholder": "Email",
  "form_message_placeholder": "Message",
  "form_submit_label": "Send message", "form_sending_label": "Sending…",
  "form_success_title": "Message sent",
  "form_success_text": "Thanks — I'll get back to you soon.",
  "form_error_text": "Couldn't send. Email me instead.",

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

## `POST /messages`

Contact-form submission. Stored in the database and shown in the admin inbox.

Request (`Content-Type: application/json`):

```json
{ "name": "Ada Lovelace", "email": "ada@example.com", "message": "Hello!", "website": "" }
```

| Field | Rules |
| --- | --- |
| `name` | required, ≤ 120 characters |
| `email` | required, valid address, ≤ 254 characters |
| `message` | required, ≤ 5000 characters |
| `website` | honeypot — must be empty or absent; bots that fill it get a fake success |

Responses:

| Status | Body |
| --- | --- |
| `201 Created` | `{ "ok": true, "id": 12 }` |
| `400 Bad Request` | `{ "error": "invalid JSON body" }` |
| `422 Unprocessable Entity` | `{ "error": "email is not a valid address" }` |
| `429 Too Many Requests` | `{ "error": "too many messages, try again in a minute" }` — 5 per minute per IP |

Example:

```sh
curl -X POST http://localhost:8080/api/v1/messages \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ada","email":"ada@example.com","message":"Hi!"}'
```

---

## CORS

The site is hosted separately from this server, so in production set
`CORS_ORIGINS` to the site's origin (comma-separated for several, or `*`).
In development the Vite proxy makes requests same-origin and no CORS is
involved.

Upload URLs (`photo`, `cv_url`, and any link under `/uploads/`) are returned
as backend-relative paths; clients on another origin must prefix them with the
backend origin (the React client does this automatically).

## Admin

The management UI is not part of the JSON API. It lives at `/admin`, uses
session cookies, and is documented in `backend/README.md`.
