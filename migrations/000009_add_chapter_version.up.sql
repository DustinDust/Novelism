CREATE TABLE IF NOT EXISTS chapter_versions (
    id SERIAL PRIMARY KEY,
    content_id INTEGER  REFERENCES contents(id) NOT NULL, -- refer to the lastest current chapter content 
    text_content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES users(id) NOT NULL
)
