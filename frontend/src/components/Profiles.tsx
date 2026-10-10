import {
  ArrowUpRight,
  AtSign,
  BarChart3,
  ChefHat,
  Code2,
  Fingerprint,
  Github,
  Globe,
  GraduationCap,
  Linkedin,
  Mail,
  Microscope,
  Terminal,
  Trophy,
} from "lucide-react";
import type { ComponentType } from "react";
import { usePortfolio } from "../api/PortfolioContext";

/** Icon keys the admin can pick for a profile link. Keep in sync with backend schema.ProfileIcons. */
export const ICONS: Record<string, ComponentType<{ size?: number }>> = {
  github: Github,
  linkedin: Linkedin,
  x: AtSign,
  scholar: GraduationCap,
  orcid: Fingerprint,
  research: Microscope,
  code: Code2,
  trophy: Trophy,
  chef: ChefHat,
  terminal: Terminal,
  chart: BarChart3,
  mail: Mail,
  globe: Globe,
};

/** Compact grouped list of profile links: icon + name, one row per group. */
export default function Profiles() {
  const { profiles } = usePortfolio();
  return (
    <dl className="space-y-2.5">
      {profiles.map((group) => (
        <div key={group.name} className="flex flex-col gap-1.5 sm:flex-row sm:gap-3">
          <dt className="w-24 shrink-0 pt-0.5 text-sm text-muted">{group.name}</dt>
          <dd className="flex flex-wrap gap-x-5 gap-y-1.5">
            {group.items.map((p) => {
              const Icon = ICONS[p.icon] ?? ArrowUpRight;
              return (
                <a
                  key={p.id}
                  href={p.href}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1.5 text-sm hover:text-accent"
                >
                  <span className="text-muted">
                    <Icon size={15} />
                  </span>
                  {p.label}
                </a>
              );
            })}
          </dd>
        </div>
      ))}
    </dl>
  );
}
