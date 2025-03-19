package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	h := &handler{}

	// router.Use(middleware.ErrorHandler())

	// health check
	router.GET("/healthz", h.ServerStatus)
}
