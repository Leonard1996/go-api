package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pack-calculator/internal/packutil"
)

type PackSizeRepository interface {
	ListPackSizes(ctx context.Context) ([]int, error)
	ReplacePackSizes(ctx context.Context, sizes []int) error
}

type SQLiteRepo struct {
	db *sql.DB
}

var (
	ErrNilDB            = errors.New("db is nil")
	ErrInvalidPackSizes = errors.New("pack sizes must be positive")
)

func NewSQLiteRepo(db *sql.DB) (*SQLiteRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}

	repo := &SQLiteRepo{db: db}
	if err := repo.init(); err != nil {
		return nil, fmt.Errorf("init sqlite repo %w", err)
	}

	return repo, nil
}

func (r *SQLiteRepo) ListPackSizes(ctx context.Context) ([]int, error) {
	rows, err := r.db.QueryContext(ctx, `select size from pack_sizes order by size asc`)
	if err != nil {
		return nil, fmt.Errorf("list pack sizes %w", err)
	}
	defer rows.Close()

	var sizes []int
	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, fmt.Errorf("list pack sizes %w", err)
		}
		sizes = append(sizes, size)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list pack sizes %w", err)
	}

	return sizes, nil
}

func (r *SQLiteRepo) ReplacePackSizes(ctx context.Context, sizes []int) error {
	clean, err := packutil.NormalizePackSizes(sizes, ErrInvalidPackSizes)
	if err != nil {
		return fmt.Errorf("replace pack sizes %w", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("replace pack sizes %w", err)
	}

	if _, err := tx.ExecContext(ctx, `delete from pack_sizes`); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("replace pack sizes %w", err)
	}

	for _, size := range clean {
		if _, err := tx.ExecContext(ctx, `insert into pack_sizes(size) values(?)`, size); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("replace pack sizes %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("replace pack sizes %w", err)
	}

	return nil
}

func (r *SQLiteRepo) init() error {
	_, err := r.db.Exec(`create table if not exists pack_sizes (size integer primary key)`)
	return err
}
