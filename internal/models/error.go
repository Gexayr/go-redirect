package models

import "time"

type ErrorLog struct {
	ID           int64     `json:"id"`
	RequestLogID int64     `json:"request_log_id"`
	ErrorMessage string    `json:"error_message"`
	ErrorType    string    `json:"error_type"`
	StackTrace   string    `json:"stack_trace"`
	CreatedAt    time.Time `json:"created_at"`
} 