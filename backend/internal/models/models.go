// Package models holds the typed shapes the public API returns. They mirror
// the TypeScript types in the frontend (src/api/types.ts).
package models

// Education is one degree.
type Education struct {
	ID     int64  `json:"id"`
	Degree string `json:"degree"`
	School string `json:"school"`
	Years  string `json:"years"`
	Note   string `json:"note"`
}

// SkillGroup is a named list of skills ("AI": PyTorch, …).
type SkillGroup struct {
	ID    int64    `json:"id"`
	Name  string   `json:"name"`
	Items []string `json:"items"`
}

// Profile is one external profile link.
type Profile struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
	Href  string `json:"href"`
	Icon  string `json:"icon"`
}

// ProfileGroup is a row of profile links sharing a group name.
type ProfileGroup struct {
	Name  string    `json:"name"`
	Items []Profile `json:"items"`
}

// Role is one experience entry on the hero timeline card.
type Role struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Company string `json:"company"`
	Period  string `json:"period"`
}

// ProjectLinks are a project's optional outbound links.
type ProjectLinks struct {
	GitHub string `json:"github,omitempty"`
	Live   string `json:"live,omitempty"`
	Paper  string `json:"paper,omitempty"`
}

// Project is one portfolio project.
type Project struct {
	ID       int64        `json:"id"`
	Slug     string       `json:"slug"`
	Title    string       `json:"title"`
	Summary  string       `json:"summary"`
	Tags     []string     `json:"tags"`
	Stack    []string     `json:"stack"`
	Featured bool         `json:"featured"`
	Links    ProjectLinks `json:"links"`
}

// PaperLinks are a paper's optional outbound links.
type PaperLinks struct {
	PDF    string `json:"pdf,omitempty"`
	Code   string `json:"code,omitempty"`
	Poster string `json:"poster,omitempty"`
	Slides string `json:"slides,omitempty"`
}

// Paper is one publication.
type Paper struct {
	ID      int64      `json:"id"`
	Title   string     `json:"title"`
	Authors []string   `json:"authors"`
	Venue   string     `json:"venue"`
	Year    int        `json:"year"`
	Status  string     `json:"status"`
	Links   PaperLinks `json:"links"`
}

// Portfolio is everything the frontend needs, in one response.
type Portfolio struct {
	Site       map[string]any `json:"site"`
	Education  []Education    `json:"education"`
	Skills     []SkillGroup   `json:"skills"`
	Profiles   []ProfileGroup `json:"profiles"`
	Experience []Role         `json:"experience"`
	Projects   []Project      `json:"projects"`
	Papers     []Paper        `json:"papers"`
}

// Message is a contact-form submission.
type Message struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Body      string `json:"body"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
}

// User is an admin account. It is never returned by the API.
type User struct {
	ID           int64
	Username     string
	PasswordHash string
}
