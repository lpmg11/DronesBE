package router

import (
	"drones-be/api/v1/auth"
	"drones-be/internal/config"
	"drones-be/internal/services"
	"drones-be/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router(cfn *config.Config, pg *storage.PostgresClient) *gin.Engine {

	authSrv := services.NewAuthService(pg, cfn)
	tokenSrv := services.NewTokenService(cfn)

	router := gin.Default()
	router.GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("api/v1/")

	authHandler := auth.NewAuthHandler(authSrv, tokenSrv)
	auth.RegisterRoutes(v1, authHandler)

	return router

}
