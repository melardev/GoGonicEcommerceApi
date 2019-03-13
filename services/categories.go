package services

import (
	"github.com/melardev/GoGonicEcommerceApi/infrastructure"
	"github.com/melardev/GoGonicEcommerceApi/models"
)

func FetchAllCategories() ([]models.Category, error) {
	database := infrastructure.GetDb()
	var categories []models.Category
	err := database.Preload("Images", "category_id IS NOT NULL").Find(&categories).Error
	return categories, err
}
