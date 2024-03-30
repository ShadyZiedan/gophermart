package models

import "time"

type Order struct {
	Id          int       `json:"id"`
	UserId      int       `json:"userId"`
	Number      int       `json:"number"`
	Status      string    `json:"status"`
	Accrual     float64   `json:"accrual"`
	UploadedAt  time.Time `json:"uploadedAt"`
	ProcessedAt time.Time `json:"processedAt"`
}
