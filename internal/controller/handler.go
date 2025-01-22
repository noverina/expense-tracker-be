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

// @BasePath /api/v1

// @Tags home
// @Accept json
// @Produce json
// @Success 200 {string} pong
// @Router /ping [get]
func Ping(g *gin.Context)  {
	g.JSON(http.StatusOK,"pong")
}

// @Accept json
// @Produce json
// @Param event body api.Event true "event information"
// @Success 200 {object} HttpResponse 
// @Failure 400 {object} HttpResponse
// @Failure 500 {object} HttpResponse
// @Router /event [post]
func UpsertEvent(c *gin.Context) {
	var event api.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	id, err := api.UpsertEvent(event)
	var response HttpResponse
	if (err != nil) {
		response = HttpResponse{IsError: true, Message: err.Error(), Data: nil}
		c.JSON(http.StatusInternalServerError, response)
	} else {
		response = HttpResponse{IsError: false, Message: "", Data: id}
		c.JSON(http.StatusOK, response)
	}
}

// @Accept json
// @Produce json
// @Param filter body map[string]interface{} true "filter string; MUST be in json format"
// @Success 200 {object} HttpResponse{data=[]api.Event}
// @Failure 400 {object} HttpResponse
// @Failure 500 {object} HttpResponse
// @Router /event/filter [post]
func GetEvent(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	events, err, code := api.GetEventFilter(input)
	var response HttpResponse
	if (err != nil) {
		response = HttpResponse{IsError: true, Message: err.Error(), Data: nil}
	} else {
		response = HttpResponse{IsError: false, Message: "", Data: events}
	}
	c.JSON(code, response)
}
