package services

import (
	"github.com/melardev/api_shop_gonic/models"
)

func CreateOne(data interface{}) error {
	database := models.GetDB()
	err := database.Create(data).Error
	return err
}

func SaveOne(data interface{}) error {
	database := models.GetDB()
	err := database.Save(data).Error
	return err
}
