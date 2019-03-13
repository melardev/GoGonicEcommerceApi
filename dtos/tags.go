package dtos

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
	"strings"
)

type CreateTag struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description" binding:"required"`
}

func CreateTagListMapDto(tags []models.Tag) map[string]interface{} {
	result := map[string]interface{}{}
	var t = make([]interface{}, len(tags))
	for i := 0; i < len(tags); i++ {
		t[i] = CreateTagDto(tags[i])
	}
	result["tags"] = t
	return CreateSuccessDto(result)
}

func CreateTagListDto(tags []models.Tag) []interface{} {
	var t = make([]interface{}, len(tags))
	for i := 0; i < len(tags); i++ {
		t[i] = CreateTagDto(tags[i])
	}
	return t
}

func CreateTagDto(tag models.Tag) map[string]interface{} {
	var imageUrls = make([]string, len(tag.Images))
	replaceAllFlag := -1
	for i := 0; i < len(tag.Images); i++ {
		imageUrls[i] = strings.Replace(tag.Images[i].FilePath, "\\", "/", replaceAllFlag)
	}
	return map[string]interface{}{
		"id":          tag.ID,
		"name":        tag.Name,
		"description": tag.Description,
		"image_urls":  imageUrls,
	}
}

func CreateTagCreatedDto(tag models.Tag) map[string]interface{} {
	return CreateSuccessWithDtoAndMessageDto(CreateTagDto(tag), "Tag created successfully")
}
