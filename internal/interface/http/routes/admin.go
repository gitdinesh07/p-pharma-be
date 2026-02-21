package routes

import "github.com/gin-gonic/gin"

func RegisterAdmin(admin *gin.RouterGroup, deps Deps) {
	orders := admin.Group("/orders")
	orders.PATCH("/:orderId/items/:itemId/status", deps.Order.UpdateOrderItemStatus)

	products := admin.Group("/products")
	products.POST("", deps.Product.CreateProduct)
	products.PATCH("/:productId", deps.Product.UpdateProduct)

	inventory := admin.Group("/inventory")
	inventory.PATCH("/:productId/stock", deps.Product.UpdateInventory)

	payments := admin.Group("/payments")
	payments.PATCH("/:paymentId/status", deps.Payment.UpdatePaymentStatus)

	consultations := admin.Group("/consultations")
	consultations.PATCH("/:consultationId/status", deps.Consultation.UpdateConsultationStatus)
}
