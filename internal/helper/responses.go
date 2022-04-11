package helper

import (
	"net/http"

	"github.com/PurotoApp/libpuroto/logHelper"
	"github.com/gin-gonic/gin"
)

// returns true, if the client requested json format, also sets the response to 406, if not
func JsonRequested(c *gin.Context) bool {
	if c.GetHeader("Content-Type") != "application/json" {
		c.AbortWithStatus(http.StatusNotAcceptable)
		logHelper.LogEvent("authfox", "Received request with wrong Content-Type header")
		return false
	}
	return true
}
