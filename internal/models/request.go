package models

import "time"

type Request struct {
	ID              int64     `json:"id"`
	Timestamp       time.Time `json:"timestamp"`
	IPAddress       string    `json:"ip_address"`
	RequestURL      string    `json:"request_url"`
	RequestMethod   string    `json:"request_method"`
	RequestHeaders  []byte    `json:"request_headers"`
	ProcessingStatus string    `json:"processing_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
} 