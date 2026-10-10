package store

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"portfolio/backend/internal/schema"
)

func res(t *testing.T, key string) *schema.Resource {
	t.Helper()
	r, ok := schema.ResourceByKey(key)
	if !ok {
		t.Fatalf("no resource %q", key)
	}
	return r
}

func project(slug, title string) map[string]string {
	return map[string]string{
		"title": title, "slug": slug, "summary": "s", "tags": "AI", "stack": "Go",
		"featured": "1", "github_url": "", "live_url": "", "paper_url": "",
	}
}

func TestOpenCreatesContentFile(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, contentFile)); err != nil {
		t.Fatalf("content.json should be created on first run: %v", err)
	}
	settings, _ := s.Settings(context.Background())
	if settings["name"] != "Rahat Almas" {
		t.Fatalf("defaults not applied: %q", settings["name"])
	}
	for _, r := range schema.Resources {
		if n, _ := s.CountRecords(context.Background(), &r); n != 0 {
			t.Fatalf("%s should start empty", r.Key)
		}
	}
}

func TestRecordsRoundTrip(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	projects := res(t, "projects")

	id1, err := s.CreateRecord(ctx, projects, project("one", "One"))
	if err != nil {
		t.Fatal(err)
	}
	id2, err := s.CreateRecord(ctx, projects, project("two", "Two"))
	if err != nil {
		t.Fatal(err)
	}
	if id1 == id2 {
		t.Fatal("ids must differ")
	}
	if _, err := s.CreateRecord(ctx, projects, project("one", "Dup")); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate slug should fail: %v", err)
	}
	if err := s.UpdateRecord(ctx, projects, id2, project("one", "Two")); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("update to a taken slug should fail: %v", err)
	}
	if err := s.UpdateRecord(ctx, projects, id2, project("two", "Two!")); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateRecord(ctx, projects, 999, project("x", "x")); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing id: %v", err)
	}

	// Order and move.
	list, _ := s.ListRecords(ctx, projects)
	if len(list) != 2 || list[0].ID != id1 || list[1].SortOrder != 2 {
		t.Fatalf("unexpected list: %+v", list)
	}
	if err := s.MoveRecord(ctx, projects, id2, "up"); err != nil {
		t.Fatal(err)
	}
	list, _ = s.ListRecords(ctx, projects)
	if list[0].ID != id2 || list[0].SortOrder != 1 {
		t.Fatalf("move up failed: %+v", list)
	}
	if err := s.MoveRecord(ctx, projects, id2, "up"); err != nil {
		t.Fatalf("moving past the top should be a no-op: %v", err)
	}

	// Typed read sees the same data.
	typed, _ := s.Projects(ctx, true)
	if len(typed) != 2 || typed[0].Slug != "two" || typed[0].Title != "Two!" || typed[0].Tags[0] != "AI" {
		t.Fatalf("typed projects: %+v", typed)
	}

	// Reopen from disk: everything persisted.
	s2, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	list2, _ := s2.ListRecords(ctx, projects)
	if len(list2) != 2 || list2[0].ID != id2 || list2[0].Values["title"] != "Two!" {
		t.Fatalf("persisted list: %+v", list2)
	}

	if err := s2.DeleteRecord(ctx, projects, id1); err != nil {
		t.Fatal(err)
	}
	if err := s2.DeleteRecord(ctx, projects, id1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("second delete: %v", err)
	}
	if n, _ := s2.CountRecords(ctx, projects); n != 1 {
		t.Fatalf("count after delete: %d", n)
	}
	// Ids are not reused while a higher one exists.
	id3, _ := s2.CreateRecord(ctx, projects, project("three", "Three"))
	if id3 <= id2 {
		t.Fatalf("new id %d should exceed %d", id3, id2)
	}
}

func TestSettingsPersist(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	s, _ := Open(dir)
	if err := s.SaveSettings(ctx, map[string]string{"name": "Ada", "show_work": "0"}); err != nil {
		t.Fatal(err)
	}
	s2, _ := Open(dir)
	settings, _ := s2.Settings(ctx)
	if settings["name"] != "Ada" || settings["show_work"] != "0" || settings["nav_about"] != "About" {
		t.Fatalf("settings after reload: %v", settings)
	}
	p, _ := s2.Portfolio(ctx)
	if p.Site["show_work"] != false || p.Site["name"] != "Ada" {
		t.Fatalf("typed site: %v", p.Site)
	}
}

