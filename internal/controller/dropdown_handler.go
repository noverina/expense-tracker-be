package controller

import (
	"expense-tracker/internal/api"

	"github.com/gin-gonic/gin"
)

// @Tags 		dropdown
// @Security 	BearerAuth
// @Summary		type: income / expense
// @Accept 		json
// @Produce 	json
// @Success 	200 	{object} 	api.HttpResponse{data=[]api.Dropdown}
// @Failure 	500 	{object} 	api.HttpResponse
// @Router 		/dropdown/type [get]
func GetTypes(c *gin.Context) {
	data, code := api.GetTypes()
	response := api.HttpResponse{IsError: false, Message: "", Data: data}
	c.JSON(code, response)
}

// @Tags		dropdown
// @Security 	BearerAuth
// @Accept 		json
// @Produce 	json
// @Success 	200 	{object} 	api.HttpResponse{data=[]api.Dropdown}
// @Failure 	500 	{object} 	api.HttpResponse
// @Router 		/dropdown/expense [get]
func GetExpenses(c *gin.Context) {
	data, code := api.GetExpenses()
	response := api.HttpResponse{IsError: false, Message: "", Data: data}
	c.JSON(code, response)
}

// @Tags 		dropdown
// @Security 	BearerAuth
// @Accept 		json
// @Produce 	json
// @Success 	200 	{object} 	api.HttpResponse{data=[]api.Dropdown}
// @Failure 	500 	{object} 	api.HttpResponse
// @Router 		/dropdown/income [get]
func GetIncomes(c *gin.Context) {
	data, code := api.GetIncomes()
	response := api.HttpResponse{IsError: false, Message: "", Data: data}
	c.JSON(code, response)
}
