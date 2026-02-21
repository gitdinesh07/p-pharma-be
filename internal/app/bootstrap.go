package app

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"ppharma/backend/internal/config"
	"ppharma/backend/internal/domain/common"
	"ppharma/backend/internal/domain/order"
	"ppharma/backend/internal/infra/auth/apikey"
	jwtinfra "ppharma/backend/internal/infra/auth/jwt"
	cachememory "ppharma/backend/internal/infra/cache/memory"
	zaplogger "ppharma/backend/internal/infra/logger/zap"
	"ppharma/backend/internal/interface/http/handlers"
	"ppharma/backend/internal/interface/http/middleware"
	"ppharma/backend/internal/interface/http/routes"
	repomemory "ppharma/backend/internal/repository/memory"
)

type Application struct {
	Engine *gin.Engine
	Config config.Config
}

func Build(cfg config.Config) (*Application, error) {
	logger, err := zaplogger.New("debug")
	if err != nil {
		return nil, err
	}

	_ = cachememory.New() // in-memory cache adapter wired and replaceable later

	seed := []*order.Order{{
		OrderID:       "ord_1",
		CustomerID:    "cust_1",
		Currency:      "INR",
		Subtotal:      2000,
		GrandTotal:    2100,
		ShippingFee:   100,
		PaymentStatus: "pending",
		DerivedStatus: order.OrderStatusPending,
		Items: []order.OrderItem{
			{ItemID: "item_1", ProductID: "prod_1", ProductSnapshot: order.ProductSnapshot{Name: "Med A", SKU: "MED-A", UnitPrice: 1000}, Qty: 1, LineTotal: 1000, Status: order.ItemStatusPending},
			{ItemID: "item_2", ProductID: "prod_2", ProductSnapshot: order.ProductSnapshot{Name: "Med B", SKU: "MED-B", UnitPrice: 1000}, Qty: 1, LineTotal: 1000, Status: order.ItemStatusPending},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}}

	repo := repomemory.NewOrderRepository(seed)
	service := order.NewService(repo, order.DefaultStatusDeriver{})
	orderHandler := handlers.NewOrderHandler(service)
	internalHandler := handlers.NewInternalHandler()
	authHandler := handlers.NewAuthHandler()
	productHandler := handlers.NewProductHandler()
	paymentHandler := handlers.NewPaymentHandler()
	consultationHandler := handlers.NewConsultationHandler()

	jwtProvider := jwtinfra.NewProvider(cfg.JWTSecret)

	var keySecrets []common.InternalAPIKeySecret
	for _, k := range cfg.InternalAPIKey {
		keySecrets = append(keySecrets, common.InternalAPIKeySecret{KeyID: k.ID, RawKey: k.Key, Scopes: k.Scopes})
	}
	if len(keySecrets) == 0 {
		keySecrets = append(keySecrets, common.InternalAPIKeySecret{KeyID: "default", RawKey: "internal-secret", Scopes: []string{"inventory.write", "orders.item_status.write", "products.write"}})
	}
	sp := apikey.NewStaticSecretProvider(keySecrets)
	keyAuth, err := apikey.NewAuthenticator(context.Background(), sp)
	if err != nil {
		return nil, err
	}

	engine := gin.New()
	engine.Use(gin.Recovery(), middleware.RequestLogger(logger))

	routeDeps := routes.Deps{
		Auth:         authHandler,
		Order:        orderHandler,
		Product:      productHandler,
		Payment:      paymentHandler,
		Consultation: consultationHandler,
		Internal:     internalHandler,
	}

	routes.RegisterHealth(engine)

	v1 := engine.Group("/api/v1")
	routes.RegisterAuth(v1, routeDeps)

	customer := v1.Group("/customer")
	customer.Use(middleware.JWTAuth(jwtProvider), middleware.RequireRole("customer"))
	routes.RegisterCustomer(customer, routeDeps)

	admin := v1.Group("/admin")
	admin.Use(middleware.JWTAuth(jwtProvider), middleware.RequireRole("admin"))
	routes.RegisterAdmin(admin, routeDeps)

	internal := v1.Group("/admin/internal")
	internal.Use(middleware.APIKeyAuth(keyAuth))
	routes.RegisterAdminInternal(internal, routeDeps)

	return &Application{Engine: engine, Config: cfg}, nil
}
