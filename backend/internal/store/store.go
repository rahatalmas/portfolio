// Package store is the data layer. Everything lives in two JSON files in
// DATA_DIR, hand-editable and small enough to keep whole in memory:
//
//   - content.json — settings plus every list resource (education, projects, …)
//   - admin.json   — admin accounts (bcrypt hashes)
//
// Reads come from memory; every write re-serialises the affected file and
// replaces it atomically (write temp, rename), so a crash mid-save cannot
// leave a half-written file. A single mutex serialises writers, which is
// plenty for one admin editing a portfolio.
//
// The method set mirrors what the admin and API handlers need:
//
//   - Generic, spec-driven CRUD over any schema.Resource (records.go section)
//   - Typed reads for the public API (lists.go)
//   - Admin users (users.go)
package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	"portfolio/backend/internal/schema"
)

// ErrNotFound is returned when an entry does not exist.
var ErrNotFound = errors.New("not found")

// ErrDuplicate is returned when a unique value (a project slug, a username)
// clashes with an existing entry.
var ErrDuplicate = errors.New("duplicate value")

const (
	contentFile = "content.json"
	adminFile   = "admin.json"
)

// Record is one entry of a list resource with every field as a string, the
// representation the admin forms work in. Values are keyed by field key.
// SortOrder is the 1-based position in display order.
type Record struct {
	ID        int64
	SortOrder int
	Values    map[string]string
}

// Store holds the data files in memory and persists changes to dir.
type Store struct {
	dir string

	mu       sync.RWMutex
	settings map[string]string
	lists    map[string][]Record // resource key → records in display order
	users    []user
}

// Open loads the data files from dir, creating the directory and an initial
// content.json (schema defaults, empty lists) when they do not exist.
func Open(dir string) (*Store, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	s := &Store{dir: abs, lists: map[string][]Record{}}
	if err := s.loadContent(); err != nil {
		return nil, fmt.Errorf("%s: %w", contentFile, err)
	}
	if err := s.loadJSON(adminFile, &s.users); err != nil {
		return nil, fmt.Errorf("%s: %w", adminFile, err)
	}
	return s, nil
}

// Dir returns the absolute data directory.
func (s *Store) Dir() string { return s.dir }

// --- files -----------------------------------------------------------------

func (s *Store) path(name string) string { return filepath.Join(s.dir, name) }

