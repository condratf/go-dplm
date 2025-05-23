package models

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type Order struct {
	ID            int    `json:"id"`
	UserLogin     int    `json:"user_login"`
	Order         string `json:"order"`
	Status        string `json:"status"`
	LoyaltyPoints int    `json:"loyalty_points"`
	CreatedAt     string `json:"created_at"`
}

type Balance struct {
	ID              int `json:"id"`
	UserID          int `json:"user_id"`
	TotalPoints     int `json:"total_points"`
	AvailablePoints int `json:"available_points"`
}

type Withdrawal struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
}

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawalRes struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type AccrualResponse struct {
	Order   string
	Status  string
	Accrual float64
}
