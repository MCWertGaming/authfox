package ginHelper

import (
	"os"

	"github.com/PurotoApp/libpuroto/logHelper"
	"github.com/gin-gonic/gin"
)

func SwitchRelMode() {
	// switch to release mode
	// TODO: Do only in prod
	if os.Getenv("RELEASE_TYPE") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
func ConfigRouter(router *gin.Engine) {

	if os.Getenv("CASCADE_RELEASE") == "production" {
		// turn on proxy support
		// TODO: allow users to specify trusted proxies
		// TODO: what if proxy behind proxy
		// TODO: what if no value specified
		logHelper.ErrorFatal(router.SetTrustedProxies(nil))
	} else {
		// turn off proxy support for debugging
		logHelper.ErrorFatal(router.SetTrustedProxies(nil))
	}
	// set health status route
	router.GET("/health", getHealth)
}
