package router

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nalawade41/secret-server/config"
	_ "github.com/nalawade41/secret-server/docs"
	"github.com/nalawade41/secret-server/internal/wire"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handler struct {
	config    *config.Config
	dbConnect *dynamodb.Client
}

func NewHandler(cfg *config.Config, db *dynamodb.Client) *Handler {
	return &Handler{config: cfg, dbConnect: db}
}

func (h *Handler) Init() *echo.Echo {
	// Init echo router
	e := echo.New()
	e.Use(
		middleware.Recover(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: os.Stdout,
		}),
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(200)),
		middleware.CORS(),
	)

	// Init router
	e.GET("/", HealthCheck)

	// Show swagger docs if APP_ENV is not production
	if h.config.Environment != config.Prod {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Init open API routes
	h.initAPI(e)

	return e
}

func (h *Handler) initAPI(e *echo.Echo) {
	secretManager := wire.InitializeRouteProvider(h.dbConnect, h.config.Database.TableName)
	api := e.Group("/api/v1")
	{
		secretManager.InitRoutes(api)
	}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Server Health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
