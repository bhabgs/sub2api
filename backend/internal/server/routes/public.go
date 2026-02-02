package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes registers public routes (no authentication required)
func RegisterPublicRoutes(v1 *gin.RouterGroup, h *handler.Handlers) {
	public := v1.Group("/public")
	{
		public.GET("/usage", h.Public.Usage)
	}
}
