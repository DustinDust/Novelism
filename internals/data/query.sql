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

-- name: UpdateUser :exec
UPDATE users
SET
    username=$2,
    password_hash=$3,
    email=$4,
    status=$5,
    verified=$6,
    verification_token=$7,
    password_reset_token=$8,
    first_name=$8,
    last_name=$9,
    date_of_birth=$9,
    gender=$10,
    profile_picture=$11
WHERE id=$1;


-- name: BulkInsertBooks :copyfrom
INSERT INTO books (
    user_id,
    title,
    description
) VALUES(
    $1, $2, $3
);

-- name: InsertBook :one
INSERT INTO books (
    user_id,
    title,
    description
) VALUES (
    $1, $2, $3
) RETURNING *;
