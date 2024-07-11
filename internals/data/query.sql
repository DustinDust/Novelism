-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: InsertUser :one
INSERT INTO users (
    username,
    password_hash,
    email,
    status,
    verified,
    verification_token,
    password_reset_token,
    first_name,
    last_name,
    date_of_birth,
    gender,
    profile_picture
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: BulkInsertBooks :copyfrom
INSERT INTO books (
    user_id,
    title,
    description
) VALUES(
    $1, $2, $3
);
