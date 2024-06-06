package storage

import (
	"database/sql"
	"fmt"
)

func NewPostStorage(dsn string) (*PostStorage, error) {
	const op = "domain.storage.NewPostStorage"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &PostStorage{db: db}, nil
}
func (s *PostStorage) Stop(db *sql.DB) error {
	return s.db.Close()
}
