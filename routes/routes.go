package routes

import (
	handlers "recipe-service/handlers"
	"recipe-service/observability"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Route struct {
	db       *pgx.Conn
	handlers *handlers.Handlers
}

func NewRoute(db *pgx.Conn, prometheusMetrics *observability.PrometheusMetrics) *Route {
	return &Route{
		db:       db,
		handlers: handlers.NewHandlers(db, prometheusMetrics),
	}
}

func (r *Route) AddRecipeRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		recipes := v1.Group("/recipe")
		{
			recipes.GET("/:id", r.handlers.GetRecipe)
			recipes.GET("/list", r.handlers.ListRecipe)
			recipes.POST("/create", r.handlers.CreateRecipe)
			recipes.PUT("/:id", r.handlers.UpdateRecipe)
			recipes.DELETE("/:id", r.handlers.DeleteRecipe)
		}
	}
}

func (r *Route) AddHealthRoutes(router *gin.Engine) {
	// Health check endpoints (no authentication required)
	router.GET("/healthz", r.handlers.HealthzHandler)
	router.GET("/readyz", r.handlers.ReadyzHandler)
}
