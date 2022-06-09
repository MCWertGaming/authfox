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
	"github.com/PurotoApp/libpuroto/libpuroto"
	"github.com/gin-gonic/gin"
)

func main() {
	// connect to the PostgreSQL
	pg_conn := libpuroto.ConnectDB()
	if pg_conn.Error != nil {
		libpuroto.ErrorPanic(pg_conn.Error)
	}
	// Connect to Redis
	redisVerify := libpuroto.Connect(1)
	redisSession := libpuroto.Connect(2)

	// check if redis can be reached
	if err := redisVerify.Ping().Err(); err != nil {
		libpuroto.ErrorPanic(err)
	} else if err := redisSession.Ping().Err(); err != nil {
		libpuroto.ErrorPanic(err)
	}

	// migrate all tables
	endpoints.AutoMigrateAuthfox(pg_conn)

	// create router
	router := gin.Default()

	// configure gin
	libpuroto.ConfigRouter(router)

	// set routes
	endpoints.SetRoutes(router, pg_conn, redisVerify, redisSession)

	// start
	if err := router.Run("0.0.0.0:3621"); err != nil {
		libpuroto.ErrorPanic(err)
	}

	// clean up
	redisVerify.Close()
	redisSession.Close()
}
