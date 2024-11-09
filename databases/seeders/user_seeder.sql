-- users
INSERT INTO public.users
(id, email, password_hash, is_verified, created_at, updated_at, deleted_at)
VALUES(1, 'some.one@mail.com', '$2a$10$UJ/8vfbswtXvPBLz6q9Wye54jQuEGBTGY4YPLCG1t0TFj8F2DtnEm', false, '2024-11-09 19:26:33.826', '2024-11-09 19:26:33.826', NULL);
INSERT INTO public.users
(id, email, password_hash, is_verified, created_at, updated_at, deleted_at)
VALUES(2, 'some.one2@mail.com', '$2a$10$8aP.KWO0u3XYp.1wcyphPOgw2XS6BdOBZJwS7EcF2y1gpl2FigmxW', false, '2024-11-09 19:26:41.220', '2024-11-09 19:26:41.220', NULL);
INSERT INTO public.users
(id, email, password_hash, is_verified, created_at, updated_at, deleted_at)
VALUES(3, 'some.one3@mail.com', '$2a$10$4weSfeEx9J64.CA5HVp9Aux3eAhRE0b7xkh.CzWrMklUtXI2MrUO6', false, '2024-11-09 19:26:48.425', '2024-11-09 19:26:48.425', NULL);
INSERT INTO public.users
(id, email, password_hash, is_verified, created_at, updated_at, deleted_at)
VALUES(4, 'some.one4@mail.com', '$2a$10$fQyLZgL.tUqRgiIBqa0vd.cH9kAn81Cm0N7gePRcpHKZnjvDCIrP2', false, '2024-11-09 19:26:58.510', '2024-11-09 19:26:58.510', NULL);
INSERT INTO public.users
(id, email, password_hash, is_verified, created_at, updated_at, deleted_at)
VALUES(5, 'some.one5@mail.com', '$2a$10$rLJ8IxnYgQQH5Gplbybgeu7X8kEmu8FJIuyxB5Fippfu3AlTi0Wce', false, '2024-11-09 19:27:06.613', '2024-11-09 19:27:06.613', NULL);

-- profiles
INSERT INTO public.profiles
(id, user_id, "name", age, bio, gender, "location", interests, photos, is_premium, created_at, updated_at, deleted_at)
VALUES(2, 2, 'Faiq Najib Kedua', 24, '', 'male'::public."gender_enum", 'Jalan Aja Dahulu No.212, Surabaya, Jawa Timur', NULL, NULL, false, '2024-11-09 19:26:41.220', '2024-11-09 19:26:41.220', NULL);
INSERT INTO public.profiles
(id, user_id, "name", age, bio, gender, "location", interests, photos, is_premium, created_at, updated_at, deleted_at)
VALUES(3, 3, 'Faiq Najib Ketiga', 24, '', 'male'::public."gender_enum", 'Jalan Aja Dahulu No.212, Surabaya, Jawa Timur', NULL, NULL, false, '2024-11-09 19:26:48.425', '2024-11-09 19:26:48.425', NULL);
INSERT INTO public.profiles
(id, user_id, "name", age, bio, gender, "location", interests, photos, is_premium, created_at, updated_at, deleted_at)
VALUES(4, 4, 'Faiq Najib Keempat', 24, '', 'male'::public."gender_enum", 'Jalan Aja Dahulu No.212, Surabaya, Jawa Timur', NULL, NULL, false, '2024-11-09 19:26:58.510', '2024-11-09 19:26:58.510', NULL);
INSERT INTO public.profiles
(id, user_id, "name", age, bio, gender, "location", interests, photos, is_premium, created_at, updated_at, deleted_at)
VALUES(5, 5, 'Faiq Najib Kelima', 24, '', 'male'::public."gender_enum", 'Jalan Aja Dahulu No.212, Surabaya, Jawa Timur', NULL, NULL, false, '2024-11-09 19:27:06.613', '2024-11-09 19:27:06.613', NULL);
INSERT INTO public.profiles
(id, user_id, "name", age, bio, gender, "location", interests, photos, is_premium, created_at, updated_at, deleted_at)
VALUES(1, 1, 'Faiq Najib Pertama', 24, '', 'male'::public."gender_enum", 'Jalan Aja Dahulu No.212, Surabaya, Jawa Timur', NULL, NULL, true, '2024-11-09 19:26:33.829', '2024-11-09 19:27:39.008', NULL);