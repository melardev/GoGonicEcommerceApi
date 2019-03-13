package controllers

// import "C"
import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/melardev/GoGonicEcommerceApi/infrastructure"
	"github.com/melardev/GoGonicEcommerceApi/models"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/melardev/GoGonicEcommerceApi/dtos"

	"github.com/melardev/GoGonicEcommerceApi/middlewares"
	"github.com/melardev/GoGonicEcommerceApi/services"

	"net/http"
	"strconv"
)

func RegisterProductRoutes(router *gin.RouterGroup) {
	router.GET("/", ProductList)
	router.GET("/:slug", GetProductDetailsBySlug)

	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("/", CreateProduct)
		router.DELETE("/:slug", ProductDelete)
	}
}

func ProductList(c *gin.Context) {

	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	productModels, modelCount, commentsCount, err := services.FetchProductsPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("products", errors.New("Invalid param")))
		return
	}

	c.JSON(http.StatusOK, dtos.CreatedProductPagedResponse(c.Request, productModels, page, pageSize, modelCount, commentsCount))
}

func GetProductDetailsBySlug(c *gin.Context) {
	productSlug := c.Param("slug")

	product := services.FetchProductDetails(&models.Product{Slug: productSlug}, true)
	if product.ID == 0 {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("products", errors.New("Invalid slug")))
		return
	}
	c.JSON(http.StatusOK, dtos.CreateProductDetailsDto(product))
}

func CreateProduct(c *gin.Context) {
	// Only admin users can create products
	user := c.Keys["currentUser"].(models.User)
	if user.IsNotAdmin() {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied, you must be admin"))
		return
	}

	var formDto dtos.CreateProduct
	if err := c.ShouldBind(&formDto); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	name := formDto.Name
	description := formDto.Description

	price := formDto.Price
	stock, err := strconv.ParseInt(c.PostForm("stock"), 10, 32)
	form, err := c.MultipartForm()

	tagCount := 0
	catCount := 0
	for key := range form.Value {
		if strings.HasPrefix(key, "tags[") {
			tagCount++
		}
		if strings.HasPrefix(key, "category[") {
			catCount++
		}
	}

	var tags = make([]models.Tag, tagCount)
	var categories = make([]models.Category, catCount)

	var rgx = regexp.MustCompile(`\[(.*?)\]`)
	database := infrastructure.GetDb()
	tagPtr := 0
	catPtr := 0

	for k, v := range form.Value {
		if strings.HasPrefix(k, "tags[") {
			result := rgx.FindStringSubmatch(k)
			var tag models.Tag
			name := result[1]
			description := v[0]
			database.Where(&models.Tag{Slug: slug.Make(name)}).
				Attrs(models.Tag{Name: name, Description: description}).
				FirstOrCreate(&tag)
			tags[tagPtr] = tag
			tagPtr++
		}

		if strings.HasPrefix(k, "category[") {
			result := rgx.FindStringSubmatch(k)
			var category models.Category
			name := result[1]
			description := v[0]
			database.Where(&models.Category{Slug: slug.Make(name)}).
				Attrs(models.Category{Name: name, Description: description}).
				FirstOrCreate(&category)
			categories[catPtr] = category
			catPtr++
		}
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateDetailedErrorDto("form_error", err))
		return
	}

	files := form.File["images[]"]
	var productImages = make([]models.FileUpload, len(files))

	for index, file := range files {
		fileName := randomString(16) + ".png"

		dirPath := filepath.Join(".", "static", "images", "products")
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
		productImages[index] = models.FileUpload{Filename: fileName, OriginalName: file.Filename, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
	}

	product := models.Product{
		Name:        name,
		Description: description,
		Tags:        tags,
		Categories:  categories,
		Price:       (int)(price),
		Stock:       (int)(stock),
		Images:      productImages,
	}

	if err := services.CreateOne(&product); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database", err))
		return
	}

	c.JSON(http.StatusOK, dtos.CreateProductCreatedDto(product))

}

func ProductDelete(c *gin.Context) {
	slug := c.Param("slug")
	err := services.DeleteProduct(&models.Product{Slug: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("products", errors.New("Invalid slug")))
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": "Delete success"})
}
