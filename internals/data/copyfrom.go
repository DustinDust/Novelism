// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: copyfrom.go

package data

import (
	"context"
)

// iteratorForBulkInsertBooks implements pgx.CopyFromSource.
type iteratorForBulkInsertBooks struct {
	rows                 []BulkInsertBooksParams
	skippedFirstNextCall bool
}

func (r *iteratorForBulkInsertBooks) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForBulkInsertBooks) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].UserID,
		r.rows[0].Title,
		r.rows[0].Cover,
		r.rows[0].Description,
		r.rows[0].Visibility,
	}, nil
}

func (r iteratorForBulkInsertBooks) Err() error {
	return nil
}

func (q *Queries) BulkInsertBooks(ctx context.Context, arg []BulkInsertBooksParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"books"}, []string{"user_id", "title", "cover", "description", "visibility"}, &iteratorForBulkInsertBooks{rows: arg})
}
