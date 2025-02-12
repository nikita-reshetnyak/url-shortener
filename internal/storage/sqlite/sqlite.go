package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
		CREATE INDEX IF NOT EXIST idx_alias ON url(alias)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}
func (s *Storage) SaveUrl(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(url,alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: preapare statement: %w", op, err)
	}
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUrlExist)
		}
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)

	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	return id, nil
}
func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	var als string
	err = stmt.QueryRow(alias).Scan(&als)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrUrlNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: no row matches the query: %w ", op, err)
	}
	return als, nil
}
func (s *Storage) DeleteUrl(alias string) (int64, error) {
	const op = "storage.sqlite.DeleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	row, err := stmt.Exec(alias)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to delete url: %w", op, err)
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}
	if rowsAffected == 0 {
		return 0, storage.ErrUrlNotFound
	}
	return rowsAffected, nil

}
