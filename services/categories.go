package services

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
	"github.com/melardev/api_blog_app/infrastructure"
)

func FetchAllCategories() ([]models.Category, error) {
	database := infrastructure.GetDB()
	var categories []models.Category
	err := database.Preload("Images", "category_id IS NOT NULL").Find(&categories).Error
	return categories, err
}
