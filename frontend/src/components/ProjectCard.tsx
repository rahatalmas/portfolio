import { useState } from "react";
import { ArrowRight, FileText, Github, Globe } from "lucide-react";
import type { Project } from "../api/types";

const iconLink = "text-muted";

/** Up to three letters from the title, for cards without an image. */
function initials(title: string) {
  return title
    .split(/[\s-]+/)
    .filter(Boolean)
    .map((w) => w[0])
    .join("")
    .slice(0, 3)
    .toUpperCase();
}

/**
 * Image on top, text below. Every block has a fixed height — the image slot is
 * 144px whatever the column width or whether the project even has an image,
 * and each text block reserves its lines — so every card is exactly the same
 * size, in any grid, at any breakpoint. Images are cropped to fill (never
 * stretched), so their own dimensions cannot change the layout.
 */
export default function ProjectCard({ project }: { project: Project }) {
  const { slug, title, summary, image, tags, stack, links } = project;
  const [failed, setFailed] = useState(false);
  const href = `/projects/${slug}`;
  const art = image && !failed ? (
    <img
      src={image}
      alt={`${title} preview`}
      loading="lazy"
      decoding="async"
      onError={() => setFailed(true)}
      className="h-full w-full object-cover"
    />
  ) : (
    <div className="flex h-full w-full items-center justify-center font-mono text-lg tracking-wide text-muted">
      {initials(title)}
    </div>
  );

  return (
    <article className="flex flex-col overflow-hidden rounded-xl border border-line bg-card">
      <div className="h-36 shrink-0 border-b border-line bg-bg">
        <a href={href} aria-label={`${title} details`} className="block h-full w-full">
          {art}
        </a>
      </div>

      <div className="flex flex-1 flex-col px-4 py-3.5">
        <div className="flex items-center justify-between gap-2">
          <h3 className="truncate text-sm font-semibold">
            <a href={href}>{title}</a>
          </h3>
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

        {/* Two lines reserved; longer summaries end in an ellipsis. */}
        <p className="mt-1.5 h-9 overflow-hidden text-[13px] leading-[18px] text-muted [display:-webkit-box] [-webkit-box-orient:vertical] [-webkit-line-clamp:2]">
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
        <div className="mt-2 flex h-4 items-center justify-between gap-2">
          <p className="truncate font-mono text-[10.5px] text-muted">{stack.join(" · ")}</p>
          <a href={href} className="inline-flex shrink-0 items-center gap-0.5 text-[10.5px] text-accent">
            Details <ArrowRight size={11} aria-hidden />
          </a>
        </div>
      </div>
    </article>
  );
}