// loadJSON reads name into v; a missing file leaves v untouched.
func (s *Store) loadJSON(name string, v any) error {
	data, err := os.ReadFile(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

// saveJSON writes v to name atomically, pretty-printed for hand editing.
func (s *Store) saveJSON(name string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomic(s.path(name), append(data, '\n'))
}

func writeAtomic(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// --- content.json ----------------------------------------------------------

// loadContent reads content.json: {"settings": {...}, "<resource>": [...]}.
// Records are flat objects ({"id": 1, "title": "…"}); scalars of any JSON
// type are accepted and stored as strings so hand edits are forgiving. Rows
// without an id get one. Missing settings fall back to schema defaults.
func (s *Store) loadContent() error {
	var raw map[string]json.RawMessage
	if err := s.loadJSON(contentFile, &raw); err != nil {
		return err
	}

	s.settings = schema.SettingDefaults()
	if rawSettings, ok := raw["settings"]; ok {
		var m map[string]any
		if err := json.Unmarshal(rawSettings, &m); err != nil {
			return fmt.Errorf("settings: %w", err)
		}
		for k, v := range m {
			s.settings[k] = stringify(v)
		}
	}

	for i := range schema.Resources {
		res := &schema.Resources[i]
		var rows []map[string]any
		if rawRows, ok := raw[res.Key]; ok {
			if err := json.Unmarshal(rawRows, &rows); err != nil {
				return fmt.Errorf("%s: %w", res.Key, err)
			}
		}
		recs := make([]Record, 0, len(rows))
		var maxID int64
		for _, row := range rows {
			rec := Record{Values: make(map[string]string, len(res.Fields))}
			if idv, ok := row["id"]; ok {
				rec.ID, _ = strconv.ParseInt(stringify(idv), 10, 64)
			}
			for _, f := range res.Fields {
				if v, ok := row[f.Key]; ok {
					rec.Values[f.Key] = stringify(v)
				} else {
					rec.Values[f.Key] = ""
				}
			}
			if rec.ID > maxID {
				maxID = rec.ID
			}
			recs = append(recs, rec)
		}
		// Assign ids to rows that have none or repeat another row's id.
		seen := map[int64]bool{}
		for j := range recs {
			if recs[j].ID <= 0 || seen[recs[j].ID] {
				maxID++
				recs[j].ID = maxID
			}
			seen[recs[j].ID] = true
			recs[j].SortOrder = j + 1
		}
		s.lists[res.Key] = recs
	}

	if raw == nil {
		// First run: materialise the file so it can be edited by hand.
		return s.saveContent()
	}
	return nil
}

// saveContent writes content.json with settings first (in schema order) and
// then each resource's records in display order with fields in form order.
func (s *Store) saveContent() error {
	var b bytes.Buffer
	b.WriteString("{\n  \"settings\": {\n")
	fields := schema.SettingFields()
	known := map[string]bool{}
	for i, f := range fields {
		known[f.Key] = true
		writeKV(&b, "    ", f.Key, s.settings[f.Key], i == len(fields)-1 && !hasExtra(s.settings, known))
	}
	// Keep unknown keys (e.g. a setting removed from the schema) rather than
	// silently dropping what someone typed into the file.
	extra := extraKeys(s.settings, known)
	for i, k := range extra {
		writeKV(&b, "    ", k, s.settings[k], i == len(extra)-1)
	}
	b.WriteString("  }")

	for _, res := range schema.Resources {
		b.WriteString(",\n  ")
		b.Write(mustJSON(res.Key))
		b.WriteString(": [")
		recs := s.lists[res.Key]
		for i, rec := range recs {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString("\n    {\n")
			writeKV(&b, "      ", "id", json.Number(strconv.FormatInt(rec.ID, 10)), len(res.Fields) == 0)
			for j, f := range res.Fields {
				writeKV(&b, "      ", f.Key, rec.Values[f.Key], j == len(res.Fields)-1)
			}
			b.WriteString("    }")
		}
		if len(recs) > 0 {
			b.WriteString("\n  ")
		}
		b.WriteString("]")
	}
	b.WriteString("\n}\n")

	// Validate what we produced before touching the file.
	if !json.Valid(b.Bytes()) {
		return errors.New("internal error: produced invalid content.json")
	}
	return writeAtomic(s.path(contentFile), b.Bytes())
}

func writeKV(b *bytes.Buffer, indent, key string, value any, last bool) {
	b.WriteString(indent)
	b.Write(mustJSON(key))
	b.WriteString(": ")
	b.Write(mustJSON(value))
	if !last {
		b.WriteString(",")
	}
	b.WriteString("\n")
}

func mustJSON(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func hasExtra(m map[string]string, known map[string]bool) bool {
	return len(extraKeys(m, known)) > 0
}

func extraKeys(m map[string]string, known map[string]bool) []string {
	var out []string
	for k := range m {
		if !known[k] {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

// stringify turns any JSON scalar into the string form the schema stores:
// bools become "1"/"0", numbers keep their text, null becomes "".
func stringify(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case bool:
		if x {
			return "1"
		}
		return "0"
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case json.Number:
		return x.String()
	default:
		return fmt.Sprint(x)
	}
}

// --- settings --------------------------------------------------------------

// Settings returns every key → value, with schema defaults for keys that
// have never been saved.
func (s *Store) Settings(ctx context.Context) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := schema.SettingDefaults()
	for k, v := range s.settings {
		out[k] = v
	}
	return out, nil
}

// SaveSettings stores the given keys.
func (s *Store) SaveSettings(ctx context.Context, values map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range values {
		s.settings[k] = v
	}
	return s.saveContent()
}

// --- records ---------------------------------------------------------------

func cloneRecord(r Record) Record {
	vals := make(map[string]string, len(r.Values))
	for k, v := range r.Values {
		vals[k] = v
	}
	return Record{ID: r.ID, SortOrder: r.SortOrder, Values: vals}
}

// ListRecords returns every entry of res in display order.
func (s *Store) ListRecords(ctx context.Context, res *schema.Resource) ([]Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	recs := s.lists[res.Key]
	out := make([]Record, len(recs))
	for i, r := range recs {
		out[i] = cloneRecord(r)
	}
	return out, nil
}

// GetRecord returns one entry by id.
func (s *Store) GetRecord(ctx context.Context, res *schema.Resource, id int64) (Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i := s.index(res, id); i >= 0 {
		return cloneRecord(s.lists[res.Key][i]), nil
	}
	return Record{}, ErrNotFound
}

// CountRecords returns how many entries res has.
func (s *Store) CountRecords(ctx context.Context, res *schema.Resource) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.lists[res.Key]), nil
}

// CreateRecord appends an entry and returns its id. values should contain
// every field of res (use Field.Clean first).
func (s *Store) CreateRecord(ctx context.Context, res *schema.Resource, values map[string]string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkUnique(res, values, 0); err != nil {
		return 0, err
	}
	var maxID int64
	for _, r := range s.lists[res.Key] {
		if r.ID > maxID {
			maxID = r.ID
		}
	}
	rec := Record{ID: maxID + 1, Values: make(map[string]string, len(res.Fields))}
	for _, f := range res.Fields {
		rec.Values[f.Key] = values[f.Key]
	}
	s.lists[res.Key] = append(s.lists[res.Key], rec)
	s.renumber(res)
	if err := s.saveContent(); err != nil {
		return 0, err
	}
	return rec.ID, nil
}

// UpdateRecord overwrites every field of the entry with id.
func (s *Store) UpdateRecord(ctx context.Context, res *schema.Resource, id int64, values map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.index(res, id)
	if i < 0 {
		return ErrNotFound
	}
	if err := s.checkUnique(res, values, id); err != nil {
		return err
	}
	for _, f := range res.Fields {
		s.lists[res.Key][i].Values[f.Key] = values[f.Key]
	}
	return s.saveContent()
}

// DeleteRecord removes the entry with id.
func (s *Store) DeleteRecord(ctx context.Context, res *schema.Resource, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.index(res, id)
	if i < 0 {
		return ErrNotFound
	}
	recs := s.lists[res.Key]
	s.lists[res.Key] = append(recs[:i:i], recs[i+1:]...)
	s.renumber(res)
	return s.saveContent()
}

// MoveRecord swaps the entry with its neighbour. dir is "up" or "down";
// moving past either end is a no-op.
func (s *Store) MoveRecord(ctx context.Context, res *schema.Resource, id int64, dir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.index(res, id)
	if i < 0 {
		return ErrNotFound
	}
	j := i - 1
	if dir == "down" {
		j = i + 1
	}
	recs := s.lists[res.Key]
	if j < 0 || j >= len(recs) {
		return nil
	}
	recs[i], recs[j] = recs[j], recs[i]
	s.renumber(res)
	return s.saveContent()
}

// index returns the position of id within res, or -1. Caller holds mu.
func (s *Store) index(res *schema.Resource, id int64) int {
	for i, r := range s.lists[res.Key] {
		if r.ID == id {
			return i
		}
	}
	return -1
}

// renumber refreshes SortOrder after any reorder. Caller holds mu.
func (s *Store) renumber(res *schema.Resource) {
	for i := range s.lists[res.Key] {
		s.lists[res.Key][i].SortOrder = i + 1
	}
}

// checkUnique enforces that Slug fields are unique within res, ignoring the
// entry with id self. Caller holds mu.
func (s *Store) checkUnique(res *schema.Resource, values map[string]string, self int64) error {
	for _, f := range res.Fields {
		if f.Kind != schema.Slug {
			continue
		}
		v := values[f.Key]
		if v == "" {
			continue
		}
		for _, r := range s.lists[res.Key] {
			if r.ID != self && r.Values[f.Key] == v {
				return ErrDuplicate
			}
		}
	}
	return nil
}
