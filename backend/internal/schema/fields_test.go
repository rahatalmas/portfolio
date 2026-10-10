package schema

import (
	"reflect"
	"testing"
)

func TestClean(t *testing.T) {
	tests := []struct {
		name    string
		f       Field
		in      string
		want    string
		wantErr bool
	}{
		{"text trims", Field{Kind: Text}, "  hi  ", "hi", false},
		{"required empty", Field{Kind: Text, Required: true, Label: "T"}, "  ", "", true},
		{"optional empty", Field{Kind: URL}, "", "", false},
		{"bool on", Field{Kind: Bool}, "on", "1", false},
		{"bool off", Field{Kind: Bool}, "", "0", false},
		{"lines drops blanks and CRLF", Field{Kind: Lines}, "a\r\n\r\n b \r\n", "a\nb", false},
		{"number ok", Field{Kind: Number}, "2025", "2025", false},
		{"number bad", Field{Kind: Number}, "20x", "", true},
		{"select ok", Field{Kind: Select, Options: []string{"a", "b"}}, "b", "b", false},
		{"select bad", Field{Kind: Select, Options: []string{"a"}}, "z", "", true},
		{"tags dedupes", Field{Kind: Tags, Options: []string{"AI", "Backend"}}, "AI, Backend,AI", "AI,Backend", false},
		{"tags unknown", Field{Kind: Tags, Options: []string{"AI"}}, "AI,Nope", "", true},
		{"url https", Field{Kind: URL}, "https://x.io", "https://x.io", false},
		{"url anchor", Field{Kind: URL}, "#research", "#research", false},
		{"url path", Field{Kind: Media}, "/uploads/a.jpg", "/uploads/a.jpg", false},
		{"url bad", Field{Kind: URL}, "javascript:alert(1)", "", true},
		{"email ok", Field{Kind: Email}, "a@b.co", "a@b.co", false},
		{"email bad", Field{Kind: Email}, "nope", "", true},
		{"slug ok", Field{Kind: Slug}, "fast-rag-2", "fast-rag-2", false},
		{"slug bad", Field{Kind: Slug}, "Fast Rag", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, problem := tt.f.Clean(tt.in)
			if (problem != "") != tt.wantErr {
				t.Fatalf("problem=%q wantErr=%v", problem, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestTyped(t *testing.T) {
	if v := (Field{Kind: Bool}).Typed("1"); v != true {
		t.Fatalf("bool: %v", v)
	}
	if v := (Field{Kind: Number}).Typed("7"); v != 7 {
		t.Fatalf("number: %v", v)
	}
	if v := (Field{Kind: Lines}).Typed("a\nb"); !reflect.DeepEqual(v, []string{"a", "b"}) {
		t.Fatalf("lines: %v", v)
	}
	if v := (Field{Kind: Tags}).Typed("AI,Backend"); !reflect.DeepEqual(v, []string{"AI", "Backend"}) {
		t.Fatalf("tags: %v", v)
	}
	if v := (Field{Kind: Lines}).Typed(""); !reflect.DeepEqual(v, []string{}) {
		t.Fatalf("empty lines should be [] not nil: %#v", v)
	}
}

func TestSettingsSpecIsConsistent(t *testing.T) {
	seen := map[string]bool{}
	for _, f := range SettingFields() {
		if f.Key == "" || f.Label == "" {
			t.Fatalf("setting with empty key/label: %+v", f)
		}
		if seen[f.Key] {
			t.Fatalf("duplicate setting key %q", f.Key)
		}
		seen[f.Key] = true
		if _, problem := f.Clean(f.Default); problem != "" {
			t.Fatalf("default for %q fails its own validation: %s", f.Key, problem)
		}
	}
}

func TestResourceSpecIsConsistent(t *testing.T) {
	keys := map[string]bool{}
	for _, r := range Resources {
		if keys[r.Key] {
			t.Fatalf("duplicate resource key %q", r.Key)
		}
		keys[r.Key] = true
		cols := map[string]bool{}
		for _, f := range r.Fields {
			if f.Key == "id" || f.Key == "sort_order" {
				t.Fatalf("%s: %q is an implicit column", r.Key, f.Key)
			}
			if cols[f.Key] {
				t.Fatalf("%s: duplicate column %q", r.Key, f.Key)
			}
			cols[f.Key] = true
			if (f.Kind == Select || f.Kind == Tags) && len(f.Options) == 0 {
				t.Fatalf("%s.%s: %s field needs Options", r.Key, f.Key, f.Kind)
			}
		}
	}
}
