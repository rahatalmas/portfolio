import { useEffect } from "react";
import { PortfolioProvider, usePortfolioState } from "./api/PortfolioContext";
import Nav from "./components/Nav";
import Footer from "./components/Footer";
import Home from "./pages/Home";
import NotFound from "./pages/NotFound";

// The site is one page at "/", with in-page anchors for the sections; there is
// no client-side router. Anything else is a 404 (the server answers those with
// a 404 status and this same bundle, so the page and the status agree).
const isKnownPath = (p: string) => p === "/" || p === "/index.html";

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
  const notFound = !isKnownPath(window.location.pathname);

  // Keep the browser tab and search snippet in sync with the admin settings.
  useEffect(() => {
    if (loading) return;
    const { meta_title, meta_description } = data.site;
    if (notFound) document.title = `Page not found · ${data.site.name}`;
    else if (meta_title) document.title = meta_title;
    const meta = document.querySelector<HTMLMetaElement>('meta[name="description"]');
    if (meta && meta_description) meta.content = meta_description;
  }, [data, loading, notFound]);

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
        <Home />
      </main>
      <Footer />
    </>
  );
}
