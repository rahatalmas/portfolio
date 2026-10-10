import { FileText, Github, Globe } from "lucide-react";
import type { Project } from "../api/types";

const iconLink = "text-muted hover:text-accent";

export default function ProjectCard({ project }: { project: Project }) {
  const { title, summary, stack, links } = project;
  return (
    <article className="flex h-full flex-col rounded-xl border border-line bg-card px-3.5 py-3 hover:border-accent">
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

      {stack.length > 0 && <p className="mt-2 truncate font-mono text-[10.5px] text-muted">{stack.join(" · ")}</p>}
    </article>
  );
}
