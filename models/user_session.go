package models

import "time"

const (
	SESSION_ACTIVE int = iota
	SESSION_DEACTIVE
)

type UserSession struct {
	ID int64 `sql:",pk"`

	Referer   string
	UserAgent string
	ClientIp  string

	UserID int64
	User   *User

	Status int

	CreatedAt time.Time
	UpdatedAt time.Time
}
