import type { MessageInput, Portfolio } from "./types";

/**
 * Base URL of the backend. Empty means same origin, which is the case both
 * in production (the Go server serves the built app) and in development (the
 * Vite proxy forwards /api and /uploads). Set VITE_API_BASE to point elsewhere.
 */
const BASE = (import.meta.env.VITE_API_BASE as string | undefined)?.replace(/\/$/, "") ?? "";

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
export function fetchPortfolio(): Promise<Portfolio> {
  return request<Portfolio>("/api/v1/portfolio");
}

/** Submit the contact form. */
export function sendMessage(input: MessageInput): Promise<{ ok: boolean; id?: number }> {
  return request("/api/v1/messages", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
}
