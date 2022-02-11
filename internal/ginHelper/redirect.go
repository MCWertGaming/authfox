package ginHelper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// redirects to the given url
func Redirect(url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, url)
	}
}
