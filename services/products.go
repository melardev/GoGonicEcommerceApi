package services

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
	"github.com/melardev/api_blog_app/infrastructure"
)

func FetchProductsPage(page int, page_size int) ([]models.Product, int, []int, error) {
	database := infrastructure.GetDB()
	var products []models.Product
	var count int
	tx := database.Begin()
	database.Model(&products).Count(&count)
	database.Offset((page - 1) * page_size).Limit(page_size).Find(&products)
	tx.Model(&products).
		Preload("Tags").Preload("Categories").Preload("Images").
		Order("created_at desc").Offset((page - 1) * page_size).Limit(page_size).Find(&products)
	commentsCount := make([]int, len(products))

	for index, product := range products {
		commentsCount[index] = tx.Model(&product).Association("Comments").Count()
	}
	err := tx.Commit().Error
	return products, count, commentsCount, err
}

func FetchProductDetails(condition interface{}, optional ...bool) models.Product {
	database := infrastructure.GetDB()
	var product models.Product

	query := database.Where(condition).
		Preload("Tags").Preload("Categories").Preload("Images").Preload("Comments")
	// Unfortunately .Preload("Comments.User") does not work as the doc states ...
	query.First(&product)
	includeUserComment := false

	if len(optional) > 0 {
		includeUserComment = optional[0]
	}

	if includeUserComment {

		for i := 0; i < len(product.Comments); i++ {
			database.Model(&product.Comments[i]).Related(&product.Comments[i].User, "UserId")
		}

		var userIds = make([]uint, len(product.Comments))
		var users []models.User
		for i := 0; i < len(product.Comments); i++ {
			userIds[i] = product.Comments[i].UserId
		}
		// WHERE users.id IN userIds; This will also work: Select([]string{"id", "username"})
		database.Select("id, username").Where(userIds).Find(&users)

		for i := 0; i < len(product.Comments); i++ {
			user := users[i]
			comment := product.Comments[i]
			if comment.UserId == user.ID {
				product.Comments[i].User = users[i]
			}
		}
	}

	return product
}

func FetchProductId(slug string) (uint, error) {
	productId := -1
	database := infrastructure.GetDB()
	err := database.Model(&models.Product{}).Where(&models.Product{Slug: slug}).Select("id").Row().Scan(&productId)
	return uint(productId), err
}

func SetTags(product *models.Product, tags []string) error {
	database := infrastructure.GetDB()
	var tagList []models.Tag
	for _, tag := range tags {
		var tagModel models.Tag
		err := database.FirstOrCreate(&tagModel, models.Tag{Name: tag}).Error
		if err != nil {
			return err
		}
		tagList = append(tagList, tagModel)
	}
	product.Tags = tagList
	return nil
}

func Update(product *models.Product, data interface{}) error {
	database := infrastructure.GetDB()
	err := database.Model(product).Update(data).Error
	return err
}

func DeleteProduct(condition interface{}) error {
	db := infrastructure.GetDB()
	err := db.Where(condition).Delete(models.Product{}).Error
	return err
}

func FetchProductsIdNameAndPrice(productIds []uint) (products []models.Product, err error) {
	database := infrastructure.GetDB()
	err = database.Select([]string{"id", "name", "slug", "price"}).Find(&products, productIds).Error
	return products, err
}
