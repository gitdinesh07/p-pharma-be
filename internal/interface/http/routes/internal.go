package routes

import (
	"github.com/gin-gonic/gin"
	"ppharma/backend/internal/interface/http/middleware"
)

func RegisterAdminInternal(internal *gin.RouterGroup, deps Deps) {
	inventory := internal.Group("/inventory")
	inventory.POST("/sync", middleware.RequireScope("inventory.write"), deps.Internal.SyncInventory)

	orders := internal.Group("/orders")
	orders.PATCH("/:orderId/items/:itemId/status", middleware.RequireScope("orders.item_status.write"), deps.Order.UpdateOrderItemStatus)

	products := internal.Group("/products")
	products.POST("/bulk-upsert", middleware.RequireScope("products.write"), deps.Internal.BulkUpsertProducts)
}
