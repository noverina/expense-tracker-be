package controller

import (
	"expense-tracker/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpResponse struct {
	IsError bool        `json:"is_error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// @BasePath 	/api/v1

// @Tags 		home
// @Produce 	json
// @Success 	200 	{string} 	pong
// @Router 		/ping [get]
func Ping(g *gin.Context)  {
	g.JSON(http.StatusOK,"pong")
}

// @Tags 		event
// @Accept 		json
// @Produce 	json
// @Param 		event 	body 		api.Event 	true 	"event information"
// @Success 	200 	{object} 	HttpResponse 
// @Failure 	400 	{object} 	HttpResponse
// @Failure 	500 	{object} 	HttpResponse
// @Router 		/event [post]
func UpsertEvent(c *gin.Context) {
	var event api.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	id, code, err := api.UpsertEvent(c, event)
	var response HttpResponse
	if (err != nil) {
		response = HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = HttpResponse{IsError: false, Message: "", Data: id}
	}
	c.JSON(code, response)
}

// @Tags 		event
// @Accept 		json
// @Produce 	json
// @Param 		filter 	body 		map[string]interface{} 	true 	"filter criteria in json format"
// @Success 	200 	{object} 	HttpResponse{data=[]api.Event}
// @Failure 	400 	{object} 	HttpResponse
// @Failure 	500 	{object} 	HttpResponse
// @Router 		/event/filter [post]
func GetEventByFilter(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	events, code, err := api.GetEventFilter(c, input)
	var response HttpResponse
	if (err != nil) {
		response = HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = HttpResponse{IsError: false, Message: "", Data: events}
	}
	c.JSON(code, response)
}

// @Tags 		event
// @Accept 		json
// @Produce 	json
// @Param 		year 		query 		string 		true 	"year"
// @Param 		month 		query 		string 		true 	"month to filter with"
// @Param 		timezone 	query 		string 		true 	"timezone"
// @Success 	200 		{object} 	HttpResponse{data=[]api.Event}
// @Failure 	400 		{object} 	HttpResponse
// @Failure 	500 		{object} 	HttpResponse
// @Router 		/event/month [get]
func GetEventByMonth(c *gin.Context) {
	year := c.Query("year")
	month := c.Query("month")
	timezone := c.Query("timezone")

	events, code, err := api.GetEventByMonth(c, year, month, timezone)
	var response HttpResponse
	if (err != nil) {
		response = HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = HttpResponse{IsError: false, Message: "", Data: events}
	}
	c.JSON(code, response)
}


// @Tags 		category
// @Accept 		json
// @Produce 	json
// @Success 	200 	{object} 	HttpResponse{data=[]api.Dropdown}
// @Failure 	500 	{object} 	HttpResponse
// @Router 		/dropdown/type [get]
func GetTypes(c *gin.Context) {
	data, code := api.GetTypes()
	response := HttpResponse{IsError: false, Message: "", Data: data}
	c.JSON(code, response)
}

// @Tags		category
// @Accept 		json
// @Produce 	json
// @Success 	200 	{object} 	HttpResponse{data=[]api.Dropdown}
// @Failure 	500 	{object} 	HttpResponse
// @Router 		/dropdown/expense [get]
func GetExpenses(c *gin.Context) {
	data, code := api.GetExpenses()
	response := HttpResponse{IsError: false, Message: "", Data: data}
	c.JSON(code, response)
}

// @Tags 		category
// @Accept 		json
// @Produce 	json
// @Success 	200 	{object} 	HttpResponse{data=[]api.Dropdown}
// @Failure 	500 	{object} 	HttpResponse
// @Router 		/dropdown/income [get]
func GetIncomes(c *gin.Context) {
	data, code := api.GetIncomes()
	response := HttpResponse{IsError: false, Message: "", Data: data}
	c.JSON(code, response)
}