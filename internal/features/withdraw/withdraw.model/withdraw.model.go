package withdrawmodel

type Withdrawal struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
}

type WithdrawalRes struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
