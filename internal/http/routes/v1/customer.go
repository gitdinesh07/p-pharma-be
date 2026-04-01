package v1

import "github.com/gin-gonic/gin"

func RegisterCustomerPublic(customer *gin.RouterGroup, deps Deps) {
	customer.POST("", deps.Customer.CreateCustomer)
}

func RegisterCustomer(customer *gin.RouterGroup, deps Deps) {
	customer.GET("/:id", deps.Customer.GetCustomer)
	customer.PUT("/:id", deps.Customer.UpdateCustomer)

	orders := customer.Group("/orders")
	orders.GET("/:orderId", deps.Order.GetCustomerOrder)
	orders.GET("/:orderId/payments", deps.Payment.ListOrderPayments)

	products := customer.Group("/products")
	products.GET("", deps.Product.ListProducts)
	products.GET("/:productId", deps.Product.GetProduct)

	consultations := customer.Group("/consultations")
	consultations.POST("", deps.Consultation.CreateConsultation)
	consultations.GET("", deps.Consultation.ListConsultations)
	consultations.GET("/:consultationId", deps.Consultation.GetConsultation)
}
