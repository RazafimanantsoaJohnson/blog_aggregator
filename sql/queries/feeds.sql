-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id )
VALUES($1, $2, $3, $4, $5, $6 ) RETURNING *;

-- name: GetAllFeeds :many
SELECT feeds.name, feeds.url, users.name AS user_name FROM feeds INNER JOIN users ON feeds.user_id= users.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url= $1 LIMIT 1;

-- name: MarkFeedFetched :one
UPDATE feeds SET last_fetched= $2, updated_at= $3 WHERE id= $1 RETURNING * ;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched ASC NULLS FIRST LIMIT 1;
