package app

import (
	"context"
	"time"

	"ppharma/backend/internal/config"
	"ppharma/backend/internal/domain/common"
	"ppharma/backend/internal/domain/order"
	"ppharma/backend/internal/http/handlers"
	"ppharma/backend/internal/http/middleware"
	"ppharma/backend/internal/http/routes"
	routesv1 "ppharma/backend/internal/http/routes/v1"
	repomemory "ppharma/backend/internal/repository/memory"
	appservice "ppharma/backend/internal/service"
	"ppharma/backend/support-pkg/auth/apikey"
	jwtinfra "ppharma/backend/support-pkg/auth/jwt"
	cachememory "ppharma/backend/support-pkg/cache/memory"
	zaplogger "ppharma/backend/support-pkg/logger/zap"
	"ppharma/backend/support-pkg/queue/filequeue"

	"github.com/gin-gonic/gin"
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
	queue := filequeue.New(cfg.QueueDir)
	internalHandler := handlers.NewInternalHandler(queue, cfg.QueueTopic)
	// Mock repo injections (to be wired into real mongo)
	authService := appservice.NewAuthService(nil, nil, jwtinfra.NewProvider(cfg.JWTSecret))
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler()
	paymentHandler := handlers.NewPaymentHandler()
	consultationHandler := handlers.NewConsultationHandler()

	jwtProvider := jwtinfra.NewProvider(cfg.JWTSecret)

	var keySecrets []common.InternalAPIKeySecret
	for _, k := range cfg.InternalAPIKey {
		keySecrets = append(keySecrets, common.InternalAPIKeySecret{KeyID: k.ID, RawKey: k.Key, Scopes: k.Scopes})
	}
	if len(keySecrets) == 0 {
		keySecrets = append(keySecrets, common.InternalAPIKeySecret{KeyID: "default", RawKey: "internal-secret", Scopes: []string{"inventory.write", "orders.item_status.write", "products.write", "queue.write"}})
	}
	sp := apikey.NewStaticSecretProvider(keySecrets)
	keyAuth, err := apikey.NewAuthenticator(context.Background(), sp)
	if err != nil {
		return nil, err
	}

	engine := gin.New()
	engine.Use(gin.Recovery(), middleware.RequestLogger(logger))

	routeDeps := routesv1.Deps{
		Auth:         authHandler,
		Order:        orderHandler,
		Product:      productHandler,
		Payment:      paymentHandler,
		Consultation: consultationHandler,
		Internal:     internalHandler,
	}

	routes.RegisterHealth(engine)
	if cfg.AppEnv != "production" {
		routes.RegisterSwagger(engine)
	}

	apiV1 := engine.Group("/api/v1")
	routesv1.RegisterAuth(apiV1, routeDeps)

	customer := apiV1.Group("/customer")
	customer.Use(middleware.JWTAuth(jwtProvider), middleware.RequireRole("customer"))
	routesv1.RegisterCustomer(customer, routeDeps)

	admin := apiV1.Group("/admin")
	admin.Use(middleware.JWTAuth(jwtProvider), middleware.RequireRole("admin"))
	routesv1.RegisterAdmin(admin, routeDeps)

	internal := apiV1.Group("/admin/internal")
	internal.Use(middleware.APIKeyAuth(keyAuth))
	routesv1.RegisterAdminInternal(internal, routeDeps)

	return &Application{Engine: engine, Config: cfg}, nil
}
