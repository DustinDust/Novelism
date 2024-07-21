ALTER TABLE books
ADD COLUMN deleted_at TIMESTAMP;

CREATE TYPE user_status AS ENUM ('active', 'idle', 'deleted');

ALTER TABLE users ADD COLUMN status USER_STATUS;
