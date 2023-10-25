package database

import "github.com/pocketbase/pocketbase/models"

var _ models.Model = (*User)(nil)

type User struct {
	models.BaseModel

	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (u *User) TableName() string {
	return "users"
}

func FromAuth(auth *models.Record) *User {
	if auth == nil {
		return nil
	}
	return &User{
		Id:   auth.GetString("id"),
		Name: auth.GetString("name"),
	}
}
