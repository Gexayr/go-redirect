package mysql

import (
	"database/sql"
	"fmt"
	"platform/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (r *ClientRepository) Create(client *models.Client) error {
	query := `
		INSERT INTO clients (username, password_hash, email)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		client.Username,
		client.PasswordHash,
		client.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	client.ID = id
	return nil
}

func (r *ClientRepository) GetByUsername(username string) (*models.Client, error) {
	query := `
		SELECT id, username, password_hash, email, created_at, updated_at
		FROM clients
		WHERE username = ?
	`

	client := &models.Client{}
	err := r.db.QueryRow(query, username).Scan(
		&client.ID,
		&client.Username,
		&client.PasswordHash,
		&client.Email,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	return client, nil
}

func (r *ClientRepository) GetByID(id int64) (*models.Client, error) {
	query := `
		SELECT id, username, password_hash, email, created_at, updated_at
		FROM clients
		WHERE id = ?
	`

	client := &models.Client{}
	err := r.db.QueryRow(query, id).Scan(
		&client.ID,
		&client.Username,
		&client.PasswordHash,
		&client.Email,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	return client, nil
}

func (r *ClientRepository) ValidatePassword(client *models.Client, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(client.PasswordHash), []byte(password))
	return err == nil
} 