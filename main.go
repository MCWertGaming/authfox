package main

import (
	authfox "github.com/PurotoApp/authfox/internal"
	"github.com/PurotoApp/authfox/internal/ginHelper"
	loghelper "github.com/PurotoApp/authfox/internal/logHelper"
	"github.com/PurotoApp/authfox/internal/mongoHelper"
	"github.com/gin-gonic/gin"
)

func main() {
	// create DB connection
	client, err := mongoHelper.ConnectDB(mongoHelper.GetDBUri())
	loghelper.ErrorFatal(err)
	// test the connection
	loghelper.ErrorFatal(mongoHelper.TestDBConnection(client))
	// close connection on program exit
	defer func() {
		loghelper.ErrorFatal(mongoHelper.DisconnectDB(client))
	}()

	// set up gin
	ginHelper.SwitchRelMode()

	// create router
	router := gin.Default()

	// configure gin
	ginHelper.ConfigRouter(router)

	// set routes
	authfox.SetRoutes(router, client)

	// start
	router.Run("localhost:3621")
}
