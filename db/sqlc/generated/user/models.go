// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package user

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID       pgtype.UUID
	Username pgtype.Text
	Email    string
	Password pgtype.Text
}