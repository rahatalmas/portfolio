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

/** Switches between loading, error and the real page. */
function Shell() {
  const { state, reload } = usePortfolioState();

  // Keep the browser tab and search snippet in sync with the admin settings.
  useEffect(() => {
    if (state.status !== "ready") return;
    const { meta_title, meta_description } = state.data.site;
    if (meta_title) document.title = meta_title;
    const meta = document.querySelector<HTMLMetaElement>('meta[name="description"]');
    if (meta && meta_description) meta.content = meta_description;
  }, [state]);

  if (state.status === "loading") {
    return (
      <div className="flex min-h-screen items-center justify-center text-sm text-muted" role="status">
        Loading…
      </div>
    );
  }
  if (state.status === "error") {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-3 px-6 text-center">
        <p className="font-medium">Couldn't load the site.</p>
        <p className="text-sm text-muted">{state.message}</p>
        <button
          onClick={reload}
          className="rounded-md border border-line px-3 py-1.5 text-sm hover:border-accent hover:text-accent"
        >
          Try again
        </button>
      </div>
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
