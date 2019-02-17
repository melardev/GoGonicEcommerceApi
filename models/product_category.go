package models

type ProductCategory struct {
	Category   User `gorm:"association_foreignkey:CategoryId"`
	CategoryId uint
	Product    User `gorm:"association_foreignkey:ProductId"`
	ProductId  uint
}

func (*ProductCategory) TableName() string {
	return "products_categories"
}
