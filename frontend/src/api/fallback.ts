import type { Portfolio } from "./types";

/**
 * Built-in content used when the API is unreachable or returns nothing, so
 * the site always renders instead of showing an error. It mirrors the
 * backend's defaults: settings from backend/internal/schema/settings.go and
 * demo content from backend/internal/store/seed.go. Keep the three in sync
 * when you change the placeholder text.
 */
export const FALLBACK: Portfolio = {
  site: {
    // identity
    name: "Rahat Almas",
    handle: "rahat",
    initials: "RA",
    role: "Software Engineer",
    bio: "I'm a software engineer and researcher working at the intersection of machine learning and distributed systems. I've shipped production backends serving millions of requests and published work on efficient model serving. I care about systems that are fast, measurable, and simple enough to reason about.",
    photo: "/profile.jpg",
    location: "City, Country",
    email: "you@example.com",
    github_url: "https://github.com/yourhandle",
    cv_url: "",
    // seo
    meta_title: "Rahat Almas — Software Engineer",
    meta_description: "Portfolio of Rahat Almas: AI and backend engineer, researcher.",
    // sections
    show_experience: true,
    show_about: true,
    show_work: true,
    show_research: true,
    show_contact: true,
    show_contact_form: true,
    // navbar
    nav_about: "About",
    nav_work: "Work",
    nav_research: "Research",
    nav_contact: "Contact",
    nav_cv_label: "CV",
    // hero
    hero_cv_label: "CV",
    hero_github_label: "GitHub",
    hero_email_label: "Email",
    experience_title: "Experience",
    stat_years_label: "Years",
    stat_papers_label: "Papers",
    stat_projects_label: "Projects",
    full_cv_label: "Full CV",
    // about
    about_title: "About",
    education_label: "Education",
    skills_label: "Skills",
    // work
    work_title: "Selected work",
    work_all_label: "All projects",
    work_all_url: "https://github.com/yourhandle?tab=repositories",
    // research
    research_title: "Research",
    research_statement: "",
    interests_label: "Interests",
    interests: ["LLM inference systems", "Retrieval-augmented generation", "Distributed storage", "ML systems reliability"],
    // contact
    contact_title: "Contact",
    contact_email_label: "Email",
    contact_location_label: "Location",
    form_title: "Send a message",
    form_subtitle: "Questions, collaborations, or opportunities.",
    form_name_placeholder: "Name",
    form_email_placeholder: "Email",
    form_message_placeholder: "Message",
    form_submit_label: "Send message",
    form_sending_label: "Sending…",
    form_success_title: "Message sent",
    form_success_text: "Thanks — I'll get back to you soon.",
    form_error_text: "Couldn't send. Email me instead.",
    // footer
    footer_updated: "October 2026",
    footer_back_to_top: "Back to top",
  },
  education: [
    {
      id: 1,
      degree: "B.Sc. Computer Science",
      school: "University Name",
      years: "2020 – 2024",
      note: "GPA 3.9/4.0 · Thesis: Efficient retrieval for long-context LLMs",
    },
  ],
  skills: [
    { id: 1, name: "AI", items: ["PyTorch", "Transformers", "vLLM", "RAG", "LangChain", "CUDA basics"] },
    { id: 2, name: "Backend", items: ["Python", "Go", "FastAPI", "PostgreSQL", "Redis", "gRPC", "Kafka"] },
    { id: 3, name: "Tools", items: ["Docker", "Kubernetes", "AWS", "GitHub Actions", "Linux"] },
  ],
  profiles: [
    {
      name: "Professional",
      items: [
        { id: 1, label: "GitHub", href: "https://github.com/yourhandle", icon: "github" },
        { id: 2, label: "LinkedIn", href: "https://linkedin.com/in/yourhandle", icon: "linkedin" },
        { id: 3, label: "X", href: "https://x.com/yourhandle", icon: "x" },
      ],
    },
    {
      name: "Academic",
      items: [
        { id: 4, label: "Google Scholar", href: "https://scholar.google.com/citations?user=XXXX", icon: "scholar" },
        { id: 5, label: "ORCID", href: "https://orcid.org/0000-0000-0000-0000", icon: "orcid" },
        { id: 6, label: "ResearchGate", href: "https://researchgate.net/profile/yourhandle", icon: "research" },
      ],
    },
    {
      name: "Coding",
      items: [
        { id: 7, label: "LeetCode", href: "https://leetcode.com/u/yourhandle", icon: "code" },
        { id: 8, label: "Codeforces", href: "https://codeforces.com/profile/yourhandle", icon: "trophy" },
        { id: 9, label: "CodeChef", href: "https://codechef.com/users/yourhandle", icon: "chef" },
        { id: 10, label: "HackerRank", href: "https://hackerrank.com/profile/yourhandle", icon: "terminal" },
        { id: 11, label: "Kaggle", href: "https://kaggle.com/yourhandle", icon: "chart" },
      ],
    },
  ],
  experience: [
    { id: 1, title: "Software Engineer", company: "Company A · ML Platform", period: "Jun 2024 – Present" },
    { id: 2, title: "Research Assistant", company: "Systems Lab, University Name", period: "Sep 2023 – May 2024" },
  ],
  projects: [
    {
      id: 1,
      slug: "fast-rag",
      title: "FastRAG",
      summary: "Retrieval pipeline that cut p95 query latency 40% with hybrid BM25 + dense search.",
      tags: ["AI", "Backend"],
      stack: ["Python", "FastAPI", "pgvector", "Redis"],
      featured: true,
      links: { github: "https://github.com/yourhandle/fast-rag", live: "https://fast-rag.example.com" },
    },
    {
      id: 2,
      slug: "kv-cache-scheduler",
      title: "KV-Cache Scheduler",
      summary: "Request scheduler for LLM serving that raised throughput 1.8× at equal latency.",
      tags: ["AI", "Research"],
      stack: ["Python", "PyTorch", "vLLM", "CUDA"],
      featured: true,
      links: { github: "https://github.com/yourhandle/kv-scheduler", live: "https://kv-scheduler.example.com" },
    },
    {
      id: 3,
      slug: "event-ledger",
      title: "Event Ledger",
      summary: "Append-only event store in Go handling 50k writes/s with exactly-once consumers.",
      tags: ["Backend"],
      stack: ["Go", "PostgreSQL", "Kafka", "gRPC"],
      featured: true,
      links: { github: "https://github.com/yourhandle/event-ledger", live: "https://event-ledger.example.com" },
    },
    {
      id: 4,
      slug: "ml-feature-store",
      title: "Feature Store",
      summary: "Online/offline feature store that removed training–serving skew for three models.",
      tags: ["AI", "Backend"],
      stack: ["Python", "Spark", "Redis", "Airflow"],
      featured: true,
      links: { github: "https://github.com/yourhandle/feature-store", live: "https://feature-store.example.com" },
    },
    {
      id: 5,
      slug: "cli-profiler",
      title: "pyprof",
      summary: "Sampling profiler CLI for Python async code with flame-graph output.",
      tags: ["Backend"],
      stack: ["Python", "Rust"],
      featured: false,
      links: { github: "https://github.com/yourhandle/pyprof", live: "https://pyprof.example.com" },
    },
  ],
  papers: [
    {
      id: 1,
      title: "Length-Aware Scheduling for Memory-Efficient LLM Inference",
      authors: ["Rahat Almas", "Co Author", "Advisor Name"],
      venue: "Workshop on ML Systems, NeurIPS",
      year: 2025,
      status: "published",
      links: { pdf: "https://arxiv.org/abs/0000.00000", code: "https://github.com/yourhandle/kv-scheduler" },
    },
    {
      id: 2,
      title: "Hybrid Retrieval for Long-Context Question Answering",
      authors: ["Rahat Almas", "Advisor Name"],
      venue: "arXiv preprint",
      year: 2025,
      status: "preprint",
      links: { pdf: "https://arxiv.org/abs/0000.00001" },
    },
    {
      id: 3,
      title: "Efficient Retrieval for Long-Context LLMs",
      authors: ["Rahat Almas"],
      venue: "B.Sc. Thesis, University Name",
      year: 2024,
      status: "published",
      links: {},
    },
  ],
};

/**
 * Fill gaps in an API response from the fallback so a partial or empty
 * payload never blanks the page. Only missing keys are filled: an empty
 * string or empty list from the API is a real choice made in the admin.
 */
export function withFallback(data: Partial<Portfolio> | null | undefined): Portfolio {
  return {
    site: { ...FALLBACK.site, ...(data?.site ?? {}) },
    education: data?.education ?? FALLBACK.education,
    skills: data?.skills ?? FALLBACK.skills,
    profiles: data?.profiles ?? FALLBACK.profiles,
    experience: data?.experience ?? FALLBACK.experience,
    projects: data?.projects ?? FALLBACK.projects,
    papers: data?.papers ?? FALLBACK.papers,
  };
}
