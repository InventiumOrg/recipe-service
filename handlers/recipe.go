package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	models "recipe-service/models/sqlc"
	"recipe-service/observability"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Ingredient struct {
	Name     string  `json:"name" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required"`
	Measure  string  `json:"measure" binding:"required"`
	Unit     string  `json:"unit" binding:"required"`
}

type CreateRecipeRequest struct {
	Name        string       `json:"name" binding:"required"`
	Ingredients []Ingredient `json:"ingredients" binding:"required"`
	Cost        int32        `json:"cost" binding:"required"`
}

type UpdateRecipeRequest struct {
	Name        string       `json:"name" binding:"required"`
	Ingredients []Ingredient `json:"ingredients" binding:"required"`
	Cost        int32        `json:"cost" binding:"required"`
}

type Handlers struct {
	queries           *models.Queries
	tracer            trace.Tracer
	db                *pgx.Conn
	prometheusMetrics *observability.PrometheusMetrics
}

func NewHandlers(db *pgx.Conn, prometheusMetrics *observability.PrometheusMetrics) *Handlers {
	return &Handlers{
		db:                db,
		queries:           models.New(db),
		tracer:            otel.Tracer("recipe-service/handlers"),
		prometheusMetrics: prometheusMetrics,
	}
}

func (h *Handlers) GetRecipe(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "GetRecipe")
	defer span.End()

	// Get recipe ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid recipe ID",
		})
		return
	}

	// Add attributes to the span
	span.SetAttributes(attribute.Int64("recipe.id", id))

	dbStart := time.Now()
	recipe, err := h.queries.GetRecipe(ctx, id)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("get", "recipe", dbDuration, err)
	}

	if err != nil {
		slog.Error("Got an error while getting recipe", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get recipe",
		})
		return
	}

	// Record successful retrieval (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordRecipeOperation("get", recipe.Name)
	}

	// Record successful operation
	span.SetAttributes(
		attribute.String("recipe.name", recipe.Name),
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(200, gin.H{
		"message": "Get Recipe Successfully",
		"data":    recipe,
	})
}

func (h *Handlers) ListRecipe(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "ListRecipes")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int("recipe.limit", 10),
		attribute.Int("recipe.offset", 0),
	)

	dbStart := time.Now()
	recipes, err := h.queries.ListRecipe(ctx, models.ListRecipeParams{
		Limit:  10,
		Offset: 0,
	})
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("list", "recipe", dbDuration, err)
	}

	if err != nil {
		slog.Error("Got an error while listing recipes", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list recipes",
		})
		return
	}

	// Record successful list operation (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordRecipeOperation("list", "list")
	}

	span.SetAttributes(
		attribute.Int("recipe.count", len(recipes)),
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(200, gin.H{
		"message": "List Recipes Successfully",
		"data":    recipes,
		"count":   len(recipes),
	})
}

func (h *Handlers) CreateRecipe(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "CreateRecipe")
	defer span.End()

	var req CreateRecipeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredientsJSON, err := json.Marshal(req.Ingredients)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredients"})
		return
	}

	param := models.CreateRecipeParams{
		Name:        req.Name,
		Ingredients: json.RawMessage(ingredientsJSON),
		Cost:        req.Cost,
	}

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("recipe.name", req.Name),
		attribute.Int("recipe.cost", int(req.Cost)),
		attribute.Int("recipe.ingredients_count", len(req.Ingredients)),
	)

	dbStart := time.Now()
	recipe, err := h.queries.CreateRecipe(ctx, param)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("create", "recipe", dbDuration, err)
	}

	if err != nil {
		slog.Error("Could not create recipe", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create recipe",
		})
		return
	}

	// Record successful creation (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordRecipeOperation("create", recipe.Name)
		h.prometheusMetrics.UpdateRecipesCount(1)
	}

	// Record successful operation
	span.SetAttributes(
		attribute.Int64("recipe.id", recipe.ID),
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Create Recipe Successfully",
		"data":    recipe,
	})
}

func (h *Handlers) UpdateRecipe(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "UpdateRecipe")
	defer span.End()

	// Get recipe ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid recipe ID",
		})
		return
	}

	var req UpdateRecipeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredientsJSON, err := json.Marshal(req.Ingredients)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredients"})
		return
	}

	param := models.UpdateRecipeParams{
		ID:          id,
		Name:        req.Name,
		Ingredients: json.RawMessage(ingredientsJSON),
		Cost:        req.Cost,
	}

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int64("recipe.id", id),
		attribute.String("recipe.name", req.Name),
	)

	dbStart := time.Now()
	recipe, err := h.queries.UpdateRecipe(ctx, param)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("update", "recipe", dbDuration, err)
	}

	if err != nil {
		slog.Error("Could not update recipe", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update recipe",
		})
		return
	}

	// Record successful update (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordRecipeOperation("update", recipe.Name)
	}

	span.SetAttributes(attribute.String("operation.status", "success"))

	ctx.JSON(200, gin.H{
		"message": "Update Recipe Successfully",
		"data":    recipe,
	})
}

func (h *Handlers) DeleteRecipe(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "DeleteRecipe")
	defer span.End()

	// Get recipe ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid recipe ID",
		})
		return
	}

	// Add attributes to the span
	span.SetAttributes(attribute.Int64("recipe.id", id))

	dbStart := time.Now()
	err = h.queries.DeleteRecipe(ctx, id)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("delete", "recipe", dbDuration, err)
	}

	if err != nil {
		slog.Error("Could not delete recipe", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete recipe",
		})
		return
	}

	// Record successful deletion (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordRecipeOperation("delete", "deleted")
		h.prometheusMetrics.UpdateRecipesCount(-1)
	}

	span.SetAttributes(attribute.String("operation.status", "success"))

	ctx.JSON(200, gin.H{
		"message": "Delete Recipe Successfully",
	})
}
