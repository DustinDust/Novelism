CREATE TYPE content_status AS ENUM('draft', 'published');

ALTER TABLE IF EXISTS contents
ADD COLUMN status content_status DEFAULT 'draft';
