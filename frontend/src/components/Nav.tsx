import { useEffect, useMemo, useState } from "react";
import { Menu, Moon, Sun, X } from "lucide-react";
import { usePortfolio } from "../api/PortfolioContext";
import { useTheme } from "../hooks/useTheme";
import { useScrollSpy } from "../hooks/useScrollSpy";

export default function Nav() {
  const { site } = usePortfolio();
  const { theme, toggle } = useTheme();
  const [scrolled, setScrolled] = useState(false);
  const [open, setOpen] = useState(false);

  // Only visible sections get a link; labels come from settings.
  const links = useMemo(
    () =>
      [
        { id: "about", label: site.nav_about, show: site.show_about },
        { id: "work", label: site.nav_work, show: site.show_work },
        { id: "research", label: site.nav_research, show: site.show_research },
        { id: "contact", label: site.nav_contact, show: site.show_contact },
      ].filter((l) => l.show),
    [site],
  );
  const ids = useMemo(() => links.map((l) => l.id), [links]);
  const active = useScrollSpy(ids);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 8);
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  const linkClass = (id: string) =>
    `relative py-0.5 text-[13px] hover:text-fg ${active === id ? "text-fg" : "text-muted"} after:absolute after:inset-x-0 after:-bottom-0.5 after:h-px after:bg-accent ${
      active === id ? "after:scale-x-100" : "after:scale-x-0"
    }`;

  return (
    <header
      className={`fixed inset-x-0 top-0 z-50 h-12 ${
        scrolled || open ? "border-b border-line bg-bg/80 backdrop-blur" : "border-b border-transparent bg-transparent"
      }`}
    >
      <nav className="mx-auto flex h-full max-w-[960px] items-center justify-between px-6">
        {/* Monogram: the full name already leads the hero */}
        <a
          href="#top"
          aria-label={`${site.name} — back to top`}
          title={site.name}
          className="font-mono text-sm font-semibold tracking-tight"
        >
          <span className="text-accent">&lt;</span>
          {site.handle}
          <span className="text-accent"> /&gt;</span>
        </a>

        {/* Desktop links */}
        <div className="hidden items-center gap-4 md:flex">
          {links.map((l) => (
            <a key={l.id} href={`#${l.id}`} className={linkClass(l.id)}>
              {l.label}
            </a>
          ))}
          {site.cv_url && (
            <a
              href={site.cv_url}
              target="_blank"
              rel="noreferrer"
              className="rounded-md border border-line px-2.5 py-0.5 text-[13px] hover:border-accent hover:text-accent"
            >
              {site.nav_cv_label}
            </a>
          )}
          <button
            onClick={toggle}
            aria-label={`Switch to ${theme === "dark" ? "light" : "dark"} mode`}
            className="rounded-md p-0.5 text-muted hover:text-fg"
          >
            {theme === "dark" ? <Sun size={16} /> : <Moon size={16} />}
          </button>
        </div>

        {/* Mobile controls */}
        <div className="flex items-center gap-3 md:hidden">
          {site.cv_url && (
            <a
              href={site.cv_url}
              target="_blank"
              rel="noreferrer"
              className="rounded-md border border-line px-2.5 py-0.5 text-[13px]"
            >
              {site.nav_cv_label}
            </a>
          )}
          <button
            onClick={() => setOpen((o) => !o)}
            aria-label={open ? "Close menu" : "Open menu"}
            aria-expanded={open}
            className="rounded-md p-1"
          >
            {open ? <X size={20} /> : <Menu size={20} />}
          </button>
        </div>
      </nav>

      {/* Mobile sheet */}
      {open && (
        <div className="fixed inset-x-0 top-12 bottom-0 flex flex-col bg-bg px-6 pt-8 md:hidden">
          {links.map((l) => (
            <a
              key={l.id}
              href={`#${l.id}`}
              onClick={() => setOpen(false)}
              className="border-b border-line py-4 text-2xl font-medium"
            >
              {l.label}
            </a>
          ))}
          <button onClick={toggle} className="mt-8 inline-flex items-center gap-2 self-start text-muted">
            {theme === "dark" ? <Sun size={18} /> : <Moon size={18} />}
            {theme === "dark" ? "Light mode" : "Dark mode"}
          </button>
        </div>
      )}
    </header>
  );
}
