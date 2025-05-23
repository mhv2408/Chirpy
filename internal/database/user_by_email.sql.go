// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user_by_email.sql

package database

import (
	"context"
)

const userByEmail = `-- name: UserByEmail :one

SELECT id, created_at, updated_at, email, hashed_password 
FROM users
WHERE email=$1
`

func (q *Queries) UserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, userByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}
