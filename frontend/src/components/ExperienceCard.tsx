import { ArrowUpRight } from "lucide-react";
import { usePortfolio } from "../api/PortfolioContext";
import type { Role } from "../api/types";

/** "Jun 2024 – Present" -> "2024 – Now", "Jun 2023 – Aug 2023" -> "2023" */
function shortPeriod(period: string) {
  const [start, end = ""] = period.split(/\s*[–-]\s*/);
  const from = start.match(/\d{4}/)?.[0] ?? start;
  const to = /present|now/i.test(end) ? "Now" : (end.match(/\d{4}/)?.[0] ?? end);
  return from === to || !to ? from : `${from} – ${to}`;
}

/** Whole years since the earliest role started. */
function yearsOfExperience(experience: Role[]) {
  const starts = experience.map((r) => Number(r.period.match(/\d{4}/)?.[0])).filter(Boolean);
  if (!starts.length) return 0;
  return new Date().getFullYear() - Math.min(...starts);
}

/** Compact timeline card for the hero: roles on top, quick stats + CV link at the bottom. */
export default function ExperienceCard({ className = "" }: { className?: string }) {
  const { site, experience, papers, projects } = usePortfolio();
  const stats = [
    { value: `${yearsOfExperience(experience)}+`, label: site.stat_years_label },
    { value: String(papers.length), label: site.stat_papers_label },
    { value: String(projects.length), label: site.stat_projects_label },
  ];

  return (
    <div id="experience" className={`flex h-full flex-col rounded-2xl border border-line bg-card p-4 ${className}`}>
      <h2 className="mb-3 font-mono text-xs font-medium uppercase tracking-wider text-muted">{site.experience_title}</h2>
      <ol className="relative ml-1 flex flex-col gap-4 border-l border-line">
        {experience.map((r, i) => (
          <li key={r.id} className="relative pl-5">
            <span
              aria-hidden
              className={`absolute top-1.5 -left-[5px] h-2.5 w-2.5 rounded-full border-2 ${
                i === 0 ? "border-accent bg-accent" : "border-line bg-card"
              }`}
            />
            <div className="flex items-baseline justify-between gap-3">
              <span className="truncate text-sm font-medium">{r.title}</span>
              <span className="shrink-0 font-mono text-[11px] text-muted">{shortPeriod(r.period)}</span>
            </div>
            <div className="truncate text-xs text-muted">{r.company}</div>
          </li>
        ))}
      </ol>

      {/* Bottom: stats + CV, pinned to the card's base */}
      <div className="mt-auto pt-4">
        <dl className="grid grid-cols-3 divide-x divide-line border-y border-line py-2 text-center">
          {stats.map((s) => (
            <div key={s.label}>
              <dt className="sr-only">{s.label}</dt>
              <dd className="text-lg font-semibold tracking-tight">{s.value}</dd>
              <dd className="font-mono text-[10px] uppercase tracking-wider text-muted">{s.label}</dd>
            </div>
          ))}
        </dl>
        {site.cv_url && (
          <a
            href={site.cv_url}
            target="_blank"
            rel="noreferrer"
            className="mt-2.5 flex items-center justify-between text-sm text-muted hover:text-accent"
          >
            {site.full_cv_label}
            <ArrowUpRight size={15} />
          </a>
        )}
      </div>
    </div>
  );
}
