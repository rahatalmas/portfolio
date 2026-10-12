import { ArrowUpRight, FileText, Mail } from "lucide-react";
import { usePortfolio } from "../api/PortfolioContext";
import Portrait from "./Portrait";
import ExperienceCard from "./ExperienceCard";
import { ICONS } from "./Profiles";

const primary = "inline-flex items-center gap-2 rounded-md bg-accent px-3.5 py-1.5 text-sm font-medium text-on-accent";
const secondary = "inline-flex items-center gap-2 rounded-md border border-line bg-card px-3.5 py-1.5 text-sm";

export default function Hero() {
  const { site, experience, profiles } = usePortfolio();
  // The first profile group (e.g. Professional: GitHub, LinkedIn, X) becomes
  // the button row; the remaining groups live in the column on the right.
  const [firstGroup, ...rest] = profiles;
  const showColumn = site.show_experience && (experience.length > 0 || rest.length > 0);
  // One filled button: the CV when there is one, otherwise the first link.
  let filled = Boolean(site.cv_url);
  const buttonClass = () => {
    if (!filled) {
      filled = true;
      return primary;
    }
    return secondary;
  };

  return (
    <section id="top">
      <div
        className={`mx-auto grid w-full max-w-[960px] items-stretch gap-8 px-6 pt-20 pb-10 md:gap-10 md:pt-24 md:pb-12 ${
          showColumn ? "md:grid-cols-[1fr_260px]" : ""
        }`}
      >
        <div>
          <div className="flex items-center gap-4">
            <Portrait className="h-16 w-16 shrink-0 md:h-[72px] md:w-[72px]" />
            <div>
              <h1 className="text-[28px] font-bold leading-tight tracking-tight md:text-4xl">{site.name}</h1>
              <p className="mt-0.5 font-mono text-[13px] text-muted">{site.role}</p>
            </div>
          </div>
          <p className="mt-4 leading-relaxed text-muted">{site.bio}</p>

          {site.interests.length > 0 && (
            <p className="mt-3 text-[13px] leading-relaxed">
              {site.interests_label && <span className="text-muted">{site.interests_label} — </span>}
              {site.interests.join(" · ")}
            </p>
          )}

          <div className="mt-5 flex flex-wrap gap-2">
            {site.cv_url && (
              <a href={site.cv_url} target="_blank" rel="noreferrer" className={primary}>
                <FileText size={15} /> {site.hero_cv_label}
              </a>
            )}
            {firstGroup?.items.map((p) => {
              const Icon = ICONS[p.icon] ?? ArrowUpRight;
              return (
                <a key={p.id} href={p.href} target="_blank" rel="noreferrer" className={buttonClass()}>
                  <Icon size={15} /> {p.label}
                </a>
              );
            })}
            {site.email && (
              <a href={`mailto:${site.email}`} className={buttonClass()}>
                <Mail size={15} /> {site.hero_email_label}
              </a>
            )}
          </div>
        </div>

        {showColumn && <ExperienceCard groups={rest} />}
      </div>
      {/* Divider kept to the content width, like every other rule on the page. */}
      <div className="mx-auto max-w-[960px] px-6">
        <div className="border-b border-line" />
      </div>
    </section>
  );
}
