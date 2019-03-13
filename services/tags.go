package services

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
	"github.com/melardev/api_blog_app/infrastructure"
)

func FetchAllTags() ([]models.Tag, error) {
	database := infrastructure.GetDB()
	var tags []models.Tag
	err := database.Preload("Images", "tag_id IS NOT NULL").Find(&tags).Error
	return tags, err
}
