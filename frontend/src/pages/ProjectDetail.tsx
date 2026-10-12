import { ArrowLeft, FileText, Github, Globe } from "lucide-react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import type { Project } from "../api/types";

const iconLink = "inline-flex items-center gap-1.5 text-[13px] text-muted";

/**
 * One project, long-form. The header repeats what the card showed, then the
 * Markdown body carries everything else — the same shape as a README.
 *
 * Markdown is rendered without raw HTML (react-markdown's default), so a body
 * is safe to render even though it is author-supplied; GFM adds tables,
 * task lists, strikethrough and autolinks. Typography lives in the .md
 * class in index.css, not in a plugin.
 */
export default function ProjectDetail({ project }: { project: Project }) {
  const { title, summary, image, tags, stack, links, body } = project;

  return (
    <article className="mx-auto max-w-[760px] px-6 pt-20 pb-10 md:pt-24">
      <a href="/#work" className="inline-flex items-center gap-1.5 text-[13px] text-muted">
        <ArrowLeft size={14} aria-hidden />
        All projects
      </a>

      <header className="mt-5">
        <h1 className="text-[28px] font-bold leading-tight tracking-tight md:text-4xl">{title}</h1>
        <p className="mt-3 leading-relaxed text-muted">{summary}</p>

        {(tags.length > 0 || stack.length > 0) && (
          <div className="mt-4 flex flex-wrap items-center gap-x-3 gap-y-2">
            <ul className="flex flex-wrap gap-1.5">
              {tags.map((t) => (
                <li key={t} className="rounded-full border border-line px-2 py-px text-[11px] text-muted">
                  {t}
                </li>
              ))}
            </ul>
            {stack.length > 0 && <p className="font-mono text-[11px] text-muted">{stack.join(" · ")}</p>}
          </div>
        )}

        {(links.github || links.live || links.paper) && (
          <div className="mt-4 flex flex-wrap gap-4 border-t border-line pt-4">
            {links.github && (
              <a href={links.github} target="_blank" rel="noreferrer" className={iconLink}>
                <Github size={14} aria-hidden /> Code
              </a>
            )}
            {links.live && (
              <a href={links.live} target="_blank" rel="noreferrer" className={iconLink}>
                <Globe size={14} aria-hidden /> Live
              </a>
            )}
            {links.paper && (
              <a href={links.paper} className={iconLink}>
                <FileText size={14} aria-hidden /> Paper
              </a>
            )}
          </div>
        )}
      </header>

      {image && (
        <img
          src={image}
          alt={`${title} preview`}
          className="mt-8 w-full rounded-xl border border-line bg-card object-cover"
        />
      )}

      {body && (
        <div className="md mt-8">
          <Markdown remarkPlugins={[remarkGfm]}>{body}</Markdown>
        </div>
      )}
    </article>
  );
}
