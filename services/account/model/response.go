package model

import "time"

//CreditLogList : credit log list response
type CreditLogList struct {
	CreditLog []CreditLog `json:"log"`
}

//CreditLog : credit log
type CreditLog struct {
	UserID     int       `json:"userID" sql:"fk_user_id"`
	Amount     float64   `json:"amt" sql:"amount"`
	Type       string    `json:"type" sql:"type"`
	Priority   int       `json:"priority" sql:"priority"`
	ExpiryDate time.Time `json:"exp_date" sql:"exp_date"`
	Status     string    `json:"status" sql:"status"`
}
