package services

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
	"github.com/melardev/api_blog_app/infrastructure"
)

func FetchAddressesPage(userId uint, page, pageSize int, includeUser bool) ([]models.Address, int) {
	var addresses []models.Address
	var totalAddressesCount int
	database := infrastructure.GetDB()
	database.Model(&models.Address{}).Where(&models.Address{UserId: uint(userId)}).Count(&totalAddressesCount)
	database.Where(&models.Address{UserId: uint(userId)}).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Preload("User").
		Find(&addresses)

	if includeUser {
		var userIds = make([]uint, len(addresses))
		var users []models.User
		for i := 0; i < len(addresses); i++ {
			userIds[i] = addresses[i].UserId
		}
		database.Select([]string{"id", "username"}).Where(userIds).Find(&users)

		// If the user gets deleted and the comment is still in the database we may have less users than addresses
		// Another scenario (the one I run into) is there is a problem with the Comment.User, the Comment.UserId does not get saved automatically
		for i := 0; i < len(addresses); i++ {
			address := addresses[i]
			for j := 0; j < len(users); j++ {
				user := users[j]
				if address.UserId == user.ID {
					addresses[i].User = users[j]
				}
			}
		}
	}
	return addresses, totalAddressesCount
}

func FetchAddress(addressId uint) (address models.Address) {
	database := infrastructure.GetDB()
	database.First(&address, addressId)
	return address
}

func FetchIdsFromAddress(addressId uint) (address models.Address) {
	database := infrastructure.GetDB()
	database.Select("id, user_id").First(&address, addressId)
	return
}
