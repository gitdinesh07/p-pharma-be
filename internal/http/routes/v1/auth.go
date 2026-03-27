package v1

import "github.com/gin-gonic/gin"

func RegisterAuth(v1 *gin.RouterGroup, deps Deps) {
	auth := v1.Group("/auth")
	auth.POST("/customer/login", deps.Auth.CustomerLogin)
	auth.POST("/user/login", deps.Auth.UserLogin)
	auth.POST("/refresh", deps.Auth.Refresh)
	auth.POST("/logout", deps.Auth.Logout)
	auth.GET("/sessions", deps.Auth.ListSessions)
	auth.DELETE("/sessions/:sessionId", deps.Auth.RevokeSession)
}
