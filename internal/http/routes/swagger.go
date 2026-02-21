package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppharma/backend/internal/http/docs"
)

func RegisterSwagger(engine *gin.Engine) {
	engine.GET("/swagger/openapi.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/yaml", []byte(docs.OpenAPISpec()))
	})

	engine.GET("/swagger", func(c *gin.Context) {
		html := `<!doctype html>
<html>
<head>
  <meta charset="utf-8" />
  <title>PPharma Swagger</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: '/swagger/openapi.yaml',
      dom_id: '#swagger-ui'
    });
  </script>
</body>
</html>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})
}
