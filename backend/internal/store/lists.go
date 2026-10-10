package store

import (
	"context"
	"fmt"

	"portfolio/backend/internal/models"
	"portfolio/backend/internal/schema"
)

// Education returns all degrees in display order.
func (s *Store) Education(ctx context.Context) ([]models.Education, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, degree, school, years, note FROM education ORDER BY sort_order, id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.Education{}
	for rows.Next() {
		var e models.Education
		if err := rows.Scan(&e.ID, &e.Degree, &e.School, &e.Years, &e.Note); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// Skills returns all skill groups with their items split into a slice.
func (s *Store) Skills(ctx context.Context) ([]models.SkillGroup, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, items FROM skill_groups ORDER BY sort_order, id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.SkillGroup{}
	for rows.Next() {
		var g models.SkillGroup
		var items string
		if err := rows.Scan(&g.ID, &g.Name, &items); err != nil {
			return nil, err
		}
		g.Items = schema.SplitLines(items)
		out = append(out, g)
	}
	return out, rows.Err()
}

// Profiles returns profile links grouped by group_name. Groups appear in the
// order of their first link; links keep their own order within a group.
func (s *Store) Profiles(ctx context.Context) ([]models.ProfileGroup, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, group_name, label, href, icon FROM profiles ORDER BY sort_order, id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.ProfileGroup{}
	index := map[string]int{}
	for rows.Next() {
		var p models.Profile
		var group string
		if err := rows.Scan(&p.ID, &group, &p.Label, &p.Href, &p.Icon); err != nil {
			return nil, err
		}
		i, ok := index[group]
		if !ok {
			i = len(out)
			index[group] = i
			out = append(out, models.ProfileGroup{Name: group, Items: []models.Profile{}})
		}
		out[i].Items = append(out[i].Items, p)
	}
	return out, rows.Err()
}

// Experience returns roles in display order (newest first by convention).
func (s *Store) Experience(ctx context.Context) ([]models.Role, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, title, company, period FROM experience ORDER BY sort_order, id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.Role{}
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(&r.ID, &r.Title, &r.Company, &r.Period); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// Projects returns projects in display order; featuredOnly limits to
// featured ones.
func (s *Store) Projects(ctx context.Context, featuredOnly bool) ([]models.Project, error) {
	q := "SELECT id, slug, title, summary, tags, stack, featured, github_url, live_url, paper_url FROM projects"
	if featuredOnly {
		q += " WHERE featured = 1"
	}
	q += " ORDER BY sort_order, id"
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.Project{}
	for rows.Next() {
		var p models.Project
		var tags, stack string
		if err := rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Summary, &tags, &stack, &p.Featured,
			&p.Links.GitHub, &p.Links.Live, &p.Links.Paper); err != nil {
			return nil, err
		}
		p.Tags = schema.SplitComma(tags)
		p.Stack = schema.SplitLines(stack)
		out = append(out, p)
	}
	return out, rows.Err()
}

// Papers returns publications newest first, then by display order.
func (s *Store) Papers(ctx context.Context) ([]models.Paper, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, title, authors, venue, year, status, pdf_url, code_url, poster_url, slides_url FROM papers ORDER BY year DESC, sort_order, id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.Paper{}
	for rows.Next() {
		var p models.Paper
		var authors string
		if err := rows.Scan(&p.ID, &p.Title, &authors, &p.Venue, &p.Year, &p.Status,
			&p.Links.PDF, &p.Links.Code, &p.Links.Poster, &p.Links.Slides); err != nil {
			return nil, err
		}
		p.Authors = schema.SplitLines(authors)
		out = append(out, p)
	}
	return out, rows.Err()
}

// Portfolio aggregates everything the frontend renders.
func (s *Store) Portfolio(ctx context.Context) (models.Portfolio, error) {
	var p models.Portfolio
	var err error
	settings, err := s.Settings(ctx)
	if err != nil {
		return p, err
	}
	p.Site = schema.TypedSettings(settings)
	if p.Education, err = s.Education(ctx); err != nil {
		return p, fmt.Errorf("education: %w", err)
	}
	if p.Skills, err = s.Skills(ctx); err != nil {
		return p, fmt.Errorf("skills: %w", err)
	}
	if p.Profiles, err = s.Profiles(ctx); err != nil {
		return p, fmt.Errorf("profiles: %w", err)
	}
	if p.Experience, err = s.Experience(ctx); err != nil {
		return p, fmt.Errorf("experience: %w", err)
	}
	if p.Projects, err = s.Projects(ctx, false); err != nil {
		return p, fmt.Errorf("projects: %w", err)
	}
	if p.Papers, err = s.Papers(ctx); err != nil {
		return p, fmt.Errorf("papers: %w", err)
	}
	return p, nil
}
