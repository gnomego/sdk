package xgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJson(c *gin.Context, data interface{}) {
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid_request", err.Error()))
		return
	}
}

func BindTargetJson(c *gin.Context, data interface{}, target string) {
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, NewTargetedErrorResponse(target, "invalid_request", err.Error()))
		return
	}
}

func Respond(c *gin.Context, response ResponseInfo) {
	if response.IsOk() {
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusInternalServerError, response)
}
