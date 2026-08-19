CREATE TABLE IF NOT EXISTS users_subscriptions (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    group_id VARCHAR(50) NOT NULL,
    PRIMARY KEY (user_id, group_id)
);