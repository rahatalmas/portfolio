import { useEffect } from "react";
import { PortfolioProvider, usePortfolioState } from "./api/PortfolioContext";
import Nav from "./components/Nav";
import Footer from "./components/Footer";
import Home from "./pages/Home";

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

  // Keep the browser tab and search snippet in sync with the admin settings.
  useEffect(() => {
    if (loading) return;
    const { meta_title, meta_description } = data.site;
    if (meta_title) document.title = meta_title;
    const meta = document.querySelector<HTMLMetaElement>('meta[name="description"]');
    if (meta && meta_description) meta.content = meta_description;
  }, [data, loading]);

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
