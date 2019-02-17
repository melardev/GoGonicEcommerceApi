package models

import "github.com/jinzhu/gorm"

type FileUpload struct {
	gorm.Model
	Filename     string
	FilePath     string
	OriginalName string
	FileSize     uint

	Tag   Tag  `gorm:"association_foreignkey:TagId"`
	TagId uint `gorm:"default:null"`

	Category   Category `gorm:"association_foreignkey:CategoryId"`
	CategoryId uint     `gorm:"default:null"`

	Product   Category `gorm:"association_foreignkey:ProductId"`
	ProductId uint     `gorm:"default:null"`
}

// Scopes, not used
func TagImages(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", "TagImage")
}

func CategoryImages(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", "CategoryImage")
}

func ProductImages(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", "ProductImage")
}

// db.Scopes(CategoryImages, ProductImages).Find(&images)
