package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"ppharma/backend/internal/domain/auth"
	"ppharma/backend/internal/domain/common"
	"ppharma/backend/pkg/api"
)

func JWTAuth(tp common.TokenProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if !strings.HasPrefix(authz, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.APIResponse[any]{
				Success: false,
				Error:   &api.APIError{Code: "UNAUTHORIZED", Message: "missing bearer token"},
			})
			return
		}
		token := strings.TrimPrefix(authz, "Bearer ")
		principal, err := tp.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.APIResponse[any]{
				Success: false,
				Error:   &api.APIError{Code: "UNAUTHORIZED", Message: "invalid token"},
			})
			return
		}
		c.Set(string(auth.ContextPrincipalKey), principal)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(string(auth.ContextPrincipalKey))
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.APIResponse[any]{
				Success: false,
				Error:   &api.APIError{Code: "UNAUTHORIZED", Message: "principal missing"},
			})
			return
		}
		p, ok := v.(*common.Principal)
		if !ok || p.Role != role {
			c.AbortWithStatusJSON(http.StatusForbidden, api.APIResponse[any]{
				Success: false,
				Error:   &api.APIError{Code: "FORBIDDEN", Message: "insufficient role"},
			})
			return
		}
		c.Next()
	}
}

func APIKeyAuth(authenticator common.APIKeyAuthenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("X-API-Key")
		principal, err := authenticator.Authenticate(c.Request.Context(), raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.APIResponse[any]{
				Success: false,
				Error:   &api.APIError{Code: "UNAUTHORIZED", Message: "invalid api key"},
			})
			return
		}
		c.Set(string(auth.ContextAPIKeyKey), principal)
		c.Next()
	}
}

func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(string(auth.ContextAPIKeyKey))
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.APIResponse[any]{
				Success: false,
				Error:   &api.APIError{Code: "UNAUTHORIZED", Message: "api key principal missing"},
			})
			return
		}
		p, ok := v.(*common.APIKeyPrincipal)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.APIResponse[any]{
				Success: false,
				Error:   &api.APIError{Code: "UNAUTHORIZED", Message: "invalid api key principal"},
			})
			return
		}
		if _, ok := p.Scopes[scope]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, api.APIResponse[any]{
				Success: false,
				Error:   &api.APIError{Code: "FORBIDDEN", Message: "missing scope"},
			})
			return
		}
		c.Next()
	}
}
