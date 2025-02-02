package models

type Order struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	OrderNumber   string `json:"order_number"`
	Status        string `json:"status"`
	LoyaltyPoints int    `json:"loyalty_points"`
	CreatedAt     string `json:"created_at"`
}
