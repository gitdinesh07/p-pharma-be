package v1

import "ppharma/backend/internal/http/handlers"

type Deps struct {
	Auth         *handlers.AuthHandler
	Order        *handlers.OrderHandler
	Product      *handlers.ProductHandler
	Payment      *handlers.PaymentHandler
	Consultation *handlers.ConsultationHandler
	Internal     *handlers.InternalHandler
	Customer     *handlers.CustomerHandler
}

