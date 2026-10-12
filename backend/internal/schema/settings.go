package schema

// Group is one fieldset on the admin Settings page.
type Group struct {
	Key         string
	Title       string
	Description string
	Fields      []Field
}

// SettingsGroups lists every single-value setting the site uses, grouped the
// way the admin shows them. Keys are the keys under "settings" in
// content.json and the JSON property names of GET /api/v1/site.
//
// Add a setting: append a Field here (with a Default). It appears in the
// admin and the API immediately; content.json gains it on the next save.
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
			{Key: "show_experience", Label: "Show the experience & profiles column in the hero", Kind: Bool, Default: "1"},
			{Key: "show_about", Label: "Show About", Kind: Bool, Default: "1"},
			{Key: "show_work", Label: "Show Projects", Kind: Bool, Default: "1"},
			{Key: "show_research", Label: "Show Research", Kind: Bool, Default: "1"},
		},
	},
	{
		Key: "nav", Title: "Navbar", Description: "Link labels. Hidden sections are dropped from the navbar automatically.",
		Fields: []Field{
			{Key: "nav_about", Label: "About link", Kind: Text, Required: true, Default: "About"},
			{Key: "nav_work", Label: "Projects link", Kind: Text, Required: true, Default: "Projects"},
			{Key: "nav_research", Label: "Research link", Kind: Text, Required: true, Default: "Research"},
			{Key: "nav_cv_label", Label: "CV button", Kind: Text, Required: true, Default: "CV"},
		},
	},
	{
		Key: "hero", Title: "Hero", Description: "Research interests and the buttons under the bio (CV, the first profile group, email), plus the column on the right.",
		Fields: []Field{
			{Key: "interests_label", Label: "Interests label", Kind: Text, Default: "Research interests", Help: "Shown under the bio, before the list"},
			{Key: "interests", Label: "Interests", Kind: Lines, Help: "One per line; rendered under the bio, separated by ·", Default: "LLM inference systems\nRetrieval-augmented generation\nDistributed storage\nML systems reliability"},
			{Key: "hero_cv_label", Label: "CV button", Kind: Text, Required: true, Default: "CV"},
			{Key: "hero_email_label", Label: "Email button", Kind: Text, Required: true, Default: "Email"},
			{Key: "experience_title", Label: "Experience heading", Kind: Text, Required: true, Default: "Experience"},
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
		Key: "work", Title: "Projects section",
		Fields: []Field{
			{Key: "work_title", Label: "Section title", Kind: Text, Required: true, Default: "Projects"},
			{Key: "work_all_label", Label: "“All projects” link label", Kind: Text, Default: "All projects", Help: "Leave empty to hide the link"},
			{Key: "work_all_url", Label: "“All projects” link URL", Kind: URL, Default: "https://github.com/yourhandle?tab=repositories"},
		},
	},
	{
		Key: "research", Title: "Research section",
		Fields: []Field{
			{Key: "research_title", Label: "Section title", Kind: Text, Required: true, Default: "Research"},
			{Key: "research_statement", Label: "Research statement", Kind: Textarea, Help: "Optional paragraph shown above the papers", Default: ""},
		},
	},
	{
		Key: "footer", Title: "Footer", Description: "Email and location are shown here, under the name.",
		Fields: []Field{
			{Key: "contact_email_label", Label: "Email label (screen readers)", Kind: Text, Required: true, Default: "Email"},
			{Key: "contact_location_label", Label: "Location label (screen readers)", Kind: Text, Required: true, Default: "Location"},
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

// SettingDefaults returns key → default value, used for keys that have not
// been saved yet.
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
