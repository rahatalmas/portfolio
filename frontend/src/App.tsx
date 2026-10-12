import { Suspense, lazy, useEffect } from "react";
import { PortfolioProvider, usePortfolioState } from "./api/PortfolioContext";
import Nav from "./components/Nav";
import Footer from "./components/Footer";
import Home from "./pages/Home";
// Lazy: the Markdown renderer is ~50 KB gzipped and only project pages need
// it, so the home page never downloads it.
const ProjectDetail = lazy(() => import("./pages/ProjectDetail"));
import NotFound from "./pages/NotFound";

// Routing is a path match, not a router: the site is one page at "/" with
// in-page anchors, plus a page per project that has Markdown details. Links
// are plain <a href>, so the server decides the status for every URL (200 for
// these, 404 for anything else) and the page it serves always agrees.
const projectSlug = (p: string) => {
  const m = /^\/projects\/([^/]+)\/?$/.exec(p);
  return m ? decodeURIComponent(m[1]) : null;
};

export default function App() {
  return (
    <PortfolioProvider>
      <Shell />
    </PortfolioProvider>
  );
}

/**
 * The page shell. Content is always available (built-in fallback until the
 * API answers), so there is no loading or error screen; the first paint is
 * the real layout and it updates in place when the API responds.
 */
function Shell() {
  const { data, loading } = usePortfolioState();
  const path = window.location.pathname;
  const home = path === "/" || path === "/index.html";
  const slug = projectSlug(path);
  // Every project has a page; the Markdown body is optional extra detail.
  const project = slug ? data.projects.find((p) => p.slug === slug) : undefined;
  const notFound = !home && !project;

  // Keep the browser tab and search snippet in sync with the admin settings.
  useEffect(() => {
    if (loading) return;
    const { meta_title, meta_description } = data.site;
    if (notFound) document.title = `Page not found · ${data.site.name}`;
    else if (project) document.title = `${project.title} · ${data.site.name}`;
    else if (meta_title) document.title = meta_title;
    const meta = document.querySelector<HTMLMetaElement>('meta[name="description"]');
    if (meta && meta_description) meta.content = meta_description;
  }, [data, loading, notFound, project]);

  if (notFound) {
    // No Nav: its links are anchors to sections this page does not have.
    return (
      <>
        <main>
          <NotFound />
        </main>
        <Footer />
      </>
    );
  }

  return (
    <>
      <Nav />
      <main>
        {project ? (
          <Suspense fallback={<div className="min-h-[60vh]" />}>
            <ProjectDetail project={project} />
          </Suspense>
        ) : (
          <Home />
        )}
      </main>
      <Footer />
    </>
  );
}
