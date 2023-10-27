package models

import "github.com/pocketbase/pocketbase/models"

type (
	Admin = models.Admin
)

var _ models.Model = (*User)(nil)

type User struct {
	models.BaseModel

	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (u *User) TableName() string {
	return "users"
}
