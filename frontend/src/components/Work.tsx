import { ArrowUpRight } from "lucide-react";
import { usePortfolio } from "../api/PortfolioContext";
import Section from "./Section";
import ProjectCard from "./ProjectCard";

export default function Work() {
  const { site, projects } = usePortfolio();
  const featured = projects.filter((p) => p.featured).slice(0, 4);
  if (featured.length === 0) return null;

  return (
    <Section
      id="work"
      title={site.work_title}
      action={
        site.work_all_label && site.work_all_url ? (
          <a
            href={site.work_all_url}
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-0.5 text-accent hover:underline"
          >
            {site.work_all_label} <ArrowUpRight size={14} />
          </a>
        ) : undefined
      }
    >
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
        {featured.map((p) => (
          <ProjectCard key={p.id} project={p} />
        ))}
      </div>
    </Section>
  );
}
