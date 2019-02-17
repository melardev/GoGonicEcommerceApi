package services

import "github.com/melardev/api_shop_gonic/models"

func FetchAllCategories() ([]models.Category, error) {
	database := models.GetDB()
	var categories []models.Category
	err := database.Preload("Images", "category_id IS NOT NULL").Find(&categories).Error
	return categories, err
}
