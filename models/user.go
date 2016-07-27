package models

import "time"

const (
	USER_CREATE int = iota
	USER_ACTIVE
	USER_DEACTIVE
)

type User struct {
	ID int64 `sql:",pk"`

	Name     string
	Email    string
	Mobile   string
	Password string
	Role     string

	Status int

	CreatedAt time.Time `sql:",null"`
	UpdatedAt time.Time `sql:",null"`
}
