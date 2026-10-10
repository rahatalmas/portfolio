package schema

// Resource describes one list-type content table (projects, papers, …).
// Every resource table has the implicit columns `id` and `sort_order`; the
// Fields are its remaining columns, in admin form order.
type Resource struct {
	Key      string  // URL segment in the admin and the public API, e.g. "projects"
	Table    string  // MySQL table name
	Singular string  // "Project"
	Plural   string  // "Projects"
	Fields   []Field // editable columns
}

// ProfileIcons are the icon keys the frontend knows how to draw for a profile
// link. Keep in sync with ICONS in src/components/Profiles.tsx.
var ProfileIcons = []string{
	"github", "linkedin", "x", "scholar", "orcid", "research",
	"code", "trophy", "chef", "terminal", "chart", "mail", "globe",
}

// ProjectTags are the categories a project can carry.
var ProjectTags = []string{"AI", "Backend", "Research"}

// PaperStatuses are the publication states a paper can be in.
var PaperStatuses = []string{"published", "preprint", "in progress"}

// Resources lists every list-type content table, in the order the admin
// sidebar shows them. Field order is form order; List marks list columns.
var Resources = []Resource{
	{
		Key: "education", Table: "education", Singular: "Education entry", Plural: "Education",
		Fields: []Field{
			{Key: "degree", Label: "Degree", Kind: Text, Required: true, List: true, Help: "e.g. B.Sc. Computer Science"},
			{Key: "school", Label: "School", Kind: Text, Required: true, List: true},
			{Key: "years", Label: "Years", Kind: Text, List: true, Help: "e.g. 2020 – 2024"},
			{Key: "note", Label: "Note", Kind: Text, Help: "GPA, thesis title, honours — shown under the degree"},
		},
	},
	{
		Key: "skills", Table: "skill_groups", Singular: "Skill group", Plural: "Skills",
		Fields: []Field{
			{Key: "name", Label: "Group name", Kind: Text, Required: true, List: true, Help: "e.g. AI, Backend, Tools"},
			{Key: "items", Label: "Skills", Kind: Lines, Required: true, List: true, Help: "One skill per line, in display order"},
		},
	},
	{
		Key: "profiles", Table: "profiles", Singular: "Profile link", Plural: "Profiles",
		Fields: []Field{
			{Key: "group_name", Label: "Group", Kind: Text, Required: true, List: true, Help: "Links with the same group name are shown on one row, e.g. Professional, Academic, Coding"},
			{Key: "label", Label: "Label", Kind: Text, Required: true, List: true},
			{Key: "href", Label: "URL", Kind: URL, Required: true, List: true},
			{Key: "icon", Label: "Icon", Kind: Select, Required: true, Options: ProfileIcons},
		},
	},
	{
		Key: "experience", Table: "experience", Singular: "Role", Plural: "Experience",
		Fields: []Field{
			{Key: "title", Label: "Title", Kind: Text, Required: true, List: true, Help: "Keep it short; the hero card is narrow"},
			{Key: "company", Label: "Company", Kind: Text, Required: true, List: true},
			{Key: "period", Label: "Period", Kind: Text, Required: true, List: true, Help: "e.g. Jun 2024 – Present; the card shows the years only"},
		},
	},
	{
		Key: "projects", Table: "projects", Singular: "Project", Plural: "Projects",
		Fields: []Field{
			{Key: "title", Label: "Title", Kind: Text, Required: true, List: true},
			{Key: "slug", Label: "Slug", Kind: Slug, Required: true, Help: "Stable identifier, lowercase-with-dashes"},
			{Key: "summary", Label: "Summary", Kind: Textarea, Required: true, Help: "One sentence with the result, e.g. “cut p95 latency 40%”"},
			{Key: "tags", Label: "Tags", Kind: Tags, Options: ProjectTags, List: true},
			{Key: "stack", Label: "Tech stack", Kind: Lines, Help: "One technology per line"},
			{Key: "featured", Label: "Featured", Kind: Bool, List: true, Help: "Featured projects appear on the site (up to 4)"},
			{Key: "github_url", Label: "Code URL", Kind: URL},
			{Key: "live_url", Label: "Live URL", Kind: URL},
			{Key: "paper_url", Label: "Paper URL", Kind: URL, Help: "Use #research to jump to the Research section"},
		},
	},
	{
		Key: "papers", Table: "papers", Singular: "Paper", Plural: "Papers",
		Fields: []Field{
			{Key: "title", Label: "Title", Kind: Text, Required: true, List: true},
			{Key: "authors", Label: "Authors", Kind: Lines, Required: true, Help: "One author per line; your own name is highlighted automatically"},
			{Key: "venue", Label: "Venue", Kind: Text, Required: true, List: true},
			{Key: "year", Label: "Year", Kind: Number, Required: true, List: true},
			{Key: "status", Label: "Status", Kind: Select, Options: PaperStatuses, Required: true, List: true},
			{Key: "pdf_url", Label: "PDF URL", Kind: URL},
			{Key: "code_url", Label: "Code URL", Kind: URL},
			{Key: "poster_url", Label: "Poster URL", Kind: URL},
			{Key: "slides_url", Label: "Slides URL", Kind: URL},
		},
	},
}

// ResourceByKey returns the resource with the given URL key.
func ResourceByKey(key string) (*Resource, bool) {
	for i := range Resources {
		if Resources[i].Key == key {
			return &Resources[i], true
		}
	}
	return nil, false
}

// Columns returns the editable column names in field order.
func (r Resource) Columns() []string {
	cols := make([]string, len(r.Fields))
	for i, f := range r.Fields {
		cols[i] = f.Key
	}
	return cols
}

// Field returns the field with the given key.
func (r Resource) Field(key string) (Field, bool) {
	for _, f := range r.Fields {
		if f.Key == key {
			return f, true
		}
	}
	return Field{}, false
}
