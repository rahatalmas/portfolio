import { ArrowUp, Mail, MapPin } from "lucide-react";
import { usePortfolio } from "../api/PortfolioContext";
import { ICONS } from "./Profiles";

const iconLink = "text-muted";

/** Icon keys shown in the footer, in order, when a matching profile exists. */
const FOOTER_ICONS = ["github", "linkedin", "scholar"];

/** Name, email and location on the left; social icons and back-to-top on the right. */
export default function Footer() {
  const { site, profiles } = usePortfolio();
  const all = profiles.flatMap((g) => g.items);
  const links = FOOTER_ICONS.map((icon) =>
    all.find((p) => p.icon === icon),
  ).filter((p) => p !== undefined);

  return (
    <footer className="mt-8">
      <div className="mx-auto max-w-[960px] px-6">
        <div className="flex flex-col gap-5 border-t border-line py-8 sm:flex-row sm:items-end sm:justify-between">
          <div className="space-y-1.5">
            <p className="text-sm font-semibold tracking-tight">{site.name}</p>
            {(site.email || site.location) && (
              <dl className="flex flex-col gap-1 text-sm sm:flex-row sm:flex-wrap sm:gap-x-5">
                {site.email && (
                  <div className="flex items-center gap-1.5">
                    <dt className="sr-only">{site.contact_email_label}</dt>
                    <Mail size={14} className="text-muted" aria-hidden />
                    <dd>
                      <a
                        href={`mailto:${site.email}`}
                        className="font-medium text-accent"
                      >
                        {site.email}
                      </a>
                    </dd>
                  </div>
                )}
                {site.location && (
                  <div className="flex items-center gap-1.5">
                    <dt className="sr-only">{site.contact_location_label}</dt>
                    <MapPin size={14} className="text-muted" aria-hidden />
                    <dd className="text-muted">{site.location}</dd>
                  </div>
                )}
              </dl>
            )}
            <p className="font-mono text-xs text-muted">
              © {new Date().getFullYear()}
              {site.footer_updated && ` · Updated ${site.footer_updated}`}
            </p>
          </div>

          <div className="flex items-center gap-4">
            {links.map((p) => {
              const Icon = ICONS[p.icon];
              return (
                <a
                  key={p.id}
                  href={p.href}
                  target="_blank"
                  rel="noreferrer"
                  aria-label={p.label}
                  title={p.label}
                  className={iconLink}
                >
                  <Icon size={17} />
                </a>
              );
            })}
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
      </div>
    </footer>
  );
}
