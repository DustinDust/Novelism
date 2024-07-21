
CREATE TYPE visibility AS ENUM ('hidden', 'visible');
ALTER TABLE books
ADD COLUMN visibility VISIBILITY DEFAULT 'visible';
