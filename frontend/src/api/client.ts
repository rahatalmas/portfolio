import type { Portfolio } from "./types";

/**
 * Origin of the Go backend. Empty means same origin, which is the normal
 * case: the Go server serves this built site at /, and in development the
 * Vite proxy forwards /api and /uploads to :8085. Set VITE_API_BASE only when
 * hosting the site on another origin (see .env.example).
 */
const BASE = (import.meta.env.VITE_API_BASE as string | undefined)?.replace(/\/$/, "") ?? "";

/** Path under which the backend serves uploaded media. */
const UPLOADS = "/uploads/";

/**
 * Uploaded files (photo, CV, …) come back from the API as backend-relative
 * paths. Same-origin they work as-is; only when VITE_API_BASE points at
 * another origin are they made absolute. Applied to every string in the
 * response, so links in projects/papers that point at uploads work too.
 */
export function absolutizeUploads<T>(value: T): T {
  if (!BASE) return value;
  if (typeof value === "string") {
    return (value.startsWith(UPLOADS) ? BASE + value : value) as T;
  }
  if (Array.isArray(value)) {
    return value.map(absolutizeUploads) as T;
  }
  if (value && typeof value === "object") {
    const out: Record<string, unknown> = {};
    for (const [k, v] of Object.entries(value as Record<string, unknown>)) out[k] = absolutizeUploads(v);
    return out as T;
  }
  return value;
}

/** Error raised for non-2xx responses, carrying the server's message. */
export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    ...init,
    headers: { Accept: "application/json", ...(init?.headers ?? {}) },
  });
  if (!res.ok) {
    let message = `${res.status} ${res.statusText}`;
    try {
      const body = (await res.json()) as { error?: string };
      if (body.error) message = body.error;
    } catch {
      /* non-JSON error body */
    }
    throw new ApiError(res.status, message);
  }
  return (await res.json()) as T;
}

/** Fetch everything the page renders in one call. */
export async function fetchPortfolio(): Promise<Portfolio> {
  return absolutizeUploads(await request<Portfolio>("/api/v1/portfolio"));
}
