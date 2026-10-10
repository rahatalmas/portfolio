package store

import (
	"context"
	"sort"
	"strconv"

	"portfolio/backend/internal/models"
	"portfolio/backend/internal/schema"
)

// snapshot returns a copy of one resource's records in display order.
func (s *Store) snapshot(key string) []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	recs := s.lists[key]
	out := make([]Record, len(recs))
	for i, r := range recs {
		out[i] = cloneRecord(r)
	}
	return out
}

// Education returns all degrees in display order.
func (s *Store) Education(ctx context.Context) ([]models.Education, error) {
	out := []models.Education{}
	for _, r := range s.snapshot("education") {
		out = append(out, models.Education{
			ID: r.ID, Degree: r.Values["degree"], School: r.Values["school"],
			Years: r.Values["years"], Note: r.Values["note"],
		})
	}
	return out, nil
}

// Skills returns all skill groups with their items split into a slice.
func (s *Store) Skills(ctx context.Context) ([]models.SkillGroup, error) {
	out := []models.SkillGroup{}
	for _, r := range s.snapshot("skills") {
		out = append(out, models.SkillGroup{ID: r.ID, Name: r.Values["name"], Items: schema.SplitLines(r.Values["items"])})
	}
	return out, nil
}

// Profiles returns profile links grouped by group_name. Groups appear in the
// order of their first link; links keep their own order within a group.
func (s *Store) Profiles(ctx context.Context) ([]models.ProfileGroup, error) {
	out := []models.ProfileGroup{}
	index := map[string]int{}
	for _, r := range s.snapshot("profiles") {
		group := r.Values["group_name"]
		i, ok := index[group]
		if !ok {
			i = len(out)
			index[group] = i
			out = append(out, models.ProfileGroup{Name: group, Items: []models.Profile{}})
		}
		out[i].Items = append(out[i].Items, models.Profile{
			ID: r.ID, Label: r.Values["label"], Href: r.Values["href"], Icon: r.Values["icon"],
		})
	}
	return out, nil
}

// Experience returns roles in display order (newest first by convention).
func (s *Store) Experience(ctx context.Context) ([]models.Role, error) {
	out := []models.Role{}
	for _, r := range s.snapshot("experience") {
		out = append(out, models.Role{ID: r.ID, Title: r.Values["title"], Company: r.Values["company"], Period: r.Values["period"]})
	}
	return out, nil
}

// Projects returns projects in display order; featuredOnly limits to
// featured ones.
func (s *Store) Projects(ctx context.Context, featuredOnly bool) ([]models.Project, error) {
	out := []models.Project{}
	for _, r := range s.snapshot("projects") {
		featured := r.Values["featured"] == "1"
		if featuredOnly && !featured {
			continue
		}
		out = append(out, models.Project{
			ID: r.ID, Slug: r.Values["slug"], Title: r.Values["title"], Summary: r.Values["summary"],
			Tags: schema.SplitComma(r.Values["tags"]), Stack: schema.SplitLines(r.Values["stack"]),
			Featured: featured,
			Links: models.ProjectLinks{
				GitHub: r.Values["github_url"], Live: r.Values["live_url"], Paper: r.Values["paper_url"],
			},
		})
	}
	return out, nil
}

// Papers returns publications newest first, then by display order.
func (s *Store) Papers(ctx context.Context) ([]models.Paper, error) {
	out := []models.Paper{}
	for _, r := range s.snapshot("papers") {
		year, _ := strconv.Atoi(r.Values["year"])
		out = append(out, models.Paper{
			ID: r.ID, Title: r.Values["title"], Authors: schema.SplitLines(r.Values["authors"]),
			Venue: r.Values["venue"], Year: year, Status: r.Values["status"],
			Links: models.PaperLinks{
				PDF: r.Values["pdf_url"], Code: r.Values["code_url"],
				Poster: r.Values["poster_url"], Slides: r.Values["slides_url"],
			},
		})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Year > out[j].Year })
	return out, nil
}

// Portfolio aggregates everything the frontend renders.
func (s *Store) Portfolio(ctx context.Context) (models.Portfolio, error) {
	var p models.Portfolio
	settings, _ := s.Settings(ctx)
	p.Site = schema.TypedSettings(settings)
	p.Education, _ = s.Education(ctx)
	p.Skills, _ = s.Skills(ctx)
	p.Profiles, _ = s.Profiles(ctx)
	p.Experience, _ = s.Experience(ctx)
	p.Projects, _ = s.Projects(ctx, false)
	p.Papers, _ = s.Papers(ctx)
	return p, nil
}
