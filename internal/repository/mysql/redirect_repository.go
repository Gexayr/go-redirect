package mysql

import (
	"database/sql"
	"fmt"
	"platform/internal/models"
)

type RedirectRepository struct {
	db *sql.DB
}

func NewRedirectRepository(db *sql.DB) *RedirectRepository {
	return &RedirectRepository{
		db: db,
	}
}

func (r *RedirectRepository) GetRedirectURL(hash string) (string, error) {
	query := `
		SELECT redirect_url 
		FROM redirect_mappings 
		WHERE hash = ?
	`

	var redirectURL string
	err := r.db.QueryRow(query, hash).Scan(&redirectURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get redirect URL: %w", err)
	}

	return redirectURL, nil
}

func (r *RedirectRepository) SaveRedirect(redirect *models.Redirect) error {
	query := `
		INSERT INTO redirect_history (
			request_log_id, original_url, redirect_url,
			redirect_type, redirect_status, redirect_timestamp
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		redirect.RequestLogID,
		redirect.OriginalURL,
		redirect.RedirectURL,
		redirect.RedirectType,
		redirect.RedirectStatus,
		redirect.RedirectTimestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to save redirect: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	redirect.ID = id
	return nil
} 