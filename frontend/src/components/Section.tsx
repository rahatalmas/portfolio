import type { ReactNode } from "react";
import { useReveal } from "../hooks/useReveal";

export const sectionTitle = "text-[28px] font-semibold tracking-tight";

interface Props {
  id: string;
  title?: string; // omit to render your own heading inside children
  children: ReactNode;
  action?: ReactNode; // e.g. "All projects →"
}

export default function Section({ id, title, children, action }: Props) {
  const ref = useReveal<HTMLElement>();
  return (
    <section id={id} ref={ref} className="reveal scroll-mt-20 py-8 md:py-10">
      <div className="mx-auto max-w-[960px] px-6">
        {title && (
          <div className="mb-5 flex items-baseline justify-between gap-4">
            <h2 className={sectionTitle}>{title}</h2>
            {action && <div className="shrink-0 text-sm">{action}</div>}
          </div>
        )}
        {children}
      </div>
    </section>
  );
}
