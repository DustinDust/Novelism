-- Create a function to update the timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to automatically update the `updated_at` field
CREATE TRIGGER trigger_update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Create a trigger to automatically update the `updated_at` field
CREATE TRIGGER trigger_update_books_updated_at
BEFORE UPDATE ON books
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER trigger_update_chapters_updated_at
BEFORE UPDATE ON chapters
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER trigger_update_contents_updated_at
BEFORE UPDATE ON contents
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER trigger_update_content_versions_updated_at
BEFORE UPDATE ON content_versions
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();
