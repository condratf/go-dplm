package storage

import (
	"database/sql"
	"fmt"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (Storage, error) {
	store := &PostgresStore{db: db}

	query := `
		CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(100) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		order_number VARCHAR(50) UNIQUE NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		-- pending, processed, rejected
		loyalty_points INT DEFAULT 0,
		-- calculated for the order
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE user_balance (
		user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
		total_points INT NOT NULL DEFAULT 0,
		available_points INT NOT NULL DEFAULT 0
	);

	CREATE TABLE withdrawals (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		amount INT NOT NULL,
		-- withdrawn
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE INDEX idx_user_orders ON orders(user_id);

	CREATE INDEX idx_user_withdrawals ON withdrawals(user_id);
	`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return store, nil
}
