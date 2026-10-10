import { usePortfolio } from "../api/PortfolioContext";
import type { ProfileGroup } from "../api/types";
import Profiles from "./Profiles";

/** "Jun 2024 – Present" -> "2024 – Now", "Jun 2023 – Aug 2023" -> "2023" */
function shortPeriod(period: string) {
  const [start, end = ""] = period.split(/\s*[–-]\s*/);
  const from = start.match(/\d{4}/)?.[0] ?? start;
  const to = /present|now/i.test(end) ? "Now" : (end.match(/\d{4}/)?.[0] ?? end);
  return from === to || !to ? from : `${from} – ${to}`;
}

const eyebrow = "mb-2 font-mono text-xs font-medium uppercase tracking-wider text-muted";

/**
 * The hero's right-hand column: roles, then the profile groups that are not
 * already buttons under the bio. No box, just a hairline running the full
 * height of the hero.
 */
export default function ExperienceCard({ groups, className = "" }: { groups: ProfileGroup[]; className?: string }) {
  const { site, experience } = usePortfolio();

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

      {groups.length > 0 && (
        <div className={experience.length > 0 ? "mt-5" : ""}>
          <h2 className={eyebrow}>{site.profiles_title}</h2>
          <Profiles groups={groups} />
        </div>
      )}
    </div>
  );
}
