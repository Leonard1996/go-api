package repo

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

func openTestDB(t *testing.T) *sql.DB {
	file, err := os.CreateTemp("", "packs-*.db")
	require.NoError(t, err)
	path := file.Name()
	require.NoError(t, file.Close())

	t.Cleanup(func() {
		_ = os.Remove(path)
	})

	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	db.SetMaxOpenConns(1)

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func TestSQLiteRepoReplaceAndList(t *testing.T) {
	db := openTestDB(t)
	repo, err := NewSQLiteRepo(db)
	require.NoError(t, err)

	ctx := context.Background()
	err = repo.ReplacePackSizes(ctx, []int{500, 250, 1000, 250})
	require.NoError(t, err)

	sizes, err := repo.ListPackSizes(ctx)
	require.NoError(t, err)
	assert.Equal(t, []int{250, 500, 1000}, sizes)
}

func TestSQLiteRepoReplaceValidation(t *testing.T) {
	db := openTestDB(t)
	repo, err := NewSQLiteRepo(db)
	require.NoError(t, err)

	ctx := context.Background()
	assert.ErrorIs(t, repo.ReplacePackSizes(ctx, []int{}), ErrInvalidPackSizes)
	assert.ErrorIs(t, repo.ReplacePackSizes(ctx, []int{0, 10}), ErrInvalidPackSizes)
}

func TestSQLiteRepoListEmpty(t *testing.T) {
	db := openTestDB(t)
	repo, err := NewSQLiteRepo(db)
	require.NoError(t, err)

	ctx := context.Background()
	sizes, err := repo.ListPackSizes(ctx)
	require.NoError(t, err)
	assert.Empty(t, sizes)
}
