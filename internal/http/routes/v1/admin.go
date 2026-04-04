package v1

import (
	"ppharma/backend/internal/domain/user"
	"ppharma/backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdmin(admin *gin.RouterGroup, deps Deps) {
	// SuperAdmin / GlobalAdmin ONLY
	users := admin.Group("/users")
	users.Use(middleware.RequireRole(string(user.RoleSuperAdmin), string(user.RoleGlobalAdmin)))
	users.POST("", deps.User.CreateUser)
	users.PUT("/:id", deps.User.UpdateUser)

	// Available to Admin, SuperAdmin and GlobalAdmin (due to global bootstrap definition)
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
