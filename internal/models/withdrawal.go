package models

import "time"

type Withdrawal struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	OrderNumber int       `json:"order_number"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
