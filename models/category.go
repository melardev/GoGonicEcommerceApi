package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Category struct {
	gorm.Model
	Name        string       `gorm:"not null"`
	Description string       `gorm:"default:null"`
	Slug        string       `gorm:"unique_index"`
	Products    []Product    `gorm:"many2many:products_categories;"`
	Images      []FileUpload `gorm:"foreignKey:CategoryId"`
	IsNewRecord bool         `gorm:"-;default:false"`
}

func (a *Category) BeforeSave() (err error) {
	a.Slug = slug.Make(a.Name)
	return
}
