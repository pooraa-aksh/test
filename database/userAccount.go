package database

//UserAccount : UserAccount Table structure as in DB
type UserAccount struct {
	tableName struct{} `sql:"user_account"`
	UserID    int      `json:"user_id" sql:"user_id,type:serial PRIMARY KEY"`
	Balance   float64  `json:"balance" sql:"balance,type:numeric NOT NULL DEFAULT 0.0"`
}
