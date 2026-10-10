import { useEffect, useState } from "react";

/** Returns the id of the section currently nearest the top of the viewport. */
export function useScrollSpy(ids: string[], offset = 120) {
  const [active, setActive] = useState<string>("");

  useEffect(() => {
    const onScroll = () => {
      let current = "";
      for (const id of ids) {
        const el = document.getElementById(id);
        if (el && el.getBoundingClientRect().top - offset <= 0) current = id;
      }
      setActive(current);
    };
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, [ids, offset]);

  return active;
}
