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
	// create collections
	collUsers := client.Database("authfox").Collection("users")
	collVerify := client.Database("authfox").Collection("verify")
	collSession := client.Database("authfox").Collection("session")
	collVerifySession := client.Database("authfox").Collection("verifySession")

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
	authfox.SetRoutes(router, collUsers, collVerify, collSession, collVerifySession)

	// start
	router.Run("localhost:3621")
}
