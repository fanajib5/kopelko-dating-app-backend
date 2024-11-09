-- Create enum types
CREATE TYPE gender_enum AS ENUM ('male', 'female', 'other');
CREATE TYPE swipe_type_enum AS ENUM ('left', 'right');
CREATE TYPE payment_status_enum AS ENUM ('completed', 'pending');
CREATE TYPE payment_provider_enum AS ENUM ('stripe', 'xendit');

-- Create tables
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    age INT CHECK (age >= 18 AND age > 0),  -- Assuming the app has an age limit
    bio TEXT,
    gender gender_enum,
    location VARCHAR(100),
    interests TEXT,  -- Comma-separated list of interests
    photos TEXT[],  -- Array of photo URLs
    is_premium BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE swipes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    swipe_type swipe_type_enum,
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

CREATE TABLE premium_features (
    id SERIAL PRIMARY KEY,
    feature_name VARCHAR(50) UNIQUE NOT NULL,  -- Examples: 'no_swipe_quota', 'verified_label'
    description TEXT
);

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feature_id INT NOT NULL REFERENCES premium_features(id),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    auto_renew BOOLEAN DEFAULT TRUE
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(10) DEFAULT 'USD',
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payment_status payment_status_enum,
    payment_provider payment_provider_enum
);

-- Indexes for faster searches by user ID
CREATE INDEX idx_user_id ON profiles(user_id);
CREATE INDEX idx_swipe_user_id ON swipes(user_id);
CREATE INDEX idx_match_user1_id ON matches(user1_id);
CREATE INDEX idx_match_user2_id ON matches(user2_id);
CREATE INDEX idx_subscription_user_id ON subscriptions(user_id);
CREATE INDEX idx_payment_user_id ON payments(user_id);
