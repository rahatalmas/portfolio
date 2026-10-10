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
import type { ProfileGroup } from "../api/types";

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

/**
 * Compact grouped profile links for the hero column: the group name as a
 * small label, then its links wrapping underneath.
 */
export default function Profiles({ groups }: { groups: ProfileGroup[] }) {
  return (
    <dl className="space-y-2.5">
      {groups.map((group) => (
        <div key={group.name}>
          <dt className="mb-1 text-[11px] text-muted">{group.name}</dt>
          <dd className="flex flex-wrap gap-x-3.5 gap-y-1">
            {group.items.map((p) => {
              const Icon = ICONS[p.icon] ?? ArrowUpRight;
              return (
                <a key={p.id} href={p.href} target="_blank" rel="noreferrer" className="inline-flex items-center gap-1 text-[13px]">
                  <span className="text-muted">
                    <Icon size={13} />
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
