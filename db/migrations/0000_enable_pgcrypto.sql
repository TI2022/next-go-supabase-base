-- 0000_enable_pgcrypto.sql
-- UUID 生成関数 gen_random_uuid() を使うための拡張

CREATE EXTENSION IF NOT EXISTS pgcrypto;

