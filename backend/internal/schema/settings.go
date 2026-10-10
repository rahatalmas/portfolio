package schema

// Group is one fieldset on the admin Settings page.
type Group struct {
	Key         string
	Title       string
	Description string
	Fields      []Field
}

// SettingsGroups lists every single-value setting the site uses, grouped the
// way the admin shows them. Keys are the `settings.key` column and the JSON
// property names of GET /api/v1/site.
//
// Add a setting: append a Field here (with a Default). It is seeded into the
// database on the next start, editable in the admin, and returned by the API.
// Mirror the default in frontend/src/api/fallback.ts, which the site renders
// when the API is unavailable.
var SettingsGroups = []Group{
	{
		Key: "identity", Title: "Identity", Description: "Who the site is about.",
		Fields: []Field{
			{Key: "name", Label: "Full name", Kind: Text, Required: true, Default: "Rahat Almas"},
			{Key: "handle", Label: "Navbar handle", Kind: Text, Required: true, Default: "rahat", Help: "Rendered as <handle /> in the navbar"},
			{Key: "initials", Label: "Initials", Kind: Text, Required: true, Default: "RA", Help: "Fallback when the photo is missing"},
			{Key: "role", Label: "Role / title", Kind: Text, Required: true, Default: "Software Engineer"},
			{Key: "bio", Label: "Bio", Kind: Textarea, Required: true, Default: "I'm a software engineer and researcher working at the intersection of machine learning and distributed systems. I've shipped production backends serving millions of requests and published work on efficient model serving. I care about systems that are fast, measurable, and simple enough to reason about."},
			{Key: "photo", Label: "Photo", Kind: Media, Default: "/profile.jpg", Help: "Upload under Media, then pick it here. Portrait crop works best. The default is the placeholder shipped with the frontend"},
			{Key: "location", Label: "Location", Kind: Text, Default: "City, Country"},
			{Key: "email", Label: "Email", Kind: Email, Required: true, Default: "you@example.com"},
			{Key: "github_url", Label: "GitHub profile URL", Kind: URL, Default: "https://github.com/yourhandle"},
			{Key: "cv_url", Label: "CV file", Kind: Media, Default: "/cv.pdf", Help: "Upload a PDF under Media, then pick it here"},
		},
	},
	{
		Key: "seo", Title: "SEO", Description: "Browser tab and search snippet.",
		Fields: []Field{
			{Key: "meta_title", Label: "Page title", Kind: Text, Required: true, Default: "Rahat Almas — Software Engineer"},
			{Key: "meta_description", Label: "Meta description", Kind: Textarea, Default: "Portfolio of Rahat Almas: AI and backend engineer, researcher."},
		},
	},
	{
		Key: "sections", Title: "Sections", Description: "Show or hide whole sections of the page.",
		Fields: []Field{
			{Key: "show_experience", Label: "Show experience card in the hero", Kind: Bool, Default: "1"},
			{Key: "show_about", Label: "Show About", Kind: Bool, Default: "1"},
			{Key: "show_work", Label: "Show Work", Kind: Bool, Default: "1"},
			{Key: "show_research", Label: "Show Research", Kind: Bool, Default: "1"},
			{Key: "show_contact", Label: "Show Contact", Kind: Bool, Default: "1"},
			{Key: "show_contact_form", Label: "Show the message form inside Contact", Kind: Bool, Default: "1"},
		},
	},
	{
		Key: "nav", Title: "Navbar", Description: "Link labels. Hidden sections are dropped from the navbar automatically.",
		Fields: []Field{
			{Key: "nav_about", Label: "About link", Kind: Text, Required: true, Default: "About"},
			{Key: "nav_work", Label: "Work link", Kind: Text, Required: true, Default: "Work"},
			{Key: "nav_research", Label: "Research link", Kind: Text, Required: true, Default: "Research"},
			{Key: "nav_contact", Label: "Contact link", Kind: Text, Required: true, Default: "Contact"},
			{Key: "nav_cv_label", Label: "CV button", Kind: Text, Required: true, Default: "CV"},
		},
	},
	{
		Key: "hero", Title: "Hero", Description: "Buttons under the bio and the experience card.",
		Fields: []Field{
			{Key: "hero_cv_label", Label: "CV button", Kind: Text, Required: true, Default: "CV"},
			{Key: "hero_github_label", Label: "GitHub button", Kind: Text, Required: true, Default: "GitHub"},
			{Key: "hero_email_label", Label: "Email button", Kind: Text, Required: true, Default: "Email"},
			{Key: "experience_title", Label: "Experience card title", Kind: Text, Required: true, Default: "Experience"},
			{Key: "stat_years_label", Label: "Stat: years label", Kind: Text, Required: true, Default: "Years"},
			{Key: "stat_papers_label", Label: "Stat: papers label", Kind: Text, Required: true, Default: "Papers"},
			{Key: "stat_projects_label", Label: "Stat: projects label", Kind: Text, Required: true, Default: "Projects"},
			{Key: "full_cv_label", Label: "Full CV link", Kind: Text, Required: true, Default: "Full CV"},
		},
	},
	{
		Key: "about", Title: "About section",
		Fields: []Field{
			{Key: "about_title", Label: "Section title", Kind: Text, Required: true, Default: "About"},
			{Key: "education_label", Label: "Education label", Kind: Text, Required: true, Default: "Education"},
			{Key: "skills_label", Label: "Skills label", Kind: Text, Required: true, Default: "Skills"},
		},
	},
	{
		Key: "work", Title: "Work section",
		Fields: []Field{
			{Key: "work_title", Label: "Section title", Kind: Text, Required: true, Default: "Selected work"},
			{Key: "work_all_label", Label: "“All projects” link label", Kind: Text, Default: "All projects", Help: "Leave empty to hide the link"},
			{Key: "work_all_url", Label: "“All projects” link URL", Kind: URL, Default: "https://github.com/yourhandle?tab=repositories"},
		},
	},
	{
		Key: "research", Title: "Research section",
		Fields: []Field{
			{Key: "research_title", Label: "Section title", Kind: Text, Required: true, Default: "Research"},
			{Key: "research_statement", Label: "Research statement", Kind: Textarea, Help: "Optional paragraph shown above the papers", Default: ""},
			{Key: "interests_label", Label: "Interests label", Kind: Text, Default: "Interests"},
			{Key: "interests", Label: "Interests", Kind: Lines, Help: "One per line", Default: "LLM inference systems\nRetrieval-augmented generation\nDistributed storage\nML systems reliability"},
		},
	},
	{
		Key: "contact", Title: "Contact section",
		Fields: []Field{
			{Key: "contact_title", Label: "Section title", Kind: Text, Required: true, Default: "Contact"},
			{Key: "contact_email_label", Label: "Email row label", Kind: Text, Required: true, Default: "Email"},
			{Key: "contact_location_label", Label: "Location row label", Kind: Text, Required: true, Default: "Location"},
			{Key: "form_title", Label: "Form title", Kind: Text, Required: true, Default: "Send a message"},
			{Key: "form_subtitle", Label: "Form subtitle", Kind: Text, Default: "Questions, collaborations, or opportunities."},
			{Key: "form_name_placeholder", Label: "Name placeholder", Kind: Text, Required: true, Default: "Name"},
			{Key: "form_email_placeholder", Label: "Email placeholder", Kind: Text, Required: true, Default: "Email"},
			{Key: "form_message_placeholder", Label: "Message placeholder", Kind: Text, Required: true, Default: "Message"},
			{Key: "form_submit_label", Label: "Submit button", Kind: Text, Required: true, Default: "Send message"},
			{Key: "form_sending_label", Label: "Submit button while sending", Kind: Text, Required: true, Default: "Sending…"},
			{Key: "form_success_title", Label: "Success title", Kind: Text, Required: true, Default: "Message sent"},
			{Key: "form_success_text", Label: "Success text", Kind: Text, Required: true, Default: "Thanks — I'll get back to you soon."},
			{Key: "form_error_text", Label: "Error text", Kind: Text, Required: true, Default: "Couldn't send. Email me instead.", Help: "Your email address is appended automatically"},
		},
	},
	{
		Key: "footer", Title: "Footer",
		Fields: []Field{
			{Key: "footer_updated", Label: "“Updated” text", Kind: Text, Default: "October 2026", Help: "Shown as “© year · Updated …”; leave empty to hide"},
			{Key: "footer_back_to_top", Label: "Back-to-top tooltip", Kind: Text, Required: true, Default: "Back to top"},
		},
	},
}

// SettingFields returns every setting field across all groups.
func SettingFields() []Field {
	var out []Field
	for _, g := range SettingsGroups {
		out = append(out, g.Fields...)
	}
	return out
}

// SettingField looks a setting up by key.
func SettingField(key string) (Field, bool) {
	for _, f := range SettingFields() {
		if f.Key == key {
			return f, true
		}
	}
	return Field{}, false
}

// SettingDefaults returns key → default value for seeding.
func SettingDefaults() map[string]string {
	out := make(map[string]string)
	for _, f := range SettingFields() {
		out[f.Key] = f.Default
	}
	return out
}

// TypedSettings converts stored strings into the JSON shape of
// GET /api/v1/site. Unknown keys are passed through as strings.
func TypedSettings(values map[string]string) map[string]any {
	out := make(map[string]any, len(values))
	for _, f := range SettingFields() {
		out[f.Key] = f.Typed(values[f.Key])
	}
	return out
}
