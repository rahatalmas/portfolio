# Portfolio Site Plan

As of 2026-10-10 · Online copy: https://claude.ai/code/artifact/f6c3fa0b-ab6e-434f-9d1e-f0294eb71176

## Goals & audiences

One site, two readers: an admissions committee skimming for research potential, and a hiring manager skimming for shipped systems. Both spend under 60 seconds on the first visit, so the home page must answer "who is this, what have they built, what do they study" above the fold.

| Reader | What they look for | Where it must be | Time budget |
| --- | --- | --- | --- |
| PhD / MS committee | Research interests, publications, advisors, GPA/thesis, technical depth | Research, Academics, CV (PDF) | 30–60 s, then CV |
| AI engineer hiring | ML/LLM projects with real results, code on GitHub, stack | Work (projects), Experience, GitHub link | 30 s, then GitHub |
| Backend engineer hiring | Systems built, scale, languages, infra, ownership | Experience, Work (projects) | 30 s, then resume |
| Recruiter / cold outreach | Name, title, location, one-line pitch, contact | Hero section | 10 s |

Design consequence: the site is one scrollable page with a fixed navbar and short sections, plus two deeper pages (Research, Projects) that only people who care will open. No one is made to click to find the basics.

## Navbar

Five items plus a CV button, left-aligned name, right-aligned links. Anything beyond six items reads as a menu, not a portfolio.

| Position | Item | Behaviour |
| --- | --- | --- |
| Left | Your name (monogram on mobile) | Scrolls to top |
| Right 1 | About | Smooth-scroll to section |
| Right 2 | Work | Smooth-scroll; section links to /projects for the full list |
| Right 3 | Research | Smooth-scroll; section links to /research for papers + interests |
| Right 4 | Experience | Smooth-scroll |
| Right 5 | Contact | Smooth-scroll |
| Right, last | CV (outlined button) | Opens PDF in new tab |

Rules:

- Sticky, 56–64 px tall, transparent at the top and gains a hairline border + blur once the page scrolls.
- Active section is underlined or tinted as you scroll (scroll-spy), so the navbar doubles as a progress indicator.
- Theme toggle (light/dark) as a small icon, not a labelled item.
- Mobile: name + CV + hamburger; the menu is a full-height sheet with the same five links in large type. No nested menus anywhere.
- Education and Academics do not get their own navbar item: they live inside About (summary) and the CV (detail). Committees go to the CV anyway; keeping the navbar short is worth more than the extra link.

Alternative to consider: drop Contact from the navbar and put email + socials in the hero and footer instead. That leaves four links and looks cleaner, at the cost of one fewer obvious call to action.

## Sections

Seven sections on the home page, each one screen or less, in the order a stranger would ask questions: who, what have you built, what do you study, where have you worked, what else, how to reach you.

| # | Section | Purpose | Height | Content cap |
| --- | --- | --- | --- | --- |
| 1 | Hero | Name, one-line role ("AI & backend engineer · researching X"), 2-sentence pitch, photo optional, CV + GitHub + email buttons | 100 vh | 40 words |
| 2 | About | 1 paragraph bio + a compact education block (degree, school, year, GPA/thesis) + skill tags grouped by AI / backend / tools | ½ screen | 120 words + 3 degrees |
| 3 | Work (selected projects) | 3–4 project cards: name, one-line result, stack tags, links to code/demo/write-up. "All projects →" link | 1 screen | 4 cards |
| 4 | Research | 2–3 publications or research projects (title, venue/year, authors, PDF/code links) + a short line of interests. "More →" link | ½–1 screen | 3 entries |
| 5 | Experience | Vertical timeline: role, company, dates, 2 bullets each. Internships and TA/RA roles included | 1 screen | 4 roles |
| 6 | Beyond | Awards, talks, teaching, open-source, hobbies: one line each. Optional | ½ screen | 6 lines |
| 7 | Contact + footer | Email, GitHub, LinkedIn, Google Scholar, location; "last updated" date | ¼ screen | — |

Sub-pages (reached from the Work and Research sections, not the navbar):

- `/projects` — every project as the same card, filterable by tag (AI, backend, research). Each card can open a `/projects/<slug>` write-up: problem, approach, result, what you'd do differently.
- `/research` — full publication list, research statement (200–300 words), areas of interest, and any preprints or thesis.
- `/cv.pdf` — the single source of truth for committees; keep it versioned.

Sections to cut, by default: testimonials, a blog (add later only if you will actually write), a stats counter, and a separate Skills page. Skills belong as tags next to the work that proves them.

## Design system

Minimal means few decisions, applied consistently: one typeface pair, one accent colour, one column, generous whitespace, and motion only where it tells the reader something.

