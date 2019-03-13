package models

type ProductTag struct {
	Tag       User `gorm:"association_foreignkey:TagId"`
	TagId     uint
	Product   Product `gorm:"association_foreignkey:ProductId"`
	ProductId uint
}

func (*ProductTag) TableName() string {
	return "products_tags"
}
