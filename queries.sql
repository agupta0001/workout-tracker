-- name: get-workouts
SELECT * FROM workouts;

-- name: get-exercises
SELECT * FROM exercises;

-- name: get-exercise-by-id
SELECT * FROM exercises WHERE id = $1;

-- name: get-user-by-id
SELECT * FROM users WHERE id = $1;

-- name: get-user-by-uuid
SELECT * FROM users WHERE uid = $1;

-- name: create-user
INSERT INTO users (uid, email, name, avatar) VALUES ($1, $2, $3, $4) RETURNING *;