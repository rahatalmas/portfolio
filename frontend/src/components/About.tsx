import { usePortfolio } from "../api/PortfolioContext";
import Section from "./Section";

function Label({ children }: { children: string }) {
  return <h3 className="mb-3 font-mono text-xs font-medium uppercase tracking-wider text-accent">{children}</h3>;
}

export default function About() {
  const { site, education, skills } = usePortfolio();
  return (
    <Section id="about" title={site.about_title}>
      {education.length > 0 && (
        <div>
          <Label>{site.education_label}</Label>
          <ul className="space-y-4">
            {education.map((e) => (
              <li key={e.id} className="flex flex-col gap-0.5 sm:flex-row sm:justify-between sm:gap-6">
                <div>
                  <span className="font-medium">{e.degree}</span>
                  <span className="text-muted"> · {e.school}</span>
                  {e.note && <div className="text-sm text-muted">{e.note}</div>}
                </div>
                <span className="shrink-0 font-mono text-xs text-muted sm:pt-1">{e.years}</span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {skills.length > 0 && (
        <div className={education.length > 0 ? "mt-8 border-t border-line pt-6" : ""}>
          <Label>{site.skills_label}</Label>
          <dl className="grid gap-6 sm:grid-cols-3">
            {skills.map((g) => (
              <div key={g.id}>
                <dt className="mb-1.5 text-sm font-semibold">{g.name}</dt>
                <dd className="text-sm leading-relaxed text-muted">{g.items.join(" · ")}</dd>
              </div>
            ))}
          </dl>
        </div>
      )}
    </Section>
  );
}
