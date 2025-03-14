package models

import "time"

type Redirect struct {
	ID               int64     `json:"id"`
	RequestLogID     int64     `json:"request_log_id"`
	OriginalURL      string    `json:"original_url"`
	RedirectURL      string    `json:"redirect_url"`
	RedirectType     string    `json:"redirect_type"`
	RedirectStatus   int       `json:"redirect_status"`
	RedirectTimestamp time.Time `json:"redirect_timestamp"`
}

type RedirectType string

const (
	RedirectTypeSite1 RedirectType = "site1"
	RedirectTypeSite2 RedirectType = "site2"
	RedirectTypeSite3 RedirectType = "site3"
)

type RedirectMapping struct {
	ID          int64     `json:"id"`
	Hash        string    `json:"hash"`
	RedirectURL string    `json:"redirect_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
} 