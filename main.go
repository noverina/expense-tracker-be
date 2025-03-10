package main

import (
	"expense-tracker/docs"
	"expense-tracker/internal/api"
	"expense-tracker/internal/controller"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				api.LogError(fmt.Sprintf("%v", r))
				c.JSON(http.StatusInternalServerError, api.HttpResponse{
					IsError: true,
					Message: "Uh oh! :( something unexpected happened",
					Data:    nil,
				})
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			if !c.IsAborted() {
				for _, value := range c.Errors {
					api.LogError(fmt.Sprintf("%v", value))
				}

				c.JSON(
					http.StatusInternalServerError,
					api.HttpResponse{
						IsError: true,
						Message: "Uh oh! :( something unexpected happened",
						Data:    nil},
				)
			}

			c.Abort()
		}
	}
}

// @title           Expense Tracker API
// @version         1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description ! IMPORTANT ! Please prepend Bearer manually. Example: "Bearer {token}"
// @BasePath  /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("unable to load .env file err=%w", err)
	}
	address := os.Getenv("FRONTEND_ADDRESS")

	api.InitDB()
	api.InitAuth()
	api.InitEvent()
	api.InitLog()

	defer api.Disconnect()

	r := gin.New()
	r.Use(ErrorHandler())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{address},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
	}))
	r.SetTrustedProxies(nil)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		home := v1.Group("")
		{
			home.GET("/ping", controller.Ping)

		}
		auth := v1.Group("/auth")
		{
			//TODO disable register endpoint later
			auth.POST("", controller.Register)
			auth.POST("/access", controller.Access)
			auth.GET("/refresh", api.JWTAuthMiddleware(), controller.Refresh)
		}
		event := v1.Group("/event")
		event.Use(api.JWTAuthMiddleware())
		{
			event.POST("", controller.UpsertEvent)
			event.POST("/filter", controller.GetEventByFilter)
			event.GET("/month", controller.GetEventByMonth)
			event.GET("/sum", controller.GetMonthSum)
		}
		category := v1.Group("/dropdown")
		category.Use(api.JWTAuthMiddleware())
		{
			category.GET("/type", controller.GetTypes)
			category.GET("/expense", controller.GetExpenses)
			category.GET("/income", controller.GetIncomes)
		}
	}

	cert := os.Getenv("CERT")
	certKey := os.Getenv("CERT_KEY")
	err := r.RunTLS(":443", cert, certKey)
	if err != nil {
		api.LogError("failed to start https", "err", err)
		os.Exit(1)
	}
	httpRedirect := gin.Default()
	httpRedirect.GET("/*any", func(c *gin.Context) {
		host := os.Getenv("HOST_ADDRESS")
		c.Redirect(http.StatusMovedPermanently, host+c.Request.RequestURI)
	})

	if err := httpRedirect.Run(":80"); err != nil {
		api.LogError("failed to start http", "err", err)
		os.Exit(1)
	}
}
