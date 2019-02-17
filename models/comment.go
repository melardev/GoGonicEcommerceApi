package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	Content   string  `gorm:"size:2048"`
	Rating    int     `gorm:"default:null"`
	Product   Product `gorm:"association_foreignkey:ProductId"`
	ProductId uint    `gorm:"not null"`
	User      User    `gorm:"association_foreignkey:UserId"`
	UserId    uint    `gorm:"not null"`
}
