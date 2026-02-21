package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"ppharma/backend/internal/config"
)

func adminToken(secret string) string {
	t := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, jwtv5.MapClaims{"sub": "admin_1", "role": "admin"})
	s, _ := t.SignedString([]byte(secret))
	return s
}

func customerToken(secret string) string {
	t := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, jwtv5.MapClaims{"sub": "cust_1", "role": "customer"})
	s, _ := t.SignedString([]byte(secret))
	return s
}

func TestInternalRequiresAPIKey(t *testing.T) {
	app, err := Build(config.Config{JWTSecret: "s", InternalAPIKey: []config.InternalAPIKeyConfig{{ID: "k1", Key: "secret", Scopes: []string{"orders.item_status.write"}}}})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/internal/orders/ord_1/items/item_1/status", bytes.NewBufferString(`{"status":"confirmed"}`))
	req.Header.Set("Authorization", "Bearer "+adminToken("s"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.Engine.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}

func TestAdminRouteRequiresJWT(t *testing.T) {
	app, err := Build(config.Config{JWTSecret: "s", InternalAPIKey: []config.InternalAPIKeyConfig{{ID: "k1", Key: "secret", Scopes: []string{"orders.item_status.write"}}}})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/orders/ord_1/items/item_1/status", bytes.NewBufferString(`{"status":"confirmed"}`))
	req.Header.Set("X-API-Key", "secret")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.Engine.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}

func TestCustomerGetOrder(t *testing.T) {
	app, err := Build(config.Config{JWTSecret: "s"})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/customer/orders/ord_1", nil)
	req.Header.Set("Authorization", "Bearer "+customerToken("s"))
	w := httptest.NewRecorder()
	app.Engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
