-- 0003_update_dummy_user_password_hash.sql
-- ダミーユーザーの password_hash を bcrypt ハッシュに更新する

UPDATE users
SET password_hash = crypt('plain-text-demo-password', gen_salt('bf'))
WHERE email = 'demo@example.com';

