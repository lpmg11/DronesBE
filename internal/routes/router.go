package router

import (
	"drones-be/api/v1/auth"
	"drones-be/internal/config"
	"drones-be/internal/middleware"
	"drones-be/internal/services"
	"drones-be/internal/storage"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router(cfn *config.Config, pg *storage.PostgresClient) *gin.Engine {

	authSrv := services.NewAuthService(pg, cfn)
	tokenSrv := services.NewTokenService(cfn)

	router := gin.Default()

	allowOriginsString := cfn.CorsOrigins
	allowOrigins := strings.Split(allowOriginsString, ";")

	router.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     allowOrigins,
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
			},
		),
	)

	router.GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("api/v1/")

	authHandler := auth.NewAuthHandler(authSrv, tokenSrv, cfn)
	auth.RegisterRoutes(v1, authHandler)

	protected := v1.Group("/protected")
	protected.Use(middleware.AuthMiddleware(tokenSrv))

	protected.GET("/ping", func(ctx *gin.Context) {
		userid, exist := ctx.Get("userID")
		if !exist {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"userID":  userid,
		})
	})

	return router

}
