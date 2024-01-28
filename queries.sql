-- name: get-workouts
SELECT * FROM workouts;

-- name: get-exercises
SELECT * FROM exercises;

-- name: get-exercise-by-id
SELECT * FROM exercises WHERE id = $1;
