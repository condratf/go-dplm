CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  login VARCHAR(100) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  user_login INT NOT NULL REFERENCES users(login) ON DELETE CASCADE,
  order VARCHAR(50) UNIQUE NOT NULL,
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

CREATE INDEX idx_user_orders ON orders(user_login);

CREATE INDEX idx_user_withdrawals ON withdrawals(user_id);