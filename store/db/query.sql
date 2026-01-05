-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByOAuthSub :one
SELECT * FROM users
WHERE oauth_sub = $1 LIMIT 1;

-- name: ListPendingUsers :many
SELECT * FROM users
WHERE status = 'pending'
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (
    email, name, oauth_provider, oauth_sub, role, status, timezone
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateUserStatus :exec
UPDATE users
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserTheme :exec
UPDATE users
SET theme_preference = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserTimezone :exec
UPDATE users
SET timezone = $2, updated_at = NOW()
WHERE id = $1;

-- name: GetMoodEntry :one
SELECT * FROM mood_entries
WHERE user_id = $1 AND date = $2
LIMIT 1;

-- name: ListMoodEntriesForYear :many
SELECT * FROM mood_entries
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC;

-- name: CreateMoodEntry :one
INSERT INTO mood_entries (
    user_id, date, rating, notes
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateMoodEntry :exec
UPDATE mood_entries
SET rating = $2, notes = $3, updated_at = NOW()
WHERE user_id = $1 AND date = sqlc.arg(entry_date);

-- name: DeleteMoodEntry :exec
DELETE FROM mood_entries
WHERE user_id = $1 AND date = $2;

-- name: ListActiveHabits :many
SELECT * FROM habits
WHERE user_id = $1 AND archived_at IS NULL
ORDER BY created_at ASC;

-- name: ListAllHabits :many
SELECT * FROM habits
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetHabit :one
SELECT * FROM habits
WHERE id = $1 LIMIT 1;

-- name: CreateHabit :one
INSERT INTO habits (
    user_id, name, description, type, goal_value, goal_unit, color, emoji
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateHabit :exec
UPDATE habits
SET name = $2, description = $3, goal_value = $4, goal_unit = $5, color = $6, emoji = $7, updated_at = NOW()
WHERE id = $1;

-- name: ArchiveHabit :exec
UPDATE habits
SET archived_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetHabitEntry :one
SELECT * FROM habit_entries
WHERE habit_id = $1 AND date = $2
LIMIT 1;

-- name: ListHabitEntriesForDateRange :many
SELECT * FROM habit_entries
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC;

-- name: ListHabitEntriesForHabit :many
SELECT * FROM habit_entries
WHERE habit_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC;

-- name: CreateHabitEntry :one
INSERT INTO habit_entries (
    habit_id, user_id, date, completed, value, notes
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateHabitEntry :exec
UPDATE habit_entries
SET completed = $2, value = $3, notes = $4, updated_at = NOW()
WHERE habit_id = $1 AND date = sqlc.arg(entry_date);

-- name: DeleteHabitEntry :exec
DELETE FROM habit_entries
WHERE habit_id = $1 AND date = $2;
