package model

//AccountActivityReq : account activity request
type AccountActivityReq struct {
	ActivityType string         `json:"activity" binding:"required"`
	Payload      ActivityDetail `json:"payload" binding:"dive"`
}

//ActivityDetail: account activity detail
type ActivityDetail struct {
	UserID   int     `json:"userId"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
	Priority int     `json:"priority"`
	Expiry   int     `json:"expiry"`
}
