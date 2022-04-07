package main

import (
	"github.com/PurotoApp/authfox/internal/endpoints"
	"github.com/PurotoApp/authfox/internal/gormHelper"
	"github.com/PurotoApp/authfox/internal/redishelper"
	"github.com/PurotoApp/libpuroto/ginHelper"
	"github.com/gin-gonic/gin"
)

func main() {
	// connect to the PostgreSQL
	pg_conn := gormHelper.ConnectDB()
	// Connect to Redis
	redisVerify := redishelper.Connect(1)
	redisSession := redishelper.Connect(2)

	// migrate all tables
	endpoints.AutoMigrateAuthfox(pg_conn)

	// create router
	router := gin.Default()

	// configure gin
	ginHelper.ConfigRouter(router)

	// set routes
	endpoints.SetRoutes(router, pg_conn, redisVerify, redisSession)

	// start
	router.Run("localhost:3621")
}
