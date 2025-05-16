-- name: CreateFeed :one
INSERT INTO feeds(id,created_at,updated_at,name,url,user_id)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = NOW(),last_fetched_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
WITH follows AS(
    SELECT feed_id FROM feed_follows WHERE 
    feed_follows.user_id = $1
)
SELECT feeds.* FROM feeds
INNER JOIN follows ON feeds.id = follows.feed_id
ORDER BY last_fetched_at NULLS FIRST LIMIT 1;
