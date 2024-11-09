-- seeder
INSERT INTO premium_features (feature_name, description) VALUES
('no_swipe_quota', 'Unlimited swipes per day'),
('verified_label', 'Verified badge on profile');

INSERT INTO public.users
(id, email, password_hash, is_verified, created_at, updated_at, deleted_at)
VALUES(1, 'some.one@mail.com', '$2a$10$zUyiG.8LbNczQFBXck/tK.pbWui5j2MVHT1zki7tEojwNiVnhZ3Qq', false, '2024-11-08 23:15:56.861', '2024-11-08 23:15:56.861', NULL);
INSERT INTO public.users
(id, email, password_hash, is_verified, created_at, updated_at, deleted_at)
VALUES(3, 'some.one1@mail.com', '$2a$10$055y/Ucb2Cn85TGlqUAVV.lee5Kn564v6tuiedzXXX48RRbIdrIPe', false, '2024-11-08 23:20:52.884', '2024-11-08 23:20:52.884', NULL);
INSERT INTO public.users
(id, email, password_hash, is_verified, created_at, updated_at, deleted_at)
VALUES(4, 'some.one2@mail.com', '$2a$10$OMg0Zw19IF/MI8928E0yveSMc4KcRKJm3XQEMTxe6luKXecgUd6AK', false, '2024-11-08 23:24:46.123', '2024-11-08 23:24:46.123', NULL);

INSERT INTO public.swipes
(id, user_id, target_user_id, swipe_type, swipe_date)
VALUES(1, 1, 3, 'left', '2024-11-08');
INSERT INTO public.swipes
(id, user_id, target_user_id, swipe_type, swipe_date)
VALUES(2, 1, 4, 'right', '2024-11-08');

INSERT INTO public.profiles
(id, user_id, "name", age, bio, gender, "location", interests, photos, is_premium, created_at, updated_at, deleted_at)
VALUES(1, 1, 'Faiq Najib', 20, '', 'male', 'Jalan Aja Dahulu No.212, Surabaya, Jawa Timur', NULL, NULL, false, '2024-11-08 23:15:56.864', '2024-11-08 23:15:56.864', NULL);
INSERT INTO public.profiles
(id, user_id, "name", age, bio, gender, "location", interests, photos, is_premium, created_at, updated_at, deleted_at)
VALUES(2, 3, 'Faiq Najib', 20, '', 'male', 'Jalan Aja Dahulu No.212, Surabaya, Jawa Timur', NULL, NULL, false, '2024-11-08 23:20:52.886', '2024-11-08 23:20:52.886', NULL);
INSERT INTO public.profiles
(id, user_id, "name", age, bio, gender, "location", interests, photos, is_premium, created_at, updated_at, deleted_at)
VALUES(3, 4, 'Faiq Najib', 20, '', 'male', 'Jalan Aja Dahulu No.212, Surabaya, Jawa Timur', NULL, NULL, false, '2024-11-08 23:24:46.125', '2024-11-08 23:24:46.125', NULL);