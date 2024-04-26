ALTER TABLE users
ADD COLUMN refresh_token VARCHAR(255),
ADD COLUMN  refresh_token_valid_until TIMESTAMP;