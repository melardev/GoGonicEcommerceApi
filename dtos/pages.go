package dtos

import "github.com/melardev/GoGonicEcommerceApi/models"

func CreateHomeResponse(tags []models.Tag, categories []models.Category) map[string]interface{} {
	return CreateSuccessDto(map[string]interface{}{
		"tags":       CreateTagListDto(tags),
		"categories": CreateCategoryListDto(categories),
	})
}
