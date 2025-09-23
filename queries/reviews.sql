-- name: CreateReview :one
INSERT INTO reviews (game_id, user_id, rating, comment)
VALUES ($1, $2, $3, $4)
    ON CONFLICT (user_id, game_id) DO UPDATE
                                          SET rating = EXCLUDED.rating,
                                          comment = EXCLUDED.comment,
                                          updated_at = NOW()
                                          RETURNING *;

-- name: GetReview :one
SELECT * FROM reviews
WHERE review_id = $1;

-- name: ListReviewsByGame :many
SELECT r.review_id, r.rating, r.comment, r.created_at, u.username
FROM reviews r
         JOIN users u ON r.user_id = u.id
WHERE r.game_id = $1
ORDER BY r.created_at DESC;

-- name: ListReviewsByUser :many
SELECT r.review_id, r.rating, r.comment, r.created_at, g.title
FROM reviews r
         JOIN games g ON r.game_id = g.id
WHERE r.user_id = $1
ORDER BY r.created_at DESC;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE review_id = $1;
