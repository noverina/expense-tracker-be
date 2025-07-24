package main

import (
	"expense-tracker/docs"
	"expense-tracker/internal/api"
	"expense-tracker/internal/controller"
	"fmt"
	"log"
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
	if os.Getenv("GIN_MODE") != "release" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, skipping")
		}
	}	
	
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("unable to load .env file err=%w", err)
	}
	address := os.Getenv("FRONTEND_ADDRESS")

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{address},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.Use(ErrorHandler())
	r.SetTrustedProxies(nil)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api.InitDB()
	api.InitAuth()
	api.InitEvent()
	api.InitLog()
	defer api.Disconnect()

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		home := v1.Group("")
		{
			home.GET("/ping", controller.Ping)

		}
		auth := v1.Group("/auth")
		{
			auth.POST("", api.RoleAuthMiddleware("client"), controller.GenerateToken)
			auth.GET("", api.RoleAuthMiddleware("admin"), controller.InvalidateToken)
		}
		event := v1.Group("/event")
		event.Use(api.RoleAuthMiddleware("client"))
		{
			event.POST("", controller.UpsertEvent)
			event.POST("/filter", controller.GetEventByFilter)
			event.GET("/month", controller.GetEventByMonth)
			event.GET("/sum", controller.GetMonthSum)
		}
		dropdown := v1.Group("/dropdown")
		dropdown.Use(api.RoleAuthMiddleware("client"))
		{
			dropdown.GET("/type", controller.GetTypes)
			dropdown.GET("/expense", controller.GetExpenses)
			dropdown.GET("/income", controller.GetIncomes)
		}
	}

	port := os.Getenv("PORT")
	r.Run("0.0.0.0:" + port)
}
