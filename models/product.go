package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Name        string       `gorm:"size:280;not null"`
	Description string       `gorm:"not null"`
	Slug        string       `gorm:"unique_index;not null"`
	Price       int          `gorm:"not null"`
	Stock       int          `gorm:"not null"`
	Tags        []Tag        `gorm:"many2many:products_tags;"`
	ProductTags []ProductTag `gorm:"foreignkey:ProductId"`

	Categories        []Category        `gorm:"many2many:products_categories;"`
	ProductCategories []ProductCategory `gorm:"foreignkey:ProductId"`

	Comments      []Comment    `gorm:"foreignKey:ProductId"`
	Images        []FileUpload `gorm:"foreignKey:ProductId"`
	CommentsCount int          `gorm:"-"`
}

func (product *Product) BeforeSave() (err error) {
	product.Slug = slug.Make(product.Name)
	return
}
