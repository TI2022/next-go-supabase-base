-- 0002_seed_dummy_user.sql
-- users テーブルにダミーユーザーを1件投入するマイグレーション

INSERT INTO users (email, password_hash, name)
VALUES (
    'demo@example.com',
    -- 実運用では bcrypt などでハッシュ化した値に置き換えること
    'plain-text-demo-password',
    'Demo User'
)
ON CONFLICT (email) DO NOTHING;