// Hand-edited files may use native JSON types, omit ids, or omit fields.
func TestLoadTolerantHandEdits(t *testing.T) {
	dir := t.TempDir()
	content := `{
  "settings": {"name": "Hand", "show_about": false, "interests": "a\nb"},
  "projects": [
    {"title": "No id", "slug": "no-id", "featured": true},
    {"id": 7, "title": "Seven", "slug": "seven", "featured": 0},
    {"id": 7, "title": "Dup id", "slug": "dup", "featured": "1"}
  ],
  "papers": [{"id": 1, "title": "P", "authors": "A", "venue": "V", "year": 2024, "status": "published"}]
}`
	if err := os.WriteFile(filepath.Join(dir, contentFile), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	settings, _ := s.Settings(ctx)
	if settings["name"] != "Hand" || settings["show_about"] != "0" || settings["nav_work"] != "Work" {
		t.Fatalf("settings: %v", settings)
	}
	list, _ := s.ListRecords(ctx, res(t, "projects"))
	if len(list) != 3 {
		t.Fatalf("want 3 projects, got %d", len(list))
	}
	seen := map[int64]bool{}
	for _, r := range list {
		if r.ID <= 0 || seen[r.ID] {
			t.Fatalf("ids must be assigned and unique: %+v", list)
		}
		seen[r.ID] = true
	}
	if list[0].Values["featured"] != "1" || list[1].Values["featured"] != "0" || list[0].Values["summary"] != "" {
		t.Fatalf("values not normalised: %+v", list)
	}
	papers, _ := s.Papers(ctx)
	if len(papers) != 1 || papers[0].Year != 2024 || papers[0].Authors[0] != "A" {
		t.Fatalf("papers: %+v", papers)
	}
	// Saving writes a well-formed file that round-trips.
	if err := s.SaveSettings(ctx, map[string]string{"name": "Hand2"}); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(dir, contentFile))
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("saved file is not valid JSON: %v\n%s", err, data)
	}
	if _, err := Open(dir); err != nil {
		t.Fatal(err)
	}
}

func TestUsers(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	s, _ := Open(dir)
	if n, _ := s.UserCount(ctx); n != 0 {
		t.Fatalf("users: %d", n)
	}
	if err := s.CreateUser(ctx, "admin", "hash1"); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateUser(ctx, "admin", "hash2"); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate user: %v", err)
	}
	if err := s.UpdatePassword(ctx, "admin", "hash3"); err != nil {
		t.Fatal(err)
	}
	s2, _ := Open(dir)
	u, err := s2.UserByUsername(ctx, "admin")
	if err != nil || u.PasswordHash != "hash3" {
		t.Fatalf("persisted user: %+v %v", u, err)
	}
	if _, err := s2.UserByUsername(ctx, "nobody"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing user: %v", err)
	}
}

// The content.json shipped in the repo must pass the same validation the
// admin applies, so a fresh install starts from a clean state.
func TestShippedContentIsValid(t *testing.T) {
	shipped := filepath.Join("..", "..", "data", contentFile)
	data, err := os.ReadFile(shipped)
	if err != nil {
		t.Skipf("no shipped content at %s: %v", shipped, err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, contentFile), data, 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	settings, _ := s.Settings(ctx)
	for _, f := range schema.SettingFields() {
		if _, problem := f.Clean(settings[f.Key]); problem != "" {
			t.Errorf("settings.%s: %s", f.Key, problem)
		}
	}
	for i := range schema.Resources {
		r := &schema.Resources[i]
		list, _ := s.ListRecords(ctx, r)
		if len(list) == 0 {
			t.Errorf("%s: shipped content should not be empty", r.Key)
		}
		for _, rec := range list {
			for _, f := range r.Fields {
				if _, problem := f.Clean(rec.Values[f.Key]); problem != "" {
					t.Errorf("%s id=%d %s: %s", r.Key, rec.ID, f.Key, problem)
				}
			}
		}
	}
}
