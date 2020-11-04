package models

import (
	"github.com/revel/revel"
)

type File struct {
	Id           int    `db:"id" json:"id"`
	Title        string `db:"title" json:"title"`
	Description  string `db:"description" json:"description"`
	CreatedBy       string `db:"createdBy" json:"createdBy"`
	CreationDate string `db:"creationDate" json:"creationDate"`
}

func (file *File) Validate(v *revel.Validation) {
	v.Check(file.Title, revel.Required{})
	v.Check(file.CreationDate, revel.Required{})
	v.Check(file.CreatedBy, revel.Required{})
}