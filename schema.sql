-- users
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
	id 			SERIAL PRIMARY KEY,
	uuid 		uuid NOT NULL UNIQUE,
	name 		VARCHAR(255) NOT NULL,
	email 		VARCHAR(255) NOT NULL UNIQUE,
	avatar 		text,
	token 		text,
	
	created_at 	TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at 	TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Exercises
DROP TABLE IF EXISTS exercises CASCADE;
CREATE TABLE exercises (
	id 			SERIAL PRIMARY KEY,
	name 		VARCHAR(255) NOT NULL UNIQUE,
	tags		VARCHAR(100)[] DEFAULT '{}',
	
	created_at 	TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at 	TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Workouts
DROP TABLE IF EXISTS workouts CASCADE;
CREATE TABLE workouts (
    id 			SERIAL PRIMARY KEY,
    set_no 		INTEGER NOT NULL,
    reps 		INTEGER NOT NULL,
	weight 		DECIMAL NOT NULL,

    exercise_id INTEGER NOT NULL REFERENCES exercises(id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id 	INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    
    created_at 	TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at 	TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
DROP INDEX IF EXISTS idx_workout_date; CREATE INDEX idx_workout_date ON workouts ((TIMEZONE('UTC', created_at)::DATE));
DROP INDEX IF EXISTS idx_workout_user; CREATE INDEX idx_workout_user ON workouts (user_id);
DROP INDEX IF EXISTS idx_workout_exercise; CREATE INDEX idx_workout_exercise ON workouts (exercise_id);
DROP INDEX IF EXISTS idx_workout_user_exercise; CREATE INDEX idx_workout_user_exercise ON workouts (user_id, exercise_id);