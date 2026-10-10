import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from "react";
import { fetchPortfolio } from "./client";
import type { Portfolio } from "./types";

type State =
  | { status: "loading" }
  | { status: "error"; message: string }
  | { status: "ready"; data: Portfolio };

interface ContextValue {
  state: State;
  reload: () => void;
}

const Ctx = createContext<ContextValue | null>(null);

/** Loads the portfolio once and shares it with every component. */
export function PortfolioProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<State>({ status: "loading" });

  const reload = useCallback(() => {
    setState({ status: "loading" });
    fetchPortfolio()
      .then((data) => setState({ status: "ready", data }))
      .catch((err: unknown) =>
        setState({ status: "error", message: err instanceof Error ? err.message : "Unknown error" }),
      );
  }, []);

  useEffect(reload, [reload]);

  return <Ctx.Provider value={{ state, reload }}>{children}</Ctx.Provider>;
}

/** Raw loading/error/ready state, for the app shell. */
export function usePortfolioState(): ContextValue {
  const v = useContext(Ctx);
  if (!v) throw new Error("usePortfolioState must be used inside PortfolioProvider");
  return v;
}

/**
 * The loaded portfolio. Only call from components rendered after the shell
 * has confirmed the data is ready.
 */
export function usePortfolio(): Portfolio {
  const { state } = usePortfolioState();
  if (state.status !== "ready") throw new Error("portfolio is not loaded yet");
  return state.data;
}
