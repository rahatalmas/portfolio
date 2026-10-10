import { usePortfolio } from "../api/PortfolioContext";
import type { Role } from "../api/types";
import Profiles from "./Profiles";

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

const eyebrow = "mb-2 font-mono text-xs font-medium uppercase tracking-wider text-muted";

/**
 * The hero's right-hand column: roles, then the profile links grouped by
 * row, then one line of numbers pinned to the bottom so it ends level with
 * the hero buttons. No box, just a hairline running the full height.
 */
export default function ExperienceCard({ className = "" }: { className?: string }) {
  const { site, experience, profiles, papers, projects } = usePortfolio();
  const stats = [
    { value: `${yearsOfExperience(experience)}+`, label: site.stat_years_label },
    { value: String(papers.length), label: site.stat_papers_label },
    { value: String(projects.length), label: site.stat_projects_label },
  ];

  return (
    <div
      id="experience"
      className={`flex h-full flex-col border-t border-line pt-5 md:border-t-0 md:border-l md:pt-0 md:pl-8 ${className}`}
    >
      {experience.length > 0 && (
        <>
          <h2 className={eyebrow}>{site.experience_title}</h2>
          <ol className="space-y-2.5">
            {experience.map((r) => (
              <li key={r.id}>
                <div className="flex items-baseline justify-between gap-3">
                  <span className="truncate text-sm font-medium">{r.title}</span>
                  <span className="shrink-0 font-mono text-[11px] text-muted">{shortPeriod(r.period)}</span>
                </div>
                <div className="truncate text-xs text-muted">{r.company}</div>
              </li>
            ))}
          </ol>
        </>
      )}

      {profiles.length > 0 && (
        <div className={experience.length > 0 ? "mt-5" : ""}>
          <h2 className={eyebrow}>{site.profiles_title}</h2>
          <Profiles />
        </div>
      )}

      <div className="mt-5 flex-1" />
      <dl className="flex flex-wrap gap-x-4 gap-y-1 border-t border-line pt-3 font-mono text-[11px] text-muted">
        {stats.map((s) => (
          <div key={s.label} className="flex items-baseline gap-1">
            <dd className="font-medium text-fg">{s.value}</dd>
            <dt className="lowercase">{s.label}</dt>
          </div>
        ))}
      </dl>
    </div>
  );
}
