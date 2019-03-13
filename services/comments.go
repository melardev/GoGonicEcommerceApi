package services

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
	"github.com/melardev/api_blog_app/infrastructure"
)

func FetchCommentsPage(productId, page int, page_size int) ([]models.Comment, int) {
	// TODO: Why Preload does not load the User? the error is can't preload field User for models.Comment

	var comments []models.Comment
	var totalCommentCount int
	database := infrastructure.GetDB()
	database.Model(&comments).Where(&models.Comment{ProductId: uint(productId)}).Count(&totalCommentCount)
	database.Where(&models.Comment{ProductId: uint(productId)}).
		Offset((page - 1) * page_size).Limit(page_size).
		Preload("User").
		Find(&comments)

	// `Where in` using other columns different than ID
	// database.Where("username in (?)", []string{"admin", "melardev"}).Find(&users)
	var userIds = make([]uint, len(comments))
	var users []models.User
	for i := 0; i < len(comments); i++ {
		userIds[i] = comments[i].UserId
	}
	database.Select("id, username").Where(userIds).Find(&users)

	// If the user gets deleted and the comment is still in the database we may have less users than comments
	// Another scenario (the one I run into) is there is a problem with the Comment.User, the Comment.UserId does not get saved automatically
	for i := 0; i < len(comments); i++ {
		comment := comments[i]
		for j := 0; j < len(users); j++ {
			user := users[j]
			if comment.UserId == user.ID {
				comments[i].User = users[j]
			}
		}

	}
	return comments, totalCommentCount
}

func FetchCommentById(id int, includes ...bool) models.Comment {
	includeUser := false
	if len(includes) > 0 {
		includeUser = includes[0]
	}
	includeProduct := false
	if len(includes) > 1 {
		includeProduct = includes[1]
	}
	database := infrastructure.GetDB()
	var comment models.Comment
	if includeProduct && includeUser {
		database.Preload("User").Preload("Product").Find(&comment, id)
	} else if includeUser {
		database.Preload("User").Find(&comment, id)
	} else if includeProduct {
		database.Preload("Product").Find(&comment, id)
	} else {
		database.Find(&comment, id)
	}
	return comment
}

func DeleteComment(condition interface{}) error {
	database := infrastructure.GetDB()
	err := database.Where(condition).Delete(models.Comment{}).Error
	return err
}
