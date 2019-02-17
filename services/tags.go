package services

import "github.com/melardev/api_shop_gonic/models"

func FetchAllTags() ([]models.Tag, error) {
	database := models.GetDB()
	var tags []models.Tag
	err := database.Preload("Images", "tag_id IS NOT NULL").Find(&tags).Error
	return tags, err
}
