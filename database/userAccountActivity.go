package database

import "time"

//UserAccountActivity : User Account Activity Table structure as in DB
type UserAccountActivity struct {
	tableName       struct{}  `sql:"user_account_activity"`
	TransactionID   int       `json:"transaction_id" sql:"transaction_id,type:serial PRIMARY KEY"`
	FkUserID        int       `json:"user_id" sql:"fk_user_id,type:int NOT NULL REFERENCES user_account(user_id) ON DELETE RESTRICT ON UPDATE CASCADE"`
	Priority        int       `json:"priority" sql:"priority,type:int DEFAULT NULL"`
	Type            string    `json:"type" sql:"type,type:varchar(255) DEFAULT NULL"`
	TransactionType string    `json:"transaction_type" sql:"transaction_type,type:varchar(255) DEFAULT NULL"`
	Status          string    `json:"status" sql:"status,type:varchar(255) NOT NULL DEFAULT 'active'"`
	Amount          float64   `json:"amount" sql:"amount,type:numeric NOT NULL DEFAULT 0.0"`
	BalanceAmt      float64   `json:"balance_amt" sql:"balance_amt,type:numeric NOT NULL DEFAULT 0.0"`
	ExpDate         time.Time `json:"exp_date" sql:"exp_date,type:timestamp DEFAULT NULL"`
	CreatedAt       time.Time `json:"-" sql:"created_at,type:timestamp NOT NULL DEFAULT NOW()"`
}
