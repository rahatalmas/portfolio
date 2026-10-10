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

  // While the mobile menu is open: lock page scroll, close on Escape, and
  // close automatically if the viewport grows to the desktop layout.
  useEffect(() => {
    if (!open) return;
    const prev = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    const onKey = (e: KeyboardEvent) => e.key === "Escape" && setOpen(false);
    const mq = window.matchMedia("(min-width: 768px)");
    const onResize = () => mq.matches && setOpen(false);
    window.addEventListener("keydown", onKey);
    mq.addEventListener("change", onResize);
    return () => {
      document.body.style.overflow = prev;
      window.removeEventListener("keydown", onKey);
      mq.removeEventListener("change", onResize);
    };
  }, [open]);

  const linkClass = (id: string) =>
    `relative py-0.5 text-[13px] ${active === id ? "text-fg" : "text-muted"} after:absolute after:inset-x-0 after:-bottom-0.5 after:h-px after:bg-accent ${
      active === id ? "after:scale-x-100" : "after:scale-x-0"
    }`;

  return (
    <header className="fixed inset-x-0 top-0 z-50">
      {/*
        The bar carries the blur. It is a separate element from the menu sheet
        because backdrop-filter makes an element the containing block for
        fixed descendants, which would collapse the sheet into the 48px bar.
      */}
      <div
        className={`h-12 ${
          scrolled || open ? "border-b border-line bg-bg/80 backdrop-blur" : "border-b border-transparent bg-transparent"
        }`}
      >
        <nav className="mx-auto flex h-full max-w-[960px] items-center justify-between gap-3 px-4 sm:px-6">
          {/* Monogram: the full name already leads the hero */}
          <a
            href="#top"
            aria-label={`${site.name} — back to top`}
            title={site.name}
            onClick={() => setOpen(false)}
            className="shrink-0 font-mono text-sm font-semibold tracking-tight"
          >
            <span className="text-accent">&lt;</span>
            {site.handle}
            <span className="text-accent"> /&gt;</span>
          </a>

          {/* Desktop links */}
          <div className="hidden min-w-0 items-center gap-4 md:flex">
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
                className="rounded-md border border-line px-2.5 py-0.5 text-[13px]"
              >
                {site.nav_cv_label}
              </a>
            )}
            <button
              onClick={toggle}
              aria-label={`Switch to ${theme === "dark" ? "light" : "dark"} mode`}
              className="rounded-md p-0.5 text-muted"
            >
              {theme === "dark" ? <Sun size={16} /> : <Moon size={16} />}
            </button>
          </div>

          {/* Mobile controls */}
          <div className="flex items-center gap-2 md:hidden">
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
              onClick={toggle}
              aria-label={`Switch to ${theme === "dark" ? "light" : "dark"} mode`}
              className="rounded-md p-1.5 text-muted"
            >
              {theme === "dark" ? <Sun size={18} /> : <Moon size={18} />}
            </button>
            <button
              onClick={() => setOpen((o) => !o)}
              aria-label={open ? "Close menu" : "Open menu"}
              aria-expanded={open}
              aria-controls="mobile-menu"
              className="-mr-1.5 rounded-md p-1.5"
            >
              {open ? <X size={20} /> : <Menu size={20} />}
            </button>
          </div>
        </nav>
      </div>

      {/* Mobile sheet: sibling of the blurred bar so `fixed` is viewport-relative */}
      {open && (
        <div
          id="mobile-menu"
          className="fixed inset-x-0 top-12 bottom-0 flex flex-col overflow-y-auto bg-bg px-6 pt-6 pb-8 md:hidden"
        >
          {links.map((l) => (
            <a
              key={l.id}
              href={`#${l.id}`}
              onClick={() => setOpen(false)}
              className={`border-b border-line py-4 text-2xl font-medium ${active === l.id ? "text-accent" : ""}`}
            >
              {l.label}
            </a>
          ))}
          {links.length === 0 && <p className="py-4 text-muted">Nothing to navigate to yet.</p>}
        </div>
      )}
    </header>
  );
}
