package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Tag struct {
	gorm.Model
	Name        string       `gorm:"not null"`
	Description string       `gorm:"default:null"`
	Slug        string       `gorm:"unique_index"`
	Products    []Product    `gorm:"many2many:products_tags;"`
	Images      []FileUpload `gorm:"foreignKey:TagId"`
	IsNewRecord bool         `gorm:"-;default:false"` // Virtual Field, so it is not persisted in the Db. This is used in FirstOrCreate()
}

func (a *Tag) BeforeSave() (err error) {
	a.Slug = slug.Make(a.Name)
	return
}
