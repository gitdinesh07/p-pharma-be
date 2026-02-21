package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"ppharma/backend/internal/domain/common"
)

type fakeAPIKeyAuth struct{}

func (f fakeAPIKeyAuth) Authenticate(_ context.Context, rawKey string) (*common.APIKeyPrincipal, error) {
	if rawKey == "ok" {
		return &common.APIKeyPrincipal{KeyID: "k1", Scopes: map[string]struct{}{"orders.item_status.write": {}}}, nil
	}
	return nil, gin.Error{}
}

func TestAPIKeyScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(APIKeyAuth(fakeAPIKeyAuth{}), RequireScope("orders.item_status.write"))
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-API-Key", "ok")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
