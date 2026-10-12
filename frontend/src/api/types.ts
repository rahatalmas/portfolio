/**
 * Shapes returned by the backend API. They mirror the Go types in
 * backend/internal/models and the settings spec in backend/internal/schema.
 * See docs/API.md.
 */

/** Every single-value setting, as returned by GET /api/v1/site. */
export interface Site {
  // identity
  name: string;
  handle: string;
  initials: string;
  role: string;
  bio: string;
  photo: string;
  location: string;
  email: string;
  cv_url: string;
  // seo
  meta_title: string;
  meta_description: string;
  // section visibility
  show_experience: boolean;
  show_about: boolean;
  show_work: boolean;
  show_research: boolean;
  // navbar
  nav_about: string;
  nav_work: string;
  nav_research: string;
  nav_cv_label: string;
  // hero
  hero_cv_label: string;
  hero_email_label: string;
  experience_title: string;
  // about
  about_title: string;
  education_label: string;
  skills_label: string;
  // work
  work_title: string;
  work_all_label: string;
  work_all_url: string;
  // research
  research_title: string;
  research_statement: string;
  interests_label: string;
  interests: string[];
  // footer contact details
  contact_email_label: string;
  contact_location_label: string;
  // footer
  footer_updated: string;
  footer_back_to_top: string;
}

export interface Education {
  id: number;
  degree: string;
  school: string;
  years: string;
  note: string;
}

export interface SkillGroup {
  id: number;
  name: string;
  items: string[];
}

export interface Profile {
  id: number;
  label: string;
  href: string;
  icon: string;
}

export interface ProfileGroup {
  name: string;
  items: Profile[];
}

export interface Role {
  id: number;
  title: string;
  company: string;
  period: string;
}

export interface Project {
  id: number;
  slug: string;
  title: string;
  summary: string;
  /** Illustration for the card: /uploads/… or any URL. Optional. */
  image?: string;
  /** Markdown details. When present the project has a page at /projects/<slug>. */
  body?: string;
  tags: string[];
  stack: string[];
  featured: boolean;
  links: { github?: string; live?: string; paper?: string };
}

export interface Paper {
  id: number;
  title: string;
  authors: string[];
  venue: string;
  year: number;
  status: "published" | "preprint" | "in progress";
  links: { pdf?: string; code?: string; poster?: string; slides?: string };
}

/** GET /api/v1/portfolio — everything the page renders. */
export interface Portfolio {
  site: Site;
  education: Education[];
  skills: SkillGroup[];
  profiles: ProfileGroup[];
  experience: Role[];
  projects: Project[];
  papers: Paper[];
}
