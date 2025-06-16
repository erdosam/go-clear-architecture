package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// responses
type paginatedResponse struct {
	Page int `json:"page"`
	Size int `json:"size"`
	Data any `json:"data"`
}

type errorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// PaginatedHATEOASResponse For HATEOAS version
type PaginatedHATEOASResponse struct {
	Data  any          `json:"data"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
	Links HATEOASLinks `json:"_links"`
}

type HATEOASLinks struct {
	Self  string `json:"self"`
	First string `json:"first"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last"`
}

func ShouldBind[T any](c *gin.Context) (*T, error) {
	var body T //don't use pointer here, or it will panic on ShouldBind
	err := c.ShouldBind(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusPreconditionFailed, errorResponse{
			Error: "Invalid type in a field",
			Code:  0,
		})
	}
	return &body, err
}

func DetailJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func ListJSON[T any](c *gin.Context, data []T) {
	if data == nil {
		data = []T{}
	}
	c.JSON(http.StatusOK, paginatedResponse{
		Data: data,
		Page: 1, //TODO adjust with requested page
		Size: len(data),
	})
}

func ErrorJSON(c *gin.Context, s int, e error, code int) {
	c.AbortWithStatusJSON(s, errorResponse{
		Error: e.Error(),
		Code:  code,
	})
}
