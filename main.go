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
		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(
				http.StatusInternalServerError, 
				controller.HttpResponse{
					IsError: true, 
					Message: "Uh oh! :( something unexpected happened", 
					Data: nil},
				)
		}
	}
}

// @title           Expense Tracker API
// @version         1.0

// @BasePath  /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("unable to load .env file err=%w", err)
	}
	address := os.Getenv("FRONTEND_ADDRESS")

	defer api.Disconnect()

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(ErrorHandler())
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{address}, 
        AllowMethods:     []string{"GET", "POST", "DELETE"}, 
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, 
        AllowCredentials: false, 
    }))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		home := v1.Group("")
		{
			home.GET("/ping", controller.Ping)
		}
		event := v1.Group("/event")
		{
			event.POST("", controller.UpsertEvent)
			event.POST("/filter", controller.GetEvent)
		}
		category := v1.Group("/dropdown")
		{
			category.GET("/type", controller.GetTypes)
			category.GET("/expense", controller.GetExpenses)
			category.GET("/income", controller.GetIncomes)
		}
	}

	r.Run(":8083")
}