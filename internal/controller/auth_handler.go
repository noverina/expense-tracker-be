package controller

import (
	"expense-tracker/internal/api"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TODO disable register endpoint
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Param 		client 		body 		api.AuthAccess	true 	"client information"
// @Success 	200 		{object} 	api.HttpResponse
// @Failure 	400 		{object} 	api.HttpResponse
// @Failure 	500 		{object} 	api.HttpResponse
// @Router 		/auth [post]
func Register(c *gin.Context) {
	var body api.AuthAccess
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	code, err := api.Register(c, body.Identifier, body.SecretKey)
	var response api.HttpResponse
	if err != nil {
		response = api.HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = api.HttpResponse{IsError: false, Message: "", Data: nil}
	}
	c.JSON(code, response)
}

// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Param 		client 		body 		api.AuthAccess	true 	"client information"
// @Success 	200 		{object} 	api.HttpResponse{data=[]api.AuthInfo}
// @Failure 	400 		{object} 	api.HttpResponse
// @Failure 	500 		{object} 	api.HttpResponse
// @Router 		/auth/access [post]
func Access(c *gin.Context) {
	var body api.AuthAccess
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	tokens, code, err := api.Access(c, body.Identifier, body.SecretKey)
	var response api.HttpResponse
	if err != nil {
		response = api.HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = api.HttpResponse{IsError: false, Message: "", Data: tokens}
	}
	c.JSON(code, response)
}

// @Tags 		auth
// @Security 	BearerAuth
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	api.HttpResponse
// @Failure 	400 		{object} 	api.HttpResponse
// @Failure 	500 		{object} 	api.HttpResponse
// @Router 		/auth/refresh [get]
func Refresh(c *gin.Context) {
	token := c.GetHeader("Authorization")
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, api.HttpResponse{
			IsError: true,
			Message: "auth must be bearer",
		})
		return
	}

	token, code, err := api.Refresh(c, tokenParts[1])
	var response api.HttpResponse
	if err != nil {
		response = api.HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = api.HttpResponse{IsError: false, Message: "", Data: token}
	}
	c.JSON(code, response)
}
