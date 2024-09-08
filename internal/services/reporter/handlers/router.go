package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/repository"
	httplib "github.com/rocky2015aaa/ethdefender/pkg/http"
)

// API contains generic methods that should be used by other API interfaces.
type API interface {
	Setup(router *gin.Engine)
}

func NewRouter(db repository.Storage) http.Handler {
	var router *gin.Engine
	if config.Default.Gin.Mode == gin.DebugMode {
		router = gin.Default()
	} else {
		router = gin.New()
	}
	// Setup service routers
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// Setup middlewares
	router.Use(httplib.CORSMiddleware())
	// Setup API routes
	NewReportsAPI(db).Setup(router)
	return router
}

func (api *ReportsAPI) Setup(router *gin.Engine) {
	go func() {
		api.CreateReports(context.Background())
	}()
	router.GET("/api/v1/report/transaction", api.GetTransactionReport)
	router.GET("/api/v1/report/pause", api.GetPauseReport)
	router.GET("/api/v1/report/slither", api.GetSlitherReport)
	router.POST("/api/v1/report/slither", api.PostSlitherReport)
}
