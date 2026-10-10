import type { ReactNode } from "react";

export const sectionTitle = "text-[28px] font-semibold tracking-tight";

/** A short, hand-drawn brush stroke that sits under every section title. */
export function BrushStroke({ className = "" }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 120 12"
      aria-hidden
      className={`mt-1.5 block h-2.5 w-20 text-accent ${className}`}
      preserveAspectRatio="none"
    >
      <path
        fill="currentColor"
        d="M1.5 6.2C14 2.1 27 8.4 41 4.9c15-3.8 31 3.9 47 1.1 10-1.8 20-2.3 30 .1l.6 3.4c-10-1.2-20-.3-30 1.4-16 2.8-32-3.4-47 .8-14 3.9-27-1.2-39.5 2.7Z"
      />
    </svg>
  );
}

interface Props {
  id: string;
  title?: string; // omit to render your own heading inside children
  children: ReactNode;
  action?: ReactNode; // e.g. "All projects →"
}

export default function Section({ id, title, children, action }: Props) {
  return (
    <section id={id} className="scroll-mt-20 py-8 md:py-10">
      <div className="mx-auto max-w-[960px] px-6">
        {title && (
          <div className="mb-6 flex items-baseline justify-between gap-4">
            <div>
              <h2 className={sectionTitle}>{title}</h2>
              <BrushStroke />
            </div>
            {action && <div className="shrink-0 text-sm">{action}</div>}
          </div>
        )}
        {children}
      </div>
    </section>
  );
}
