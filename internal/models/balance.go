package models

type Balance struct {
	ID              int `json:"id"`
	UserID          int `json:"user_id"`
	TotalPoints     int `json:"total_points"`
	AvailablePoints int `json:"available_points"`
}
