package models

import "github.com/pocketbase/pocketbase/models"

var _ models.Model = (*Task)(nil)

type Task struct {
	models.BaseModel
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Board       string `db:"board" json:"board"`
	Owner       string `db:"owner" json:"owner"`
	Completed   bool   `db:"completed" json:"completed"`
}

func (t *Task) TableName() string {
	return "tasks"
}
