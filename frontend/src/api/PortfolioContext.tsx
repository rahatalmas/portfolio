import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from "react";
import { fetchPortfolio } from "./client";
import { FALLBACK, withFallback } from "./fallback";
import type { Portfolio } from "./types";

interface ContextValue {
  /** Content to render. Starts as the built-in fallback, replaced by the API. */
  data: Portfolio;
  /** True while the first request is in flight. */
  loading: boolean;
  /** True when the API could not be reached and the fallback is shown. */
  offline: boolean;
  reload: () => void;
}

const Ctx = createContext<ContextValue | null>(null);

/**
 * Loads the portfolio once and shares it with every component. The page is
 * never blocked on the API: it renders the fallback immediately and swaps in
 * the real content when it arrives. If the request fails the fallback stays
 * and nothing is shown to the visitor beyond a console warning.
 */
export function PortfolioProvider({ children }: { children: ReactNode }) {
  const [data, setData] = useState<Portfolio>(FALLBACK);
  const [loading, setLoading] = useState(true);
  const [offline, setOffline] = useState(false);

  const reload = useCallback(() => {
    setLoading(true);
    fetchPortfolio()
      .then((p) => {
        setData(withFallback(p));
        setOffline(false);
      })
      .catch((err: unknown) => {
        console.warn("portfolio: API unavailable, showing built-in content", err);
        setOffline(true);
      })
      .finally(() => setLoading(false));
  }, []);

  useEffect(reload, [reload]);

  return <Ctx.Provider value={{ data, loading, offline, reload }}>{children}</Ctx.Provider>;
}

/** Loading/offline flags plus reload, for the app shell. */
export function usePortfolioState(): ContextValue {
  const v = useContext(Ctx);
  if (!v) throw new Error("usePortfolioState must be used inside PortfolioProvider");
  return v;
}

/** The portfolio to render. Always available. */
export function usePortfolio(): Portfolio {
  return usePortfolioState().data;
}
