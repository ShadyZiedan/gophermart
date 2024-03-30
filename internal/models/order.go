package models

import "time"

type Order struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Number      int       `json:"number"`
	Status      string    `json:"status"`
	Accrual     float64   `json:"accrual"`
	UploadedAt  time.Time `json:"uploaded_at"`
	ProcessedAt time.Time `json:"processed_at"`
}
