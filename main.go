/* <AuthFox - a simple authentication and session server for Puroto>
   Copyright (C) 2022  PurotoApp

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"github.com/PurotoApp/authfox/internal/endpoints"
	"github.com/PurotoApp/authfox/internal/helper"
	"github.com/PurotoApp/libpuroto/libpuroto"
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
	libpuroto.ConfigRouter(router)

	// set routes
	endpoints.SetRoutes(router, pg_conn, redisVerify, redisSession)

	// start
	router.Run("0.0.0.0:3621")
}
