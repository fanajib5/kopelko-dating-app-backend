CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    age INT CHECK (age >= 18),  -- Assuming the app has an age limit
    bio TEXT,
    gender VARCHAR(20),
    location VARCHAR(100),
    interests TEXT,  -- Comma-separated list of interests
    photos TEXT[],  -- Array of photo URLs
    is_premium BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE swipes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    swipe_type VARCHAR(10) CHECK (swipe_type IN ('left', 'right')),
    swipe_date DATE NOT NULL DEFAULT CURRENT_DATE,
    UNIQUE(user_id, target_user_id, swipe_date)  -- Prevents duplicate swipes on the same day
);

CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    user1_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    matched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user1_id, user2_id)  -- Prevents duplicate matches
);

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_type VARCHAR(50),  -- e.g., "premium", "standard"
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    auto_renew BOOLEAN DEFAULT TRUE
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payment_status VARCHAR(50),  -- e.g., "completed", "pending"
    payment_provider VARCHAR(50)  -- e.g., "stripe", "paypal"
);

-- Indexes for faster searches by user ID
CREATE INDEX idx_user_id ON profiles(user_id);
CREATE INDEX idx_swipe_user_id ON swipes(user_id);
CREATE INDEX idx_match_user1_id ON matches(user1_id);
CREATE INDEX idx_match_user2_id ON matches(user2_id);
CREATE INDEX idx_subscription_user_id ON subscriptions(user_id);
CREATE INDEX idx_payment_user_id ON payments(user_id);

