ALTER TABLE books
ADD COLUMN deleted_at TIMESTAMP;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
        CREATE TYPE user_status AS ENUM ('active', 'idle', 'deleted');
    END IF;
END $$;

ALTER TABLE users ADD COLUMN status USER_STATUS
