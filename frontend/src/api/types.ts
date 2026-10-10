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
  github_url: string;
  cv_url: string;
  // seo
  meta_title: string;
  meta_description: string;
  // section visibility
  show_experience: boolean;
  show_about: boolean;
  show_work: boolean;
  show_research: boolean;
  show_contact: boolean;
  show_contact_form: boolean;
  // navbar
  nav_about: string;
  nav_work: string;
  nav_research: string;
  nav_contact: string;
  nav_cv_label: string;
  // hero
  hero_cv_label: string;
  hero_github_label: string;
  hero_email_label: string;
  experience_title: string;
  stat_years_label: string;
  stat_papers_label: string;
  stat_projects_label: string;
  full_cv_label: string;
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
  // contact
  contact_title: string;
  contact_email_label: string;
  contact_location_label: string;
  form_title: string;
  form_subtitle: string;
  form_name_placeholder: string;
  form_email_placeholder: string;
  form_message_placeholder: string;
  form_submit_label: string;
  form_sending_label: string;
  form_success_title: string;
  form_success_text: string;
  form_error_text: string;
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

/** POST /api/v1/messages body. */
export interface MessageInput {
  name: string;
  email: string;
  message: string;
  /** Honeypot — leave empty. */
  website?: string;
}
