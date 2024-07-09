CREATE TABLE IF NOT EXISTS books (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id),
	title TEXT,
	description TEXT,
  	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP
	-- genres
	-- tags
	-- ratings
)
