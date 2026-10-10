package store

import (
	"testing"

	"portfolio/backend/internal/schema"
)

// Every demo row must pass the same validation the admin applies, and must
// only use columns the resource actually has.
func TestDemoContentMatchesSchema(t *testing.T) {
	for key, rows := range demoContent {
		res, ok := schema.ResourceByKey(key)
		if !ok {
			t.Fatalf("demo content for unknown resource %q", key)
		}
		for i, row := range rows {
			for col := range row {
				if _, ok := res.Field(col); !ok {
					t.Fatalf("%s[%d]: unknown column %q", key, i, col)
				}
			}
			for _, f := range res.Fields {
				if _, problem := f.Clean(row[f.Key]); problem != "" {
					t.Fatalf("%s[%d].%s: %s", key, i, f.Key, problem)
				}
			}
		}
	}
}
