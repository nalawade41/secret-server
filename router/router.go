package router

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nalawade41/secret-server/config"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handler struct {
	ContextTimeout time.Duration
	Config         *config.Config
	// TODO: add any additional required values
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{Config: cfg}
}

func (h *Handler) Init() *echo.Echo {
	// Init echo router
	e := echo.New()
	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(200)),
		middleware.CORS(),
	)

	// Init router
	e.GET("/", HealthCheck)

	// Show swagger docs if APP_ENV is not production
	if h.Config.Environment != config.Prod {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Init open API routes
	h.initAPI(e)

	return e
}

func (h *Handler) initAPI(e *echo.Echo) {
	// TODO: add API routes
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
