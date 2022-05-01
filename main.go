package main

import (
	"github.com/PurotoApp/authfox/internal/endpoints"
	"github.com/PurotoApp/authfox/internal/helper"
	"github.com/PurotoApp/libpuroto/ginHelper"
	"github.com/gin-gonic/gin"
)

func main() {
	// connect to the PostgreSQL
	pg_conn := helper.ConnectDB()
	// Connect to Redis
	redisVerify := helper.Connect(1)
	redisSession := helper.Connect(2)

	// migrate all tables
	endpoints.AutoMigrateAuthfox(pg_conn)

	// create router
	router := gin.Default()

	// configure gin
	ginHelper.ConfigRouter(router)

	// set routes
	endpoints.SetRoutes(router, pg_conn, redisVerify, redisSession)

	// start
	router.Run("0.0.0.0:3621")
}
