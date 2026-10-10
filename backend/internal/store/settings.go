package store

import (
	"context"
	"fmt"

	"portfolio/backend/internal/schema"
)

// Settings returns every stored key → value, with schema defaults filled in
// for keys that have not been saved yet.
func (s *Store) Settings(ctx context.Context) (map[string]string, error) {
	out := schema.SettingDefaults()
	rows, err := s.db.QueryContext(ctx, "SELECT `key`, `value` FROM settings")
	if err != nil {
		return nil, fmt.Errorf("load settings: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, rows.Err()
}

// SaveSettings upserts the given keys in one transaction.
func (s *Store) SaveSettings(ctx context.Context, values map[string]string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO settings (`key`, `value`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `value` = VALUES(`value`)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for k, v := range values {
		if _, err := stmt.ExecContext(ctx, k, v); err != nil {
			return fmt.Errorf("save setting %q: %w", k, err)
		}
	}
	return tx.Commit()
}

// SeedSettings inserts defaults for keys that are missing, leaving existing
// values untouched. Run on every start so new settings appear automatically.
func (s *Store) SeedSettings(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT IGNORE INTO settings (`key`, `value`) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for k, v := range schema.SettingDefaults() {
		if _, err := stmt.ExecContext(ctx, k, v); err != nil {
			return fmt.Errorf("seed setting %q: %w", k, err)
		}
	}
	return tx.Commit()
}
