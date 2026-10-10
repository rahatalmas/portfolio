import { FileText, Github, Mail } from "lucide-react";
import { usePortfolio } from "../api/PortfolioContext";
import Portrait from "./Portrait";
import ExperienceCard from "./ExperienceCard";

const primary = "inline-flex items-center gap-2 rounded-md bg-accent px-3.5 py-1.5 text-sm font-medium text-white";
const secondary = "inline-flex items-center gap-2 rounded-md border border-line bg-card px-3.5 py-1.5 text-sm";

export default function Hero() {
  const { site, experience } = usePortfolio();
  const showCard = site.show_experience && experience.length > 0;
  // One filled button: the CV when there is one, otherwise GitHub.
  const githubPrimary = !site.cv_url;

  return (
    <section id="top">
      <div
        className={`mx-auto grid w-full max-w-[960px] items-stretch gap-8 px-6 pt-20 pb-6 md:gap-10 md:pt-24 ${
          showCard ? "md:grid-cols-[1fr_280px]" : ""
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

          <div className="mt-5 flex flex-wrap gap-2">
            {site.cv_url && (
              <a href={site.cv_url} target="_blank" rel="noreferrer" className={primary}>
                <FileText size={15} /> {site.hero_cv_label}
              </a>
            )}
            {site.github_url && (
              <a href={site.github_url} target="_blank" rel="noreferrer" className={githubPrimary ? primary : secondary}>
                <Github size={15} /> {site.hero_github_label}
              </a>
            )}
            {site.email && (
              <a href={`mailto:${site.email}`} className={secondary}>
                <Mail size={15} /> {site.hero_email_label}
              </a>
            )}
          </div>
        </div>

        {showCard && <ExperienceCard />}
      </div>
    </section>
  );
}
