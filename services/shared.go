package services

import (
	"github.com/melardev/api_blog_app/infrastructure"
)

func CreateOne(data interface{}) error {
	database := infrastructure.GetDB()
	err := database.Create(data).Error
	return err
}

func SaveOne(data interface{}) error {
	database := infrastructure.GetDB()
	err := database.Save(data).Error
	return err
}
