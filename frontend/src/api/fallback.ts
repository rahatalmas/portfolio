import content from "../../../backend/data/content.json";
import type { Paper, Portfolio, Profile, ProfileGroup, Project, Site } from "./types";

/**
 * Offline content: the repo's backend/data/content.json, bundled at build
 * time and converted into the API's shape. It paints first and stays if the
 * API is unreachable, so the site never shows an error. There is one source
 * of truth for content; nothing here needs to be kept in sync by hand.
 *
 * content.json stores every value as the admin does: strings, with "1"/"0"
 * for booleans, one item per line for lists, and comma-separated tags. The
 * converters below mirror backend/internal/schema (Field.Typed, SplitLines,
 * SplitComma) and the typed reads in backend/internal/store/lists.go.
 */

type Raw = Record<string, unknown>;

const str = (v: unknown): string =>
  v == null ? "" : typeof v === "string" ? v : typeof v === "boolean" ? (v ? "1" : "0") : String(v);
const bool = (v: unknown): boolean => str(v) === "1";
const lines = (v: unknown): string[] =>
  str(v)
    .split("\n")
    .map((s) => s.trim())
    .filter(Boolean);
const comma = (v: unknown): string[] =>
  str(v)
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean);
const num = (v: unknown): number => Number.parseInt(str(v), 10) || 0;

/** Settings keys the API returns as something other than a string. */
const BOOL_KEYS = new Set(["show_experience", "show_about", "show_work", "show_research", "show_contact"]);
const LINES_KEYS = new Set(["interests"]);

function site(raw: Raw): Site {
  const out: Raw = {};
  for (const [k, v] of Object.entries(raw)) {
    out[k] = BOOL_KEYS.has(k) ? bool(v) : LINES_KEYS.has(k) ? lines(v) : str(v);
  }
  return out as unknown as Site;
}

/** Rows without an id get a sequential one, as the backend does on load. */
function withIds(rows: Raw[] | undefined): (Raw & { id: number })[] {
  let next = Math.max(0, ...(rows ?? []).map((r) => num(r.id)));
  return (rows ?? []).map((r) => ({ ...r, id: num(r.id) || ++next }));
}

function profiles(rows: Raw[] | undefined): ProfileGroup[] {
  const groups: ProfileGroup[] = [];
  for (const r of withIds(rows)) {
    const name = str(r.group_name);
    let g = groups.find((x) => x.name === name);
    if (!g) groups.push((g = { name, items: [] }));
    const p: Profile = { id: r.id, label: str(r.label), href: str(r.href), icon: str(r.icon) };
    g.items.push(p);
  }
  return groups;
}

function projects(rows: Raw[] | undefined): Project[] {
  return withIds(rows).map((r) => ({
    id: r.id,
    slug: str(r.slug),
    title: str(r.title),
    summary: str(r.summary),
    tags: comma(r.tags),
    stack: lines(r.stack),
    featured: bool(r.featured),
    links: {
      ...(str(r.github_url) && { github: str(r.github_url) }),
      ...(str(r.live_url) && { live: str(r.live_url) }),
      ...(str(r.paper_url) && { paper: str(r.paper_url) }),
    },
  }));
}

function papers(rows: Raw[] | undefined): Paper[] {
  return withIds(rows)
    .map(
      (r): Paper => ({
        id: r.id,
        title: str(r.title),
        authors: lines(r.authors),
        venue: str(r.venue),
        year: num(r.year),
        status: (["published", "preprint", "in progress"].includes(str(r.status)) ? str(r.status) : "published") as Paper["status"],
        links: {
          ...(str(r.pdf_url) && { pdf: str(r.pdf_url) }),
          ...(str(r.code_url) && { code: str(r.code_url) }),
          ...(str(r.poster_url) && { poster: str(r.poster_url) }),
          ...(str(r.slides_url) && { slides: str(r.slides_url) }),
        },
      }),
    )
    .sort((a, b) => b.year - a.year);
}

const raw = content as unknown as Raw & { settings?: Raw } & Partial<Record<string, Raw[]>>;

export const FALLBACK: Portfolio = {
  site: site(raw.settings ?? {}),
  education: withIds(raw.education).map((r) => ({
    id: r.id,
    degree: str(r.degree),
    school: str(r.school),
    years: str(r.years),
    note: str(r.note),
  })),
  skills: withIds(raw.skills).map((r) => ({ id: r.id, name: str(r.name), items: lines(r.items) })),
  profiles: profiles(raw.profiles),
  experience: withIds(raw.experience).map((r) => ({
    id: r.id,
    title: str(r.title),
    company: str(r.company),
    period: str(r.period),
  })),
  projects: projects(raw.projects),
  papers: papers(raw.papers),
};

/**
 * Fill gaps in an API response from the fallback so a partial payload never
 * blanks the page. Only missing keys are filled: an empty string or empty
 * list from the API is a real choice made in the admin.
 */
export function withFallback(data: Partial<Portfolio> | null | undefined): Portfolio {
  return {
    site: { ...FALLBACK.site, ...(data?.site ?? {}) },
    education: data?.education ?? FALLBACK.education,
    skills: data?.skills ?? FALLBACK.skills,
    profiles: data?.profiles ?? FALLBACK.profiles,
    experience: data?.experience ?? FALLBACK.experience,
    projects: data?.projects ?? FALLBACK.projects,
    papers: data?.papers ?? FALLBACK.papers,
  };
}
