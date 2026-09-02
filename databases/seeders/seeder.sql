-- Premium Features Seeder
INSERT INTO premium_features (id, feature_name, description) VALUES
(1, 'no_swipe_quota', 'Unlimited swipes per day'),
(2, 'verified_label', 'Verified badge on profile')
ON CONFLICT (feature_name) DO NOTHING;

-- Reset sequence for premium_features
SELECT setval('premium_features_id_seq', (SELECT COALESCE(MAX(id), 1) FROM premium_features));

-- Dummy Users Seeder
INSERT INTO users (id, email, password_hash, is_verified, created_at, updated_at) VALUES
(1, 'alice@example.com', '$2a$10$UJ/8vfbswtXvPBLz6q9Wye54jQuEGBTGY4YPLCG1t0TFj8F2DtnEm', true, NOW(), NOW()),
(2, 'bob@example.com', '$2a$10$8aP.KWO0u3XYp.1wcyphPOgw2XS6BdOBZJwS7EcF2y1gpl2FigmxW', false, NOW(), NOW()),
(3, 'charlie@example.com', '$2a$10$4weSfeEx9J64.CA5HVp9Aux3eAhRE0b7xkh.CzWrMklUtXI2MrUO6', false, NOW(), NOW()),
(4, 'diana@example.com', '$2a$10$fQyLZgL.tUqRgiIBqa0vd.cH9kAn81Cm0N7gePRcpHKZnjvDCIrP2', true, NOW(), NOW()),
(5, 'eve@example.com', '$2a$10$rLJ8IxnYgQQH5Gplbybgeu7X8kEmu8FJIuyxB5Fippfu3AlTi0Wce', false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Reset sequence for users
SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 1) FROM users));

-- Dummy Profiles Seeder
INSERT INTO profiles (id, user_id, name, age, bio, gender, location, interests, photos, is_premium, created_at, updated_at) VALUES
(1, 1, 'Alice Wonderland', 24, 'Coffee and adventure lover', 'female', 'Surabaya', ARRAY['coffee', 'traveling'], ARRAY['https://images.unsplash.com/photo-1494790108377-be9c29b29330'], true, NOW(), NOW()),
(2, 2, 'Bob Builder', 26, 'Passionate software engineer', 'male', 'Jakarta', ARRAY['coding', 'gadgets'], ARRAY['https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d'], false, NOW(), NOW()),
(3, 3, 'Charlie Parker', 23, 'Jazz and good books', 'male', 'Bandung', ARRAY['music', 'reading'], ARRAY['https://images.unsplash.com/photo-1500648767791-00dcc994a43e'], false, NOW(), NOW()),
(4, 4, 'Diana Prince', 27, 'Fitness and art enthusiast', 'female', 'Bali', ARRAY['gym', 'art'], ARRAY['https://images.unsplash.com/photo-1534528741775-53994a69daeb'], true, NOW(), NOW()),
(5, 5, 'Eve Polastri', 25, 'Foodie and weekend runner', 'female', 'Surabaya', ARRAY['food', 'running'], ARRAY['https://images.unsplash.com/photo-1517841905240-472988babdf9'], false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Reset sequence for profiles
SELECT setval('profiles_id_seq', (SELECT COALESCE(MAX(id), 1) FROM profiles));
