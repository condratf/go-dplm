package models

type Withdrawal struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
}
