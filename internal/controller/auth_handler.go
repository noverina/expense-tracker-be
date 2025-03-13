package controller

import (
	"expense-tracker/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags 		auth
// @Summary		generate token
// @Accept 		json
// @Produce 	json
// @Param 		client 		body 		api.Auth	true 	"client information"
// @Success 	200 		{object} 	api.HttpResponse{data=string}
// @Failure 	400 		{object} 	api.HttpResponse
// @Failure 	500 		{object} 	api.HttpResponse
// @Router 		/auth 		[post]
func GenerateToken(c *gin.Context) {
	var body api.Auth
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	token, code, err := api.GenerateToken(c, body.Identifier, body.SecretKey)
	var response api.HttpResponse
	if err != nil {
		response = api.HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = api.HttpResponse{IsError: false, Message: "", Data: token}
	}

	c.JSON(code, response)
}

// @Tags 		auth
// @Security 	BearerAuth
// @Summary		invalidate token
// @Accept 		json
// @Produce 	json
// @Param 		token 		query 		string 		true 	"the token string to invalidate"
// @Success 	200 		{object} 	api.HttpResponse{data=string}
// @Failure 	400 		{object} 	api.HttpResponse
// @Failure 	500 		{object} 	api.HttpResponse
// @Router 		/auth 		[get]
func InvalidateToken(c *gin.Context) {
	token := c.Query("token")

	code, err := api.InvalidateToken(c, token)
	var response api.HttpResponse
	if err != nil {
		response = api.HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = api.HttpResponse{IsError: false, Message: "", Data: nil}
	}

	c.JSON(code, response)
}
