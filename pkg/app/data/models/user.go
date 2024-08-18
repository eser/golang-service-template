package models

import "time"

type UserKind string

const (
	UserKindAdmin   UserKind = "admin"
	UserKindEditor  UserKind = "editor"
	UserKindRegular UserKind = "regular"
)

type User struct {
	Id   string
	Kind UserKind

	Name  string
	Email string
	Phone *string

	IndividualProfileId *string
	GithubRemoteId      *string
	GithubHandle        *string
	XRemoteId           *string
	XHandle             *string

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
