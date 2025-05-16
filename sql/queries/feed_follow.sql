-- name: CreateFeedFollow :one
WITH inserted AS
    (INSERT INTO feed_follows(id,created_at,updated_at,user_id,feed_id) VALUES(
        $1,
        NOW(),
        NOW(),
        $2,
        $3
    )RETURNING *
),
inserted_user AS(
    SELECT inserted.*,users.name AS user_name FROM
    inserted INNER JOIN
    users ON inserted.user_id = users.id
)
SELECT inserted_user.*,feeds.name AS feed_name FROM inserted_user
INNER JOIN feeds ON inserted_user.feed_id = feeds.id;

-- name: GetFollowedFeeds :many
WITH follows AS(
    SELECT * FROM feed_follows WHERE feed_follows.user_id = $1
)
SELECT feeds.* FROM follows INNER JOIN feeds ON follows.feed_id = feeds.id;

-- name: DeleteFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;
