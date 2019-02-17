package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/melardev/api_shop_gonic/dtos"
	"github.com/melardev/api_shop_gonic/middlewares"
	"github.com/melardev/api_shop_gonic/models"
	"github.com/melardev/api_shop_gonic/services"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

func RegisterTagRoutes(router *gin.RouterGroup) {
	router.GET("", TagList)
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("", CreateTag)
	}
}

func TagList(c *gin.Context) {
	tags, err := services.FetchAllTags()
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("fetch_error", err))
		return
	}
	c.JSON(http.StatusOK, dtos.CreateTagListDto(tags))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CreateTag(c *gin.Context) {
	user := c.Keys["currentUser"].(models.User)
	if user.IsNotAdmin() {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied, you must be admin"))
		return
	}
	var createForm dtos.CreateTag
	// name := c.PostForm("name")
	// description := c.PostForm("description")

	// If you wanna know more about how binding is done internally check gin-gonic/bin/binding.formBinding.Bind at form.go
	if err := c.ShouldBind(&createForm); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["images[]"]
	var tagImages = make([]models.FileUpload, len(files))
	for index, file := range files {
		fileName := randomString(16) + ".png"

		dirPath := filepath.Join(".", "static", "images", "tags")
		filePath := filepath.Join(dirPath, fileName)
		if _, err = os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, os.ModeDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, dtos.CreateDetailedErrorDto("io_error", err))
				return
			}
		}
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusBadRequest, dtos.CreateDetailedErrorDto("upload_error", err))
			return
		}
		fileSize := (uint)(file.Size)
		tagImages[index] = models.FileUpload{Filename: fileName, OriginalName: file.Filename, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
	}

	database := models.GetDB()
	tag := models.Tag{Name: createForm.Name, Description: createForm.Description, Images: tagImages}
	// TODO: Why it is performing a SELECT SQL Query per image?
	// Even worse, it is selecting category_id, why??
	// SELECT "category_id", "product_id" FROM "file_uploads"  WHERE (id = insertedFileUploadId)
	err = database.Create(&tag).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.CreateDetailedErrorDto("db_error", err))
	}

	c.JSON(http.StatusOK, dtos.CreateTagCreatedDto(tag))
}
