CREATE TABLE IF NOT EXISTS chapters (
	id SERIAL PRIMARY KEY,
	book_id INT REFERENCES books(id),
	author_id INT REFERENCES users(id), -- user id
	title TEXT,
	description TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contents (
	id SERIAL PRIMARY KEY,
	chapter_id INT UNIQUE REFERENCES chapters(id),
	text_content TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP
)
