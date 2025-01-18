// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: add_last_fetched.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
SELECT feeds.id AS feed_id, feeds.url AS URL, feeds.name AS feed_name
FROM feeds
ORDER BY feeds.last_fetched_at NULLS FIRST
LIMIT 1
`

type GetNextFeedToFetchRow struct {
	FeedID   uuid.UUID
	Url      string
	FeedName string
}

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (GetNextFeedToFetchRow, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var i GetNextFeedToFetchRow
	err := row.Scan(&i.FeedID, &i.Url, &i.FeedName)
	return i, err
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE feeds.id = $1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}
