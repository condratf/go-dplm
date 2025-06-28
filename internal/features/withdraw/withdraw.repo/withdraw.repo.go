package withdrawrepo

import (
	"database/sql"
	"time"

	"github.com/condratf/go-musthave-diploma-tpl/internal/custerrors"
	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
)

func NewWithdrawRepository(
	db *sql.DB,
) WithdrawRepository {
	return &withdrawRepository{
		db: db,
	}
}

func (r *withdrawRepository) Withdraw(login, order string, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	var available float64

	err = tx.QueryRow(`
		SELECT u.id, ub.available_points 
		FROM users u
		JOIN user_balance ub ON u.id = ub.user_id
		WHERE u.login = $1
	`, login).Scan(&userID, &available)
	if err != nil {
		return err
	}

	if available < amount {
		return custerrors.ErrInsufficientFunds
	}

	_, err = tx.Exec(`
		UPDATE user_balance 
		SET available_points = available_points - $1 
		WHERE user_id = $2
	`, amount, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO withdrawals (user_id, amount, order_number, created_at) 
		VALUES ($1, $2, $3, NOW())
	`, userID, amount, order)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *withdrawRepository) GetWithdrawals(login string) ([]models.WithdrawalRes, error) {
	query := `
		SELECT w.order_number, w.amount, w.created_at
		FROM users u
		JOIN withdrawals w ON u.id = w.user_id
		WHERE u.login = $1
		ORDER BY w.created_at DESC
	`
	rows, err := r.db.Query(query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []models.WithdrawalRes
	for rows.Next() {
		var w models.WithdrawalRes
		var createdAt time.Time
		if err := rows.Scan(&w.Order, &w.Sum, &createdAt); err != nil {
			return nil, err
		}
		w.ProcessedAt = createdAt.Format(time.RFC3339)
		withdrawals = append(withdrawals, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return withdrawals, nil
}
