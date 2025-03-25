package router

import (
	admindrone "drones-be/api/v1/admin/drone"
	admintransactions "drones-be/api/v1/admin/transactions"
	adminwarehouse "drones-be/api/v1/admin/warehouse"
	"drones-be/api/v1/auth"
	providerproducts "drones-be/api/v1/provider/products"
	providerprofile "drones-be/api/v1/provider/profile"
	userbudget "drones-be/api/v1/user/budget"
	storeproducts "drones-be/api/v1/user/store/products"
	"drones-be/internal/config"
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

	warehouseHandler := adminwarehouse.NewWarehouseHandler(pg)
	adminwarehouse.RegisterRoutes(v1, warehouseHandler, tokenSrv)

	droneHandler := admindrone.NewDroneHandler(pg.DB)
	admindrone.RegisterRoutes(v1, droneHandler, tokenSrv)

	profileHandler := providerprofile.NewProfileHandler(pg.DB)
	providerprofile.RegisterRoutes(v1, profileHandler, tokenSrv)

	providerProductsHandler := providerproducts.NewProductHandler(pg.DB)
	providerproducts.RegisterRoutes(v1, providerProductsHandler, tokenSrv)

	productsHandler := storeproducts.NewProductHandler(pg.DB)
	storeproducts.RegisterRoutes(v1, productsHandler, tokenSrv)

	budgetHandler := userbudget.NewBudgetHandler(pg.DB)
	userbudget.RegisterRoutes(v1, budgetHandler, tokenSrv)

	adminTransactionsHandler := admintransactions.NewTransactionHandler(pg.DB)
	admintransactions.RegisterRoutes(v1, adminTransactionsHandler, tokenSrv)

	return router

}
