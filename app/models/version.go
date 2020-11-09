package models

import (
	"github.com/revel/revel"
)

type Version struct {
	Id          	int		`db:"id" json:"id"`
	RawData    		string	`db:"rawData" json:"rawData"`
	CreatedBy          string	`db:"createdBy" json:"createdBy"`
	CreationDate    string	`db:"creationDate" json:"creationDate"`
	File			File	`db:"file" json:"file"`
}

func (version *Version) Validate(v *revel.Validation) {
	v.Check(version.RawData, revel.Required{})
	v.Check(version.CreationDate, revel.Required{})
	v.Check(version.File, revel.Required{})
	v.Check(version.CreatedBy, revel.Required{})
}

