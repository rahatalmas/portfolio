import { usePortfolio } from "../api/PortfolioContext";
import type { Paper } from "../api/types";

const linkClass = "text-accent hover:underline";

/** One compact row: title + year, then authors · venue · links on a single line. */
export default function PaperEntry({ paper }: { paper: Paper }) {
  const { site } = usePortfolio();
  const { title, authors, venue, year, status, links } = paper;
  const linkList = [
    { label: "PDF", href: links.pdf },
    { label: "Code", href: links.code },
    { label: "Poster", href: links.poster },
    { label: "Slides", href: links.slides },
  ].filter((l): l is { label: string; href: string } => Boolean(l.href));

  return (
    <li className="border-b border-line py-3 first:pt-0 last:border-0 last:pb-0">
      <div className="flex items-baseline justify-between gap-4">
        <h3 className="text-[15px] font-medium">{title}</h3>
        <span className="shrink-0 font-mono text-xs text-muted">{year}</span>
      </div>
      <p className="mt-0.5 text-[13px] text-muted">
        {authors.map((a, i) => (
          <span key={`${a}-${i}`}>
            {a === site.name ? <span className="font-medium text-fg">{a}</span> : a}
            {i < authors.length - 1 ? ", " : ""}
          </span>
        ))}
        <span> · </span>
        <span className="italic">{venue}</span>
        {status !== "published" && (
          <span className="ml-1.5 rounded border border-line px-1 py-px font-mono text-[10px] not-italic">{status}</span>
        )}
        {linkList.map((l) => (
          <span key={l.label}>
            <span> · </span>
            <a href={l.href} target="_blank" rel="noreferrer" className={linkClass}>
              {l.label}
            </a>
          </span>
        ))}
      </p>
    </li>
  );
}
