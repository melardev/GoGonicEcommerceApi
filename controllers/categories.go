package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/melardev/api_shop_gonic/dtos"
	"github.com/melardev/api_shop_gonic/middlewares"
	"github.com/melardev/api_shop_gonic/models"
	"github.com/melardev/api_shop_gonic/services"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func RegisterCategoryRoutes(router *gin.RouterGroup) {
	router.GET("", CategoryList)
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("", CreateCategory)
	}
}

func CategoryList(c *gin.Context) {
	tags, err := services.FetchAllCategories()
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("fetch_error", err))
		return
	}
	c.JSON(http.StatusOK, dtos.CreateCategoryListDto(tags))
}

func CreateCategory(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)
	if user.IsNotAdmin() {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied, you must be admin"))
		return
	}
	name := c.PostForm("name")
	description := c.PostForm("description")

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["images[]"]
	var categoryImages = make([]models.FileUpload, len(files))
	for index, file := range files {
		fileName := randomString(16) + ".png"

		dirPath := filepath.Join(".", "static", "images", "categories")
		filePath := filepath.Join(dirPath, fileName)
		// Create directory if does not exist
		if _, err = os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, os.ModeDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, dtos.CreateDetailedErrorDto("io_error", err))
				return
			}
		}
		// Create file that will hold the image
		outputFile, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		// Open the temporary file that contains the uploaded image
		inputFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusOK, dtos.CreateDetailedErrorDto("io_error", err))
		}
		defer inputFile.Close()

		// Copy the temporary image to the permanent location outputFile
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			log.Fatal(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		fileSize := (uint)(file.Size)
		categoryImages[index] = models.FileUpload{Filename: file.Filename, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
	}

	database := models.GetDB()
	category := models.Category{Name: name, Description: description, Images: categoryImages}

	// TODO: Why it is performing a SELECT SQL Query per image?
	// Even worse, it is selecting category_id, why??
	// SELECT "tag_id", "product_id" FROM "file_uploads"  WHERE (id = insertedFileUploadId)
	err = database.Create(&category).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.CreateDetailedErrorDto("db_error", err))
	}
	c.JSON(http.StatusOK, dtos.CreateCategoryCreatedDto(category))
}
