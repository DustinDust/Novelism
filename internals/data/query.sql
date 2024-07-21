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
    first_name=$9,
    last_name=$10,
    date_of_birth=$11,
    gender=$12,
    profile_picture=$13
WHERE id=$1;


-- name: BulkInsertBooks :copyfrom
INSERT INTO books (
    user_id,
    title,
    cover,
    description,
    visibility
) VALUES(
    $1, $2, $3, $4, $5
);

-- name: InsertBook :one
INSERT INTO books (
    user_id,
    title,
    description,
    cover,
    visibility
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetBookById :one
SELECT * FROM books WHERE id=$1 LIMIT 1;

-- name: BrowseBooks :many
SELECT * FROM books WHERE visibility = 'visible' LIMIT $1 OFFSET $2;

-- name: CountBrowsableBooks :one
SELECT count(*) FROM BOOKS WHERE visibility = 'visible';

-- name: FindBooksByUserId :many
SELECT * FROM books WHERE user_id = $1;

-- name: UpdateBook :exec
UPDATE books
SET
    title = $2,
    description = $3,
    updated_at = $4,
    cover = $5,
    visibility = $6
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteBook :exec
UPDATE books
SET
    deleted_at = now()
WHERE id = $1;

-- name: FindChaptersByBookId :many
SELECT * FROM chapters WHERE chapters.book_id = $1 AND deleted_at IS NULL;

-- name: InsertChapter :one
INSERT INTO chapters (
    book_id, author_id, title, description
) VALUES (
    $1, $2, $3, $4
) RETURNING *;
