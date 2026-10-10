import { ArrowUp, Mail } from "lucide-react";
import { usePortfolio } from "../api/PortfolioContext";
import { ICONS } from "./Profiles";

const iconLink = "text-muted hover:text-accent";

/** Icon keys shown in the footer, in order, when a matching profile exists. */
const FOOTER_ICONS = ["github", "linkedin", "scholar"];

export default function Footer() {
  const { site, profiles } = usePortfolio();
  const all = profiles.flatMap((g) => g.items);
  const links = FOOTER_ICONS.map((icon) => all.find((p) => p.icon === icon)).filter((p) => p !== undefined);

  return (
    <footer className="mt-8 border-t border-line">
      <div className="mx-auto flex max-w-[960px] flex-col gap-4 px-6 py-6 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p className="text-sm font-semibold tracking-tight">{site.name}</p>
          <p className="font-mono text-xs text-muted">
            © {new Date().getFullYear()}
            {site.footer_updated && ` · Updated ${site.footer_updated}`}
          </p>
        </div>

        <div className="flex items-center gap-4">
          {links.map((p) => {
            const Icon = ICONS[p.icon];
            return (
              <a key={p.id} href={p.href} target="_blank" rel="noreferrer" aria-label={p.label} title={p.label} className={iconLink}>
                <Icon size={17} />
              </a>
            );
          })}
          {site.email && (
            <a href={`mailto:${site.email}`} aria-label="Email" title="Email" className={iconLink}>
              <Mail size={17} />
            </a>
          )}
          <span className="h-4 w-px bg-line" aria-hidden />
          <button
            onClick={() => window.scrollTo({ top: 0 })}
            aria-label={site.footer_back_to_top}
            title={site.footer_back_to_top}
            className={iconLink}
          >
            <ArrowUp size={17} />
          </button>
        </div>
      </div>
    </footer>
  );
}
