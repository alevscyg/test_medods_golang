package sqlstorage

import (
	"database/sql"

	"medos/storage"

	_ "github.com/lib/pq" // ...
)

// Storage ...
type Storage struct {
	db                *sql.DB
	RefreshRepository *RefreshTokenRepository
}

// New ...
func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// Auth ...
func (s *Storage) Auth() storage.RefreshTokenRepository {
	if s.RefreshRepository != nil {
		return s.RefreshRepository
	}

	return &RefreshTokenRepository{
		storage: s,
	}
}
