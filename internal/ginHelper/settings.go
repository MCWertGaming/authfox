package ginHelper

import (
	"os"

	"github.com/PurotoApp/authfox/internal/logHelper"
	"github.com/gin-gonic/gin"
)

func ConfigRouter(router *gin.Engine) {

	if os.Getenv("GIN_MODE") == "release" {
		// turn on proxy support
		// TODO: allow users to specify trusted proxies
		// TODO: what if proxy behind proxy
		// TODO: what if no value specified
		logHelper.ErrorFatal("Router", router.SetTrustedProxies(nil))
	} else {
		// turn off proxy support for debugging
		logHelper.ErrorFatal("Router", router.SetTrustedProxies(nil))
	}
	// set health status route
	router.GET("/health", getHealth)
}
