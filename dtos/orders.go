package dtos

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
	"net/http"
)

type CreateOrderRequestDto struct {
	FirstName     string `form:"first_name" json:"first_name" xml:"first_name"`
	LastName      string `form:"last_name" json:"last_name" xml:"last_name"`
	Country       string `form:"country" json:"country" xml:"country"`
	City          string `form:"city" json:"city" xml:"city"`
	StreetAddress string `form:"street_address" json:"street_address" xml:"street_address" `
	ZipCode       string `form:"zip_code" json:"zip_code" xml:"zip_code" `
	AddressId     uint   `form:"address_id" json:"address_id" xml:"address_id" `
	CartItems     []struct {
		Id       uint `form:"id" json:"id" binding:"required"`
		Quantity int  `form:"quantity" json:"quantity" binding:"required"`
	} `json:"cart_items"`
}

func CreateOrderPagedResponse(request *http.Request, orders []models.Order, page, page_size, totalOrdersCount int, includes ...bool) map[string]interface{} {
	var resources = make([]interface{}, len(orders))
	for index, order := range orders {

		includeAddress, includeOrderItems, includeUser := getIncludeFlags(includes...)

		resources[index] = CreateOrderDto(&order, includeAddress, includeOrderItems, includeUser)
	}
	return CreatePagedResponse(request, resources, "orders", page, page_size, totalOrdersCount)
}

func CreateOrderDto(order *models.Order, includes ...bool) map[string]interface{} {

	includeAddress, includeOrderItems, includeUser := getIncludeFlags(includes...)

	result := map[string]interface{}{
		"id":              order.ID,
		"tracking_number": order.TrackingNumber,
		"order_status":    order.GetOrderStatusAsString(),
	}

	if includeAddress {
		result["address"] = map[string]interface{}{
			"first_name":     order.Address.FirstName,
			"last_name":      order.Address.LastName,
			"street_address": order.Address.StreetAddress,
			"city":           order.Address.City,
			"country":        order.Address.Country,
			"zip_code":       order.Address.ZipCode,
		}
	}

	if includeOrderItems {
		orderItems := make([]map[string]interface{}, len(order.OrderItems))
		for i := 0; i < len(order.OrderItems); i++ {
			oi := order.OrderItems[i]
			orderItems[i] = map[string]interface{}{
				"name":  oi.ProductName,
				"slug":  oi.Slug,
				"price": oi.Price,
			}
		}
		result["order_items"] = orderItems
	} else {
		result["order_items_count"] = order.OrderItemsCount
	}

	if includeUser {
		result["user"] = map[string]interface{}{
			"id":       order.UserId,
			"username": order.User.Username,
		}
	}

	return CreateSuccessDto(result)
}

func CreateOrderDetailsDto(order *models.Order) map[string]interface{} {
	// includeUser -> false
	// includeOrderItems -> true
	// includeUser -> false
	return CreateSuccessDto(CreateOrderDto(order, true, true, false))
}

func getIncludeFlags(includes ...bool) (includeAddress, includeOrderItems, includeUser bool) {

	if len(includes) > 0 {
		includeAddress = includes[0]
	}

	if len(includes) > 1 {
		includeOrderItems = includes[1]
	}

	if len(includes) > 2 {
		includeUser = includes[2]
	}
	return
}

func CreateOrderCreatedDto(order *models.Order) map[string]interface{} {
	return CreateSuccessWithDtoAndMessageDto(CreateOrderDetailsDto(order), "Order created successfully")
}
