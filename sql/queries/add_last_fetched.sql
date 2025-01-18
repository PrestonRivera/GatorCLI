-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE feeds.id = $1;


-- name: GetNextFeedToFetch :one
SELECT feeds.id AS feed_id, feeds.url AS URL, feeds.name AS feed_name
FROM feeds
ORDER BY feeds.last_fetched_at NULLS FIRST
LIMIT 1;

