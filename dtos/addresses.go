package dtos

import (
	"github.com/melardev/api_shop_gonic/models"
	"net/http"
)

type CreateAddress struct {
	FirstName     string `form:"first_name" json:"first_name" xml:"first_name"`
	LastName      string `form:"last_name" json:"last_name" xml:"last_name"`
	Country       string `form:"country" json:"country" xml:"country" binding:"required"`
	City          string `form:"city" json:"city" xml:"city" binding:"required"`
	StreetAddress string `form:"address" json:"address" xml:"address" binding:"required"`
	ZipCode       string `form:"zip_code" json:"zip_code" xml:"zip_code" binding:"required"`
}

func CreateAddressPagedResponse(request *http.Request, addresses []models.Address, page, page_size, count int, includeUser bool) map[string]interface{} {
	var resources = make([]interface{}, len(addresses))
	for index, address := range addresses {
		resources[index] = GetAddressDto(&address, includeUser)
	}
	return CreatePagedResponse(request, resources, "addresses", page, page_size, count)
}

func GetAddressDto(address *models.Address, includeUser bool) map[string]interface{} {
	dto := map[string]interface{}{
		"id":         address.ID,
		"first_name": address.FirstName,
		"last_name":  address.LastName,
		"zip_code":   address.ZipCode,
		"country":    address.Country,
		"city":       address.City,
	}

	if includeUser {
		dto["user"] = map[string]interface{}{
			"id":       address.UserId,
			"username": address.User.Username,
		}
	}
	return dto
}

func GetAddressCreatedDto(address *models.Address, includeUser bool) map[string]interface{} {
	return CreateSuccessWithDtoAndMessageDto(GetAddressDto(address, includeUser), "StreetAddress created successfully")
}
