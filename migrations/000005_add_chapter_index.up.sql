CREATE UNIQUE INDEX idx_unique_chapter_chapter_no_book_id
ON chapters (chapter_no, book_id);