| Element | Decision | Why |
| --- | --- | --- |
| Layout | Single centred column, max 720 px for text, 960 px for card grids; 24 px gutters on mobile | Reads like a paper, not a landing page |
| Type | Inter (or Geist) for UI and body; JetBrains Mono for tags, dates, code. Body 17 px / 1.6 line-height | Clean, free, loads fast; mono signals "engineer" |
| Scale | H1 40–48 px, H2 28 px, H3 20 px, body 17 px, small 14 px | Five sizes is enough |
| Colour | Near-black on off-white (light); off-white on near-black (dark); one accent (e.g. #2563eb blue or a muted green) for links and the active nav item | Zero colour noise; the work carries the visual interest |
| Spacing | 8 px base; section padding 96 px top/bottom on desktop, 64 px on mobile | Breathing room is the main "design" |
| Cards | No shadows; 1 px border, 12 px radius, hover lifts the border colour to the accent | Flat and fast to scan |
| Icons | Lucide, 16–20 px, only for external links and theme toggle | Icons without labels become noise |
| Motion | Fade-up on first scroll into view (200 ms, once), nav underline slide, nothing else | Motion that repeats is distracting |
| Images | Only a headshot (optional) and project thumbnails; no stock art, no gradients, no 3D blobs | Every pixel should be yours |
| Dark mode | Respects system setting, toggle persists | Expected by engineers in 2026 |

Reference sites in this style: [Brittany Chiang](https://brittanychiang.com), [Lee Robinson](https://leerob.com), [Paco Coursey](https://paco.me), and academic ones like [Andrej Karpathy](https://karpathy.ai) and [Lilian Weng](https://lilianweng.github.io). The first three show the engineer register, the last two the researcher register; the goal is the midpoint.

Accessibility floor: contrast ≥ 4.5:1, keyboard-navigable nav, visible focus rings, `prefers-reduced-motion` disables the fade-ups.

## Showcasing each content type

Every item follows the same formula: what it is, what it did (a number where possible), and a link to proof. No paragraphs where a line will do.

| Content | Format | Fields | Tip |
| --- | --- | --- | --- |
| Bio | Hero line + one About paragraph, first person | Who you are now, what you build, what you want next (grad school or role) | Write two versions of the last sentence and pick per season; committees like "I want to study X with Y" |
| Academics | Compact education block inside About; detail in CV | Degree, institution, years, GPA (if ≥ 3.5/4), thesis title, 2–3 key courses | Don't list every course; link the thesis PDF |
| Works / projects | Card grid, 3–4 featured | Name, one-line outcome ("cut p95 latency 40%"), stack tags, links: GitHub · demo · write-up | Lead with results, not features; each featured project needs a public repo or demo |
| Experience | Vertical timeline, newest first | Role, company, dates, 2 bullets (scope + impact), stack tags | Internships, RA/TA, and freelance all count; keep bullets under 20 words |
| Research | Publication list + statement on /research | Title, authors (you in bold), venue, year, links: PDF · code · poster; 200–300-word research statement | Preprints and in-progress work are fine if labelled; link Google Scholar |
| Areas of interest | 4–6 short phrases, as tags under the research statement | e.g. "LLM inference systems", "retrieval", "distributed storage" | Interests must match at least one project or paper above them |
| Others | "Beyond" section, one line each | Awards, talks, teaching, open-source contributions, hackathons, one hobby | Cap at six; if it needs explaining, it goes in the CV instead |
| CV | Single PDF, linked from navbar, hero, footer | Academic CV (2–3 pages) for committees; optionally a 1-page resume for jobs | Both can live at /cv.pdf and /resume.pdf; keep dates in sync with the site |

Content sources to keep in one place: a `content/` folder of Markdown/JSON (projects, papers, roles). The site renders from it, and the CV is regenerated from the same data, so nothing drifts.

## Tech stack & structure

Recommendation: Astro + Tailwind, content in Markdown, deployed to GitHub Pages or Vercel. It ships zero JavaScript by default, which keeps the page fast and the Lighthouse score near 100, and the project itself reads as a small, well-built system on your GitHub.

| Option | Fit | Trade-off |
| --- | --- | --- |
| Astro + Tailwind (recommended) | Static, content-first, Markdown collections for projects and papers, islands only where needed | Smaller ecosystem than Next; fine for a portfolio |
| Next.js (App Router) + Tailwind | Familiar to hiring managers, easy to add a blog or API route later | Heavier; more JS than a portfolio needs |
| Plain HTML + CSS | Zero build, nothing to maintain | Repeating cards and lists by hand gets tedious; no content collections |

Proposed structure:

```
portfolio/
  src/
    content/
      projects/*.md      # one file per project (frontmatter: title, summary, tags, links, featured)
      papers/*.md        # one file per publication
      experience.json    # roles
      site.json          # name, tagline, socials, interests
    components/          # Nav, Hero, ProjectCard, Timeline, PaperEntry, Footer
    layouts/Base.astro
    pages/
      index.astro        # home: 7 sections
      projects/index.astro
      projects/[slug].astro
      research.astro
  public/
    cv.pdf
    resume.pdf
```

Other pieces:

- Hosting: GitHub Pages via Actions (free, custom domain) or Vercel (free, preview URLs).
- Domain: `firstnamelastname.com` or `.dev`; set it before sharing any link in applications.
- Analytics: none, or a privacy-light one (Plausible/Umami) if you want to see which projects get opened.
- SEO basics: title, description, Open Graph image (a plain card with your name), `sitemap.xml`, Google Scholar and GitHub as verified links.
- CV: keep as LaTeX or Typst in the repo so edits are versioned; export to `public/cv.pdf` on build.

## Open questions

Answer these and the build can start; everything else above has a default I'm happy to go with.

- [ ] Stack: Astro (recommended), Next.js, or plain HTML?
- [ ] Navbar: keep Contact as a fifth link, or move contact to hero + footer and run four links?
- [ ] Which reader comes first this year: grad applications or job search? It decides the hero's last sentence and whether Research sits above or below Work.
- [ ] Do you have publications or preprints yet? If none, Research becomes "Research interests" with projects as evidence, and the /research page waits.
- [ ] Headshot: yes or no?
- [ ] Dark mode: system-default with toggle (recommended), or light only?
- [ ] One CV PDF or CV + 1-page resume?
- [ ] Domain name you own or want, and hosting preference (GitHub Pages vs Vercel).
- [ ] Content I can start from: existing resume/CV, GitHub profile, LinkedIn, list of projects and papers.

Suggested next steps once answered: scaffold the project, build the Nav + Hero + one ProjectCard with placeholder content, review the look together, then fill in the rest from your real content.
