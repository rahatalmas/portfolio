// Package schema describes the portfolio's content model: the editable
// settings (single values such as the site name or a section title) and the
// list resources (projects, papers, …). The admin UI, the data layer and the
// API are all driven by these specs, so adding a field means editing one
// place here. Existing content.json files simply gain the field (empty, or
// the Default for settings) on the next save.
package schema

import (
	"net/mail"
	"regexp"
	"strconv"
	"strings"
)

// Kind is the input type of a field. It decides how the admin renders the
// field, how a submitted value is validated/normalised, and how the public
// API encodes it (see Typed).
type Kind string

const (
	Text     Kind = "text"     // single line
	Textarea Kind = "textarea" // multi-line prose
	Lines    Kind = "lines"    // one item per line; API returns []string
	Bool     Kind = "bool"     // checkbox; stored "1"/"0"; API returns bool
	Number   Kind = "number"   // integer; API returns number
	Select   Kind = "select"   // one of Options
	Tags     Kind = "tags"     // subset of Options; stored comma-separated; API returns []string
	URL      Kind = "url"      // absolute URL, or relative path (/…), or anchor (#…)
	Email    Kind = "email"    // RFC 5322 address
	Slug     Kind = "slug"     // lowercase letters, digits, dashes
	Media    Kind = "media"    // path to an uploaded file (/uploads/…) or any URL
)

// Field is one editable value: a field of a resource entry, or one key of
// the settings.
type Field struct {
	Key      string   // key in content.json / form input name
	Label    string   // shown in the admin
	Kind     Kind     // see Kind
	Help     string   // optional hint under the input
	Required bool     // reject empty values
	Options  []string // for Select and Tags
	Default  string   // settings only: value used until one is saved
	List     bool     // resources only: show as a column in the admin list
}

var slugRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// Clean validates and normalises a raw form value for f. It returns the value
// to store and a human-readable problem ("" when the value is fine).
func (f Field) Clean(raw string) (string, string) {
	v := strings.TrimSpace(raw)
	if f.Kind == Textarea || f.Kind == Lines {
		v = strings.ReplaceAll(raw, "\r\n", "\n")
		v = strings.TrimSpace(v)
	}

	if f.Kind == Bool {
		if v == "1" || v == "on" || v == "true" {
			return "1", ""
		}
		return "0", ""
	}
	if v == "" {
		if f.Required {
			return "", f.Label + " is required."
		}
		return "", ""
	}

	switch f.Kind {
	case Lines:
		var out []string
		for _, line := range strings.Split(v, "\n") {
			if line = strings.TrimSpace(line); line != "" {
				out = append(out, line)
			}
		}
		return strings.Join(out, "\n"), ""
	case Number:
		if _, err := strconv.Atoi(v); err != nil {
			return "", f.Label + " must be a whole number."
		}
	case Select:
		if !contains(f.Options, v) {
			return "", f.Label + " has an unknown value."
		}
	case Tags:
		var keep []string
		for _, t := range strings.Split(v, ",") {
			if t = strings.TrimSpace(t); t == "" {
				continue
			}
			if !contains(f.Options, t) {
				return "", f.Label + ": unknown option " + t + "."
			}
			if !contains(keep, t) {
				keep = append(keep, t)
			}
		}
		return strings.Join(keep, ","), ""
	case URL, Media:
		if !validURL(v) {
			return "", f.Label + " must be a URL (https://…), a path (/…) or an anchor (#…)."
		}
	case Email:
		if _, err := mail.ParseAddress(v); err != nil {
			return "", f.Label + " must be a valid email address."
		}
	case Slug:
		if !slugRe.MatchString(v) {
			return "", f.Label + " may only contain lowercase letters, digits and dashes."
		}
	}
	return v, ""
}

// Typed converts a stored string into the JSON value the public API exposes.
func (f Field) Typed(stored string) any {
	switch f.Kind {
	case Bool:
		return stored == "1"
	case Number:
		n, _ := strconv.Atoi(stored)
		return n
	case Lines:
		return SplitLines(stored)
	case Tags:
		return SplitComma(stored)
	}
	return stored
}

// SplitLines turns a Lines value into its items, dropping blanks.
func SplitLines(s string) []string {
	out := []string{}
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// SplitComma turns a Tags value into its items, dropping blanks.
func SplitComma(s string) []string {
	out := []string{}
	for _, l := range strings.Split(s, ",") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

func validURL(v string) bool {
	return strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") ||
		strings.HasPrefix(v, "mailto:") || strings.HasPrefix(v, "/") || strings.HasPrefix(v, "#")
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
