import { FileText, Github, Globe } from "lucide-react";
import type { Project } from "../api/types";

const iconLink = "text-muted";

/**
 * Fixed-size tile: every text block has a reserved height and is cut with an
 * ellipsis, so cards in a row are always identical in size regardless of how
 * long the title, summary, tags or stack are.
 */
export default function ProjectCard({ project }: { project: Project }) {
  const { title, summary, tags, stack, links } = project;
  return (
    <article className="flex h-[196px] flex-col rounded-xl border border-line bg-card px-4 py-3.5">
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

      {/* Three lines reserved; longer summaries end in an ellipsis. */}
      <p className="mt-1.5 h-[54px] overflow-hidden text-[13px] leading-[18px] text-muted [display:-webkit-box] [-webkit-box-orient:vertical] [-webkit-line-clamp:3]">
        {summary}
      </p>

      {/* One row of tags; overflow is hidden rather than wrapping. */}
      <ul className="mt-auto flex h-5 gap-1 overflow-hidden whitespace-nowrap">
        {tags.map((t) => (
          <li key={t} className="shrink-0 rounded-full border border-line px-2 py-px text-[10.5px] text-muted">
            {t}
          </li>
        ))}
      </ul>
      <p className="mt-2 h-4 truncate font-mono text-[10.5px] text-muted">{stack.join(" · ")}</p>
    </article>
  );
}
