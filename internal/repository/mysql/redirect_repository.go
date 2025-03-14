package mysql

import (
	"database/sql"
	"fmt"
	"math/rand"
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

func (r *RedirectRepository) CreateRedirectMapping(clientID int64, mapping *models.RedirectMapping) error {
	// Generate a unique 6-character hash
	hash := r.generateUniqueHash()
	mapping.Hash = hash
	mapping.ClientID = clientID

	query := `
		INSERT INTO redirect_mappings (
			client_id, hash, redirect_url, redirect_url_black
		) VALUES (?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		mapping.ClientID,
		mapping.Hash,
		mapping.RedirectURL,
		mapping.RedirectURLBlack,
	)
	if err != nil {
		return fmt.Errorf("failed to create redirect mapping: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	mapping.ID = id
	return nil
}

func (r *RedirectRepository) GetClientRedirectMappings(clientID int64) ([]models.RedirectMapping, error) {
	query := `
		SELECT id, client_id, hash, redirect_url, redirect_url_black, created_at, updated_at
		FROM redirect_mappings
		WHERE client_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get redirect mappings: %w", err)
	}
	defer rows.Close()

	var mappings []models.RedirectMapping
	for rows.Next() {
		var mapping models.RedirectMapping
		err := rows.Scan(
			&mapping.ID,
			&mapping.ClientID,
			&mapping.Hash,
			&mapping.RedirectURL,
			&mapping.RedirectURLBlack,
			&mapping.CreatedAt,
			&mapping.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan redirect mapping: %w", err)
		}
		mappings = append(mappings, mapping)
	}

	return mappings, nil
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

// generateUniqueHash generates a unique 6-character hash
func (r *RedirectRepository) generateUniqueHash() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	for {
		hash := make([]byte, length)
		for i := range hash {
			hash[i] = charset[rand.Intn(len(charset))]
		}

		// Check if hash already exists
		var exists bool
		err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM redirect_mappings WHERE hash = ?)", string(hash)).Scan(&exists)
		if err != nil {
			// If there's an error, try again
			continue
		}

		if !exists {
			return string(hash)
		}
	}
} 