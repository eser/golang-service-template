// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package data

import (
	"database/sql"
)

type Channel struct {
	Id   string         `json:"id"`
	Name sql.NullString `json:"name"`
}