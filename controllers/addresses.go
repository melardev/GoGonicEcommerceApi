package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/melardev/api_shop_gonic/dtos"
	"github.com/melardev/api_shop_gonic/middlewares"
	"github.com/melardev/api_shop_gonic/models"
	"github.com/melardev/api_shop_gonic/services"

	"net/http"
	"strconv"
)

func RegisterAddressesRoutes(router *gin.RouterGroup) {

	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.GET("/addresses", ListAddresses)
		router.POST("/addresses", CreateAddress)
	}

}

func ListAddresses(c *gin.Context) {

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

	// userId:= c.Keys["currentUserId"].(uint) // or
	userId := c.MustGet("currentUserId").(uint)
	includeUser := false
	addresses, totalCommentCount := services.FetchAddressesPage(userId, page, pageSize, includeUser)

	c.JSON(http.StatusOK, dtos.CreateAddressPagedResponse(c.Request, addresses, page, pageSize, totalCommentCount, includeUser))
}

func CreateAddress(c *gin.Context) {

	user := c.MustGet("currentUser").(models.User)

	var json dtos.CreateAddress
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}
	firstName := json.FirstName
	lastName := json.LastName
	if firstName == "" {
		firstName = user.FirstName
	}
	if lastName == "" {
		lastName = user.LastName
	}
	address := models.Address{
		FirstName:     firstName,
		LastName:      lastName,
		Country:       json.Country,
		City:          json.City,
		StreetAddress: json.StreetAddress,
		ZipCode:       json.ZipCode,
		User:          user,
		UserId:        user.ID,
	}

	if err := services.SaveOne(&address); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database_error", err))
		return
	}

	c.JSON(http.StatusOK, dtos.GetAddressCreatedDto(&address, false))
}
