package models

import "github.com/pocketbase/pocketbase/models"

var _ models.Model = (*Board)(nil)

type Board struct {
	models.BaseModel
	Name string `db:"name" json:"name"`
}

func (t *Board) TableName() string {
	return "boards"
}
