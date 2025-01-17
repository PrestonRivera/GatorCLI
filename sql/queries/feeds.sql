-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,   
    $6
)
RETURNING *;


-- name: ListFeeds :many
SELECT feeds.name AS feed_name, feeds.url AS url, users.name AS users_name
FROM feeds
INNER JOIN users ON feeds.user_id = users.id
ORDER BY feeds.updated_at DESC, feeds.created_at DESC;


-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;