package db

import "testing"

func TestSplitStatements(t *testing.T) {
	in := `-- comment
CREATE TABLE a (
  id INT -- trailing
);

INSERT INTO a VALUES (1);
SELECT 1`
	got := splitStatements(in)
	if len(got) != 3 {
		t.Fatalf("want 3 statements, got %d: %q", len(got), got)
	}
	if got[1] != "INSERT INTO a VALUES (1)" {
		t.Fatalf("unexpected second statement: %q", got[1])
	}
	if got[2] != "SELECT 1" {
		t.Fatalf("unexpected last statement: %q", got[2])
	}
}

func TestMigrationFilesAreNumbered(t *testing.T) {
	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("no migrations embedded")
	}
	for _, e := range entries {
		if len(e.Name()) < 6 || e.Name()[4] != '_' {
			t.Fatalf("migration %q must be named NNNN_description.sql", e.Name())
		}
	}
}
