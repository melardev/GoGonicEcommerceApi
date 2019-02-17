package services

import (
	"github.com/melardev/api_blog_app/db"
	"github.com/melardev/api_shop_gonic/models"
)

// You could input the conditions and it will return an User in database with error info.
// 	userModel, err := FindOneUser(&User{Username: "username0"})
func FindOneUser(condition interface{}) (models.User, error) {
	database := models.GetDB()
	var user models.User

	err := database.Where(condition).Preload("Roles").First(&user).Error
	return user, err
}

// You could update properties of an User to database returning with error info.
//  err := db.Model(userModel).Update(User{Username: "wangzitian0"}).Error
func UpdateUser(user models.User, data interface{}) error {
	database := db.GetDB()
	err := database.Model(user).Update(data).Error
	return err
}
