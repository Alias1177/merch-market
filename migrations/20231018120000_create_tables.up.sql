CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) UNIQUE NOT NULL,
                                     password_hash TEXT NOT NULL,
                                     coins INT DEFAULT 0 CHECK (coins >= 0),
                                     created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username_lower ON users(LOWER(username));

CREATE TABLE IF NOT EXISTS items (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(255) UNIQUE NOT NULL,
                                     price INT NOT NULL CHECK (price > 0)
);

CREATE INDEX IF NOT EXISTS idx_items_price ON items(price);

INSERT INTO items (name, price) VALUES
                                    ('t-shirt', 80),
                                    ('cup', 20),
                                    ('book', 50),
                                    ('pen', 10),
                                    ('powerbank', 200),
                                    ('hoody', 300),
                                    ('umbrella', 200),
                                    ('socks', 10),
                                    ('wallet', 50),
                                    ('pink-hoody', 500);

CREATE TABLE IF NOT EXISTS inventory (
                                         id SERIAL PRIMARY KEY,
                                         user_id INT REFERENCES users(id) ON DELETE CASCADE,
                                         item_id INT REFERENCES items(id) ON DELETE CASCADE,
                                         quantity INT DEFAULT 1 CHECK (quantity > 0),
                                         UNIQUE(user_id, item_id)
);

-- Индексы для ускорения поиска по инвентарю
CREATE INDEX IF NOT EXISTS idx_inventory_user_id ON inventory(user_id);
CREATE INDEX IF NOT EXISTS idx_inventory_item_id ON inventory(item_id);

CREATE TABLE IF NOT EXISTS transactions (
                                            id SERIAL PRIMARY KEY,
                                            sender_id INT REFERENCES users(id) ON DELETE SET NULL,
                                            receiver_id INT REFERENCES users(id) ON DELETE SET NULL,
                                            amount INT NOT NULL CHECK (amount > 0),
                                            created_at TIMESTAMP DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_transactions_sender_id ON transactions(sender_id);
CREATE INDEX IF NOT EXISTS idx_transactions_receiver_id ON transactions(receiver_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);