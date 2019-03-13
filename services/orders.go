package services

import (
	"github.com/melardev/GoGonicEcommerceApi/infrastructure"
	"github.com/melardev/GoGonicEcommerceApi/models"
)

func FetchOrdersPage(userId uint, page, pageSize int) (orders []models.Order, totalOrdersCount int, err error) {
	database := infrastructure.GetDb()

	totalOrdersCount = 0

	query := database.Model(&models.Order{}).Where(&models.Order{UserId: userId})
	query.Count(&totalOrdersCount)

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).
		// TODO: Why Preload("Address") does not work?, perhaps OrderItems does
		// Preload("OrderItems").Preload("Address").
		Find(&orders).Error
	if err != nil {
		return
	}

	var orderIds = make([]uint, len(orders))
	for i := 0; i < len(orders); i++ {
		orderIds[i] = orders[i].ID
	}

	var orderItems []models.OrderItem
	if len(orders) > 0 {
		//
		database.Select("id, order_id").Where("order_id in (?)", orderIds).Find(&orderItems)

		for i := 0; i < len(orderItems); i++ {
			oi := orderItems[i]
			for j := 0; j < len(orders); j++ {
				if oi.OrderId == orders[j].ID {
					orders[j].OrderItemsCount = orders[j].OrderItemsCount + 1
				}
			}
		}
	}
	return orders, totalOrdersCount, err
}

func FetchOrderDetails(orderId uint) (order models.Order, err error) {
	database := infrastructure.GetDb()
	err = database.Model(models.Order{}).Preload("OrderItems").First(&order, orderId).Error
	var address models.Address
	database.Model(&order).Related(&address)
	order.Address = address
	return order, err
}
