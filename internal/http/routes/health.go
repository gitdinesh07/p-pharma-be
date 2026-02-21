package routes

import "github.com/gin-gonic/gin"

func RegisterHealth(engine *gin.Engine) {
	engine.GET("/health/live", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	engine.GET("/health/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })
}
