import { FileText, Github, Globe } from "lucide-react";
import type { Project } from "../api/types";

/**
 * Dark gradients the cards cycle through, so a row of projects reads as a set
 * of distinct tiles. Text on them is always white, in both themes.
 */
const GRADES = [
  "linear-gradient(135deg, #0f172a 0%, #1e3a8a 100%)", // navy → blue
  "linear-gradient(135deg, #1e1b4b 0%, #5b21b6 100%)", // indigo → violet
  "linear-gradient(135deg, #042f2e 0%, #0f766e 100%)", // deep teal
  "linear-gradient(135deg, #3b0764 0%, #9d174d 100%)", // plum → rose
  "linear-gradient(135deg, #1c1917 0%, #7c2d12 100%)", // charcoal → rust
  "linear-gradient(135deg, #082f49 0%, #0369a1 100%)", // ocean
];

const iconLink = "text-white/70 hover:text-white";

/**
 * Fixed-size tile: every text block has a reserved height and is cut with an
 * ellipsis, so cards in a row are always identical in size regardless of how
 * long the title, summary, tags or stack are.
 */
export default function ProjectCard({ project, index = 0 }: { project: Project; index?: number }) {
  const { title, summary, tags, stack, links } = project;
  return (
    <article
      style={{ backgroundImage: GRADES[index % GRADES.length] }}
      className="flex h-[196px] flex-col rounded-xl px-4 py-3.5 text-white shadow-[0_14px_34px_-18px_rgba(0,0,0,0.6)] transition hover:-translate-y-0.5 hover:shadow-[0_20px_40px_-18px_rgba(0,0,0,0.7)]"
    >
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
      <p className="mt-1.5 h-[54px] overflow-hidden text-[13px] leading-[18px] text-white/75 [display:-webkit-box] [-webkit-box-orient:vertical] [-webkit-line-clamp:3]">
        {summary}
      </p>

      {/* One row of tags; overflow is hidden rather than wrapping. */}
      <ul className="mt-auto flex h-5 gap-1 overflow-hidden whitespace-nowrap">
        {tags.map((t) => (
          <li key={t} className="shrink-0 rounded-full bg-white/15 px-2 py-px text-[10.5px] font-medium">
            {t}
          </li>
        ))}
      </ul>
      <p className="mt-2 h-4 truncate font-mono text-[10.5px] text-white/60">{stack.join(" · ")}</p>
    </article>
  );
}
