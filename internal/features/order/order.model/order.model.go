package ordermodel

type Order struct {
	ID            int    `json:"id"`
	UserLogin     int    `json:"user_login"`
	Order         string `json:"order"`
	Status        string `json:"status"`
	LoyaltyPoints int    `json:"loyalty_points"`
	CreatedAt     string `json:"created_at"`
}
