-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT p.*, f.name AS feed_name FROM posts p
JOIN feeds f ON f.id = p.feed_id
JOIN feed_follows f_f ON f_f.feed_id = f.id
WHERE f_f.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;
