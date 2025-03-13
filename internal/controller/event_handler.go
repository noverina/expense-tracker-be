package controller

import (
	"expense-tracker/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath 	/api/v1
// @Security 	BearerAuth
// @Tags 		home
// @Produce 	json
// @Success 	200 	{string} 	pong
// @Router 		/ping [get]
func Ping(g *gin.Context) {
	g.JSON(http.StatusOK, "pong")
}

// @Tags 		event
// @Security 	BearerAuth
// @Summary		upsert event
// @Accept 		json
// @Produce 	json
// @Param 		event 	body 		api.Event 	true 	"event information"
// @Success 	200 	{object} 	api.HttpResponse
// @Failure 	400 	{object} 	api.HttpResponse
// @Failure 	500 	{object} 	api.HttpResponse
// @Router 		/event [post]
func UpsertEvent(c *gin.Context) {
	var event api.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	id, code, err := api.UpsertEvent(c, event)
	var response api.HttpResponse
	if err != nil {
		response = api.HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = api.HttpResponse{IsError: false, Message: "", Data: id}
	}
	c.JSON(code, response)
}

// @Tags 		event
// @Security 	BearerAuth
// @Summary		get events by filter
// @Description	please input filter criteria body must be in JSON format. example: {"_id": "123"}
// @Accept 		json
// @Produce 	json
// @Param 		filter 	body 		map[string]interface{} 	true 	"filter criteria in json format"
// @Success 	200 	{object} 	api.HttpResponse{data=[]api.Event}
// @Failure 	400 	{object} 	api.HttpResponse
// @Failure 	500 	{object} 	api.HttpResponse
// @Router 		/event/filter [post]
func GetEventByFilter(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	events, code, err := api.GetEventFilter(c, input)
	var response api.HttpResponse
	if err != nil {
		response = api.HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = api.HttpResponse{IsError: false, Message: "", Data: events}
	}
	c.JSON(code, response)
}

// @Tags 		event
// @Security 	BearerAuth
// @Summary		get all events in a given month
// @Accept 		json
// @Produce 	json
// @Param 		year 		query 		string 		true 	"year"
// @Param 		month 		query 		string 		true 	"month to filter with"
// @Param 		timezone 	query 		string 		true 	"timezone"
// @Success 	200 		{object} 	api.HttpResponse{data=[]api.Event}
// @Failure 	400 		{object} 	api.HttpResponse
// @Failure 	500 		{object} 	api.HttpResponse
// @Router 		/event/month [get]
func GetEventByMonth(c *gin.Context) {
	year := c.Query("year")
	month := c.Query("month")
	timezone := c.Query("timezone")

	events, code, err := api.GetEventByMonth(c, year, month, timezone)
	var response api.HttpResponse
	if err != nil {
		response = api.HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = api.HttpResponse{IsError: false, Message: "", Data: events}
	}
	c.JSON(code, response)
}

// @Tags 		event
// @Security 	BearerAuth
// @Summary		income + expense summary in a given month
// @Description	get income + expense summary for each category in a given month; also output total income + spending in said month
// @Accept 		json
// @Produce 	json
// @Param 		year 		query 		string 		true 	"year"
// @Param 		month 		query 		string 		true 	"month to filter with"
// @Param 		timezone 	query 		string 		true 	"timezone"
// @Success 	200 		{object} 	api.HttpResponse{data=[]api.Sum}
// @Failure 	400 		{object} 	api.HttpResponse
// @Failure 	500 		{object} 	api.HttpResponse
// @Router 		/event/sum [get]
func GetMonthSum(c *gin.Context) {
	year := c.Query("year")
	month := c.Query("month")
	timezone := c.Query("timezone")

	sum, code, err := api.GetMonthSum(c, year, month, timezone)
	var response api.HttpResponse
	if err != nil {
		response = api.HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = api.HttpResponse{IsError: false, Message: "", Data: sum}
	}
	c.JSON(code, response)
}
