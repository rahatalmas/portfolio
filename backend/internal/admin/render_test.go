package admin

import (
	"log/slog"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"portfolio/backend/internal/media"
	"portfolio/backend/internal/schema"
)

// Every page must parse and execute with representative data, so template
// mistakes fail here instead of at a user's first click.
func TestAllPagesRender(t *testing.T) {
	r, err := newRenderer(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	if err != nil {
		t.Fatal(err)
	}
	proj, _ := schema.ResourceByKey("projects")
	base := view{Title: "T", User: "admin", Handle: "rahat", CSRF: "tok", Resources: schema.Resources, Flash: "hi", FlashKind: "ok"}

	fields := func(values map[string]string) []fieldView {
		var out []fieldView
		for _, f := range proj.Fields {
			fv := fieldView{F: f, Value: values[f.Key], Checked: map[string]bool{}}
			if f.Kind == schema.Tags {
				for _, tag := range schema.SplitComma(fv.Value) {
					fv.Checked[tag] = true
				}
			}
			out = append(out, fv)
		}
		return out
	}

	var groups []groupView
	for _, g := range schema.SettingsGroups {
		gv := groupView{G: g}
		for _, f := range g.Fields {
			gv.Fields = append(gv.Fields, fieldView{F: f, Value: f.Default, Checked: map[string]bool{}})
		}
		groups = append(groups, gv)
	}

	cases := map[string]any{
		"login":     map[string]string{"Username": "x"},
		"dashboard": dashboardData{Counts: []resourceCount{{Res: *proj, Count: 3}}},
		"settings":  settingsData{Groups: groups, Media: []string{"/uploads/a.jpg"}},
		"list": listData{Res: *proj, Columns: proj.Fields[:2], Rows: []listRow{
			{ID: 1, Cells: []string{"One", "one"}, First: true}, {ID: 2, Cells: []string{"Two", "two"}, Last: true},
		}},
		"form":     formData{Res: *proj, ID: 7, Fields: fields(map[string]string{"title": "T", "tags": "AI,Backend", "featured": "1"}), Media: []string{"/uploads/a.jpg"}},
		"media":    []media.File{{Name: "a.jpg", URL: "/uploads/a.jpg", Size: 2048, IsImage: true}, {Name: "cv.pdf", URL: "/uploads/cv.pdf", Size: 5 << 20}},
		"password": nil,
	}
	for page, data := range cases {
		t.Run(page, func(t *testing.T) {
			v := base
			v.Data = data
			if page == "login" {
				v.User = ""
			}
			rec := httptest.NewRecorder()
			r.render(rec, 200, page, v)
			body := rec.Body.String()
			if !strings.Contains(body, "</html>") {
				t.Fatalf("page %s did not render fully:\n%s", page, body)
			}
			if page != "login" && !strings.Contains(body, `name="_csrf" value="tok"`) {
				t.Fatalf("page %s is missing the CSRF token", page)
			}
		})
	}
	for page := range r.pages {
		if _, ok := cases[page]; !ok {
			t.Errorf("page %q has no render test case", page)
		}
	}
}

func TestFormTagsRenderChecked(t *testing.T) {
	r, err := newRenderer(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	if err != nil {
		t.Fatal(err)
	}
	proj, _ := schema.ResourceByKey("projects")
	tags, _ := proj.Field("tags")
	fv := fieldView{F: tags, Value: "AI", Checked: map[string]bool{"AI": true}}
	rec := httptest.NewRecorder()
	r.render(rec, 200, "form", view{User: "u", CSRF: "t", Resources: schema.Resources, Data: formData{Res: *proj, Fields: []fieldView{fv}}})
	body := rec.Body.String()
	if !strings.Contains(body, `value="AI" checked`) || strings.Contains(body, `value="Backend" checked`) {
		t.Fatalf("tags checkbox state wrong:\n%s", body)
	}
}
