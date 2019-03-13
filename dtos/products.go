package dtos

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
	"net/http"
	"strings"
	"time"
)

type ManagedModel models.Product

type CreateProduct struct {
	Name        string `form:"name" json:"name" xml:"name" binding:"required"`
	Description string `form:"description" json:"description" xml:"description" binding:"required"`
	Price       int    `form:"price" json:"price" xml:"price" binding:"required"`
	Stock       int    `form:"stock" json:"stock" xml:"stock" binding:"required"`
}

func CreatedProductPagedResponse(request *http.Request, products []models.Product, page, page_size, count int, commentsCount []int) interface{} {
	var resources = make([]interface{}, len(products))
	for index, product := range products {
		resources[index] = CreateProductDto(&product, commentsCount[index])
	}
	return CreatePagedResponse(request, resources, "products", page, page_size, count)
}

func CreateProductDto(product *models.Product, commentCount int) map[string]interface{} {

	var tags = make([]map[string]interface{}, len(product.Tags))
	var categories = make([]map[string]interface{}, len(product.Categories))
	var images = make([]string, len(product.Images))

	for index, tag := range product.Tags {
		tags[index] = map[string]interface{}{
			"id":   tag.ID,
			"name": tag.Name,
			"slug": tag.Slug,
		}
	}

	for index, category := range product.Categories {
		categories[index] = map[string]interface{}{
			"id":   category.ID,
			"name": category.Name,
			"slug": category.Slug,
		}
	}
	replaceAllFlag := -1
	for index, image := range product.Images {
		images[index] = strings.Replace(image.FilePath, "\\", "/", replaceAllFlag)
	}

	for index, tag := range product.Tags {
		tags[index] = map[string]interface{}{
			"id":   tag.ID,
			"name": tag.Name,
			"slug": tag.Slug,
		}
	}

	result := map[string]interface{}{
		"id":         product.ID,
		"name":       product.Name,
		"slug":       product.Slug,
		"price":      product.Price,
		"stock":      product.Stock,
		"tags":       tags,
		"categories": categories,
		"image_urls": images,
		"created_at": product.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		"updated_at": product.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}

	if commentCount >= 0 {
		// "comments_count": product.CommentsCount,
		result["comments_count"] = commentCount
	}
	return result
}

func CreateProductDetailsDto(product models.Product) map[string]interface{} {
	result := CreateProductDto(&product, -1)
	result["description"] = product.Description
	comments := make([]map[string]interface{}, len(product.Comments))
	for index, comment := range product.Comments {
		comments[index] = GetSummary(&comment, true, false)
	}

	result["comments"] = comments
	return result
}
func CreateProductCreatedDto(product models.Product) map[string]interface{} {
	return CreateSuccessWithDtoAndMessageDto(CreateProductDetailsDto(product), "Product crated successfully")
}
