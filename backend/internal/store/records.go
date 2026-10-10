package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"

	"portfolio/backend/internal/schema"
)

// Record is one row of a resource table with every column as a string, the
// representation the admin forms work in. Values are keyed by column name.
type Record struct {
	ID        int64
	SortOrder int
	Values    map[string]string
}

// ErrDuplicate is returned when a unique column (e.g. projects.slug) clashes.
var ErrDuplicate = errors.New("duplicate value")

// ListRecords returns every row of res ordered by sort_order, then id.
func (s *Store) ListRecords(ctx context.Context, res *schema.Resource) ([]Record, error) {
	cols := res.Columns()
	q := fmt.Sprintf("SELECT id, sort_order, %s FROM `%s` ORDER BY sort_order, id", quoteAll(cols), res.Table)
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list %s: %w", res.Key, err)
	}
	defer rows.Close()

	var out []Record
	for rows.Next() {
		rec, err := scanRecord(rows, cols)
		if err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

// GetRecord returns one row by id.
func (s *Store) GetRecord(ctx context.Context, res *schema.Resource, id int64) (Record, error) {
	cols := res.Columns()
	q := fmt.Sprintf("SELECT id, sort_order, %s FROM `%s` WHERE id = ?", quoteAll(cols), res.Table)
	rows, err := s.db.QueryContext(ctx, q, id)
	if err != nil {
		return Record{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return Record{}, ErrNotFound
	}
	return scanRecord(rows, cols)
}

// CreateRecord inserts a row at the end of the sort order and returns its id.
// values must contain every column of res (use Field.Clean first).
func (s *Store) CreateRecord(ctx context.Context, res *schema.Resource, values map[string]string) (int64, error) {
	cols := res.Columns()
	args := make([]any, 0, len(cols))
	for _, c := range cols {
		args = append(args, values[c])
	}
	q := fmt.Sprintf(
		"INSERT INTO `%s` (%s, sort_order) SELECT %s, COALESCE(MAX(sort_order), 0) + 1 FROM `%s`",
		res.Table, quoteAll(cols), placeholders(len(cols)), res.Table,
	)
	result, err := s.db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, wrapDuplicate(err)
	}
	return result.LastInsertId()
}

// UpdateRecord overwrites every editable column of the row with id.
func (s *Store) UpdateRecord(ctx context.Context, res *schema.Resource, id int64, values map[string]string) error {
	cols := res.Columns()
	sets := make([]string, len(cols))
	args := make([]any, 0, len(cols)+1)
	for i, c := range cols {
		sets[i] = fmt.Sprintf("`%s` = ?", c)
		args = append(args, values[c])
	}
	args = append(args, id)
	q := fmt.Sprintf("UPDATE `%s` SET %s WHERE id = ?", res.Table, strings.Join(sets, ", "))
	result, err := s.db.ExecContext(ctx, q, args...)
	if err != nil {
		return wrapDuplicate(err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		// Either missing or unchanged; distinguish with a cheap existence check.
		var one int
		if err := s.db.QueryRowContext(ctx, fmt.Sprintf("SELECT 1 FROM `%s` WHERE id = ?", res.Table), id).Scan(&one); err == sql.ErrNoRows {
			return ErrNotFound
		}
	}
	return nil
}

// DeleteRecord removes the row with id.
func (s *Store) DeleteRecord(ctx context.Context, res *schema.Resource, id int64) error {
	result, err := s.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM `%s` WHERE id = ?", res.Table), id)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// MoveRecord swaps the row's position with its neighbour. dir is "up" or
// "down"; moving past either end is a no-op.
func (s *Store) MoveRecord(ctx context.Context, res *schema.Resource, id int64, dir string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Load ids in display order so gaps/duplicates in sort_order don't matter.
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT id FROM `%s` ORDER BY sort_order, id FOR UPDATE", res.Table))
	if err != nil {
		return err
	}
	var ids []int64
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			rows.Close()
			return err
		}
		ids = append(ids, v)
	}
	rows.Close()

	pos := -1
	for i, v := range ids {
		if v == id {
			pos = i
		}
	}
	if pos == -1 {
		return ErrNotFound
	}
	swap := pos - 1
	if dir == "down" {
		swap = pos + 1
	}
	if swap < 0 || swap >= len(ids) {
		return nil
	}
	ids[pos], ids[swap] = ids[swap], ids[pos]

	// Renumber everything 1..n; cheap for the few dozen rows a portfolio has
	// and it heals any earlier inconsistencies.
	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf("UPDATE `%s` SET sort_order = ? WHERE id = ?", res.Table))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i, v := range ids {
		if _, err := stmt.ExecContext(ctx, i+1, v); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// CountRecords returns the number of rows in res.
func (s *Store) CountRecords(ctx context.Context, res *schema.Resource) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM `%s`", res.Table)).Scan(&n)
	return n, err
}

func scanRecord(rows *sql.Rows, cols []string) (Record, error) {
	rec := Record{Values: make(map[string]string, len(cols))}
	dest := make([]any, 0, len(cols)+2)
	dest = append(dest, &rec.ID, &rec.SortOrder)
	vals := make([]sql.NullString, len(cols))
	for i := range vals {
		dest = append(dest, &vals[i])
	}
	if err := rows.Scan(dest...); err != nil {
		return rec, err
	}
	for i, c := range cols {
		rec.Values[c] = vals[i].String
	}
	return rec, nil
}

func quoteAll(cols []string) string {
	q := make([]string, len(cols))
	for i, c := range cols {
		q[i] = "`" + c + "`"
	}
	return strings.Join(q, ", ")
}

func placeholders(n int) string {
	return strings.TrimSuffix(strings.Repeat("?, ", n), ", ")
}

func wrapDuplicate(err error) error {
	var me *mysql.MySQLError
	if errors.As(err, &me) && me.Number == 1062 {
		return ErrDuplicate
	}
	return err
}
