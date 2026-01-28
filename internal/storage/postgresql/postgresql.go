package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

func New(dsn string) (*PostgresStorage, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Save(id, url string) error {
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO url(short_url, full_url) VALUES($1, $2)", id, url)
	return  err
}

func (s *PostgresStorage) Get(id string) (string, bool, error) {
	var fullURL string
	err := s.db.QueryRow(context.Background(),
		"SELECT full_url FROM url WHERE short_url = $1", id).Scan(&fullURL)

	if err != nil {
		return "", false, err
	}
	return fullURL, true, nil
}

