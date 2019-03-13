package dtos

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
	"net/http"
	"time"
)

type CreateComment struct {
	Content string `form:"content" json:"content" xml:"content"  binding:"required"`
}

func CreateCommentPagedResponse(request *http.Request, comments []models.Comment, page, page_size, count int, bools ...bool) map[string]interface{} {
	var resources = make([]interface{}, len(comments))
	for index, comment := range comments {
		includeUser := false
		if len(bools) > 0 {
			includeUser = bools[0]
		}
		includeProduct := false
		if len(bools) > 1 {
			includeProduct = bools[1]
		}

		resources[index] = GetSummary(&comment, includeUser, includeProduct)
	}
	return CreatePagedResponse(request, resources, "comments", page, page_size, count)
}

func GetCommentDetailsDto(comment *models.Comment, includes ...bool) map[string]interface{} {
	includeUser := false
	if len(includes) > 0 {
		includeUser = includes[0]
	}
	includeProduct := false
	if len(includes) > 1 {
		includeProduct = includes[1]
	}
	return GetSummary(comment, includeUser, includeProduct)
}

func GetSummary(comment *models.Comment, includeUser, includeProduct bool) map[string]interface{} {
	result := map[string]interface{}{
		"id":         comment.ID,
		"content":    comment.Content,
		"created_at": comment.CreatedAt.UTC().Format(time.RFC1123),
		"updated_at": comment.UpdatedAt.UTC().Format(time.RFC1123),
	}
	if includeUser == true {
		result["user"] = map[string]interface{}{
			"id":       comment.User.ID,
			"username": comment.User.Username,
		}
	}
	if includeProduct == true {
		result["product"] = map[string]interface{}{
			"id":   comment.Product.ID,
			"name": comment.Product.Name,
			"slug": comment.Product.Slug,
		}
	}
	return result
}

func CreateCommentCreatedDto(comment *models.Comment, includes ...bool) map[string]interface{} {
	return CreateSuccessWithDtoAndMessageDto(GetCommentDetailsDto(comment, includes...), "Comment created successfully")
}
