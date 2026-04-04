package middleware

import (
	"ppharma/backend/internal/domain/common"
	"strings"

	"github.com/gin-gonic/gin"
)

const ClientInfoKey = "client_info"

// ClientInfo extracts tracking headers securely from generic requests abstracting generic device metrics.
func ClientInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		info := common.ClientAppInfo{
			DeviceID: c.GetHeader("X-Device-Id"),
			Source:   c.GetHeader("X-Source"),
		}

		appInfo := c.GetHeader("X-App-Info")
		if appInfo != "" && strings.Contains(appInfo, ",") {
			parts := strings.Split(appInfo, ",")
			if len(parts) == 2 {
				info.DeviceType = parts[0]
				info.AppVersion = parts[1]
			}
		}

		if info.DeviceType == "" {
			info.DeviceType = "android"
		}
		// Inject info object for downstream controller boundaries.
		c.Set(ClientInfoKey, info)

		c.Next()
	}
}
