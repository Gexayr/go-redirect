package mysql

import (
	"database/sql"
	"fmt"
	"platform/internal/models"
)

type RequestRepository struct {
	db *sql.DB
}

func NewRequestRepository(db *sql.DB) *RequestRepository {
	return &RequestRepository{
		db: db,
	}
}

func (r *RequestRepository) SaveRequest(request *models.Request) error {
	query := `
		INSERT INTO request_logs (
			timestamp, ip_address, request_url, request_method,
			request_headers, processing_status
		) VALUES (?, ?, ?, ?, ?, 'processed')
	`

	result, err := r.db.Exec(
		query,
		request.Timestamp,
		request.IPAddress,
		request.RequestURL,
		request.RequestMethod,
		request.RequestHeaders,
	)
	if err != nil {
		return fmt.Errorf("failed to save request: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	request.ID = id
	return nil
} 