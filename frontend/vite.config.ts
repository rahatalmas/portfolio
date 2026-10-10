import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// In development the Go backend runs on :8080; the proxy forwards API and
// upload requests there so the app can use same-origin URLs everywhere.
// Change this if you run the backend on another port.
const backend = "http://localhost:8080";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    proxy: {
      "/api": backend,
      "/uploads": backend,
    },
  },
});
