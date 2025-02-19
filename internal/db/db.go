package db

import (
	"database/sql"
	"strings"
)

func CreateTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_url TEXT NOT NULL UNIQUE,
		shortened_url TEXT NOT NULL
	)`
	_, err := db.Exec(query)
	return err
}

func StoreURL(db *sql.DB, originalURL, shortenedURL string) error {
	originalURL = strings.TrimSpace(strings.ToLower(originalURL))
	shortenedURL = strings.TrimSpace(strings.ToLower(shortenedURL))
	query := `
		INSERT INTO urls (original_url, shortened_url)
		VALUES ($1, $2)
		ON CONFLICT (original_url)
    	DO UPDATE SET shortened_url = EXCLUDED.shortened_url
	`
	_, err := db.Exec(query, originalURL, shortenedURL)
	return err
}

func GetOriginalURL(db *sql.DB, shortenedURL string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE shortened_url = $1`
	err := db.QueryRow(query, shortenedURL).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}
