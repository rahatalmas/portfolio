import { ArrowLeft } from "lucide-react";
import { usePortfolio } from "../api/PortfolioContext";

/**
 * Shown for any path the site does not have. The server answers these with a
 * 404 status and this page as the body, so crawlers and visitors agree.
 */
export default function NotFound() {
  const { site } = usePortfolio();
  const path = window.location.pathname;

  return (
    <section className="mx-auto flex min-h-[70vh] max-w-[960px] flex-col justify-center px-6 py-20">
      <p className="font-mono text-[13px] text-accent">404</p>
      <h1 className="mt-2 text-[28px] font-bold leading-tight tracking-tight md:text-4xl">Page not found</h1>
      <p className="mt-4 max-w-[52ch] leading-relaxed text-muted">
        There is nothing at{" "}
        <span className="font-mono text-[13px] break-all text-fg">{path}</span>. It may have moved, or the link that
        brought you here may be wrong.
      </p>
      <div className="mt-6 flex flex-wrap items-center gap-4">
        <a
          href="/"
          className="inline-flex items-center gap-2 rounded-md bg-accent px-3.5 py-1.5 text-sm font-medium text-on-accent"
        >
          <ArrowLeft size={16} aria-hidden />
          Back to {site.name}
        </a>
        {site.email && (
          <a href={`mailto:${site.email}`} className="text-sm text-muted">
            Think this is a mistake? Email me.
          </a>
        )}
      </div>
    </section>
  );
}
