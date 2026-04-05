package v1

import "github.com/gin-gonic/gin"

func RegisterAuth(v1 *gin.RouterGroup, deps Deps) {
	auth := v1.Group("/auth")
	auth.POST("/send-otp", deps.Auth.GenerateAndSendOtp)
	auth.POST("/login", deps.Auth.LoginHandler)
	auth.POST("/reset-password", deps.Auth.ResetCustomerPassword)
	auth.POST("/generate-totp", deps.Auth.GenerateTOTP)
	auth.POST("/refresh", deps.Auth.Refresh)
	auth.POST("/logout", deps.Auth.Logout)
	auth.GET("/sessions", deps.Auth.ListSessions)
	auth.DELETE("/sessions/:sessionId", deps.Auth.RevokeSession)
}
