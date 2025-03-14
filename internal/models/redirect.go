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
	CreatedAt        time.Time `json:"created_at"`
}

type RedirectMapping struct {
	ID              int64     `json:"id"`
	ClientID        int64     `json:"client_id"`
	Hash            string    `json:"hash"`
	RedirectURL     string    `json:"redirect_url"`
	RedirectURLBlack string    `json:"redirect_url_black"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type RedirectMappingCreate struct {
	RedirectURL     string `json:"redirect_url" binding:"required,url"`
	RedirectURLBlack string `json:"redirect_url_black" binding:"required,url"`
} 