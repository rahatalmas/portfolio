/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** Backend origin; empty/undefined = same origin. */
  readonly VITE_API_BASE?: string;
}
