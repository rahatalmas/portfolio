import { FileText, Github, Globe } from "lucide-react";
import type { Project } from "../api/types";

const iconLink = "text-muted hover:text-accent";

export default function ProjectCard({ project }: { project: Project }) {
  const { title, summary, tags, stack, links } = project;
  return (
    <article className="flex h-full flex-col rounded-xl border border-line bg-card px-3.5 py-3 transition hover:-translate-y-0.5 hover:border-accent hover:shadow-[0_12px_30px_-14px_var(--color-accent)]">
      <div className="flex items-center justify-between gap-2">
        <h3 className="truncate text-sm font-semibold">{title}</h3>
        <div className="flex shrink-0 gap-2">
          {links.github && (
            <a href={links.github} target="_blank" rel="noreferrer" aria-label="Code" title="Code" className={iconLink}>
              <Github size={14} />
            </a>
          )}
          {links.live && (
            <a href={links.live} target="_blank" rel="noreferrer" aria-label="Live site" title="Live site" className={iconLink}>
              <Globe size={14} />
            </a>
          )}
          {links.paper && (
            <a href={links.paper} aria-label="Paper" title="Paper" className={iconLink}>
              <FileText size={14} />
            </a>
          )}
        </div>
      </div>

      <p className="mt-1 line-clamp-3 flex-1 text-[13px] leading-snug text-muted">{summary}</p>

      {tags.length > 0 && (
        <ul className="mt-2.5 flex flex-wrap gap-1">
          {tags.map((t) => (
            <li key={t} className="rounded-full bg-accent-soft px-2 py-px text-[10.5px] font-medium text-accent">
              {t}
            </li>
          ))}
        </ul>
      )}
      {stack.length > 0 && <p className="mt-2 truncate font-mono text-[10.5px] text-muted">{stack.join(" · ")}</p>}
    </article>
  );
}
