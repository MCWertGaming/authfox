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

package endpoints

import (
	"net/http"

	"github.com/PurotoApp/libpuroto/libpuroto"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func SetRoutes(router *gin.Engine, pg_conn *gorm.DB, redisVerify, redisSession *redis.Client) {
	router.POST("/v1/user", registerUser(pg_conn, redisVerify, redisSession))
	router.POST("/v1/user/login", loginUser(pg_conn, redisVerify, redisSession))
	router.POST("/v1/user/verify", verifyUser(pg_conn, redisVerify, redisSession))
	router.POST("/v1/user/validate", validateSession(redisVerify, redisSession))
	router.PATCH("/v1/user", updatePassword(pg_conn, redisVerify, redisSession))
	// router.POST("/v1/user/delete", accountDeletion(collVerifySession, collSession, collUsers, collProfiles))
	// swagger docs
	router.Static("/swagger", "swagger/")
	// user redirects
	router.GET("/", libpuroto.Redirect("/swagger"))
	router.GET("/v1", libpuroto.Redirect("/swagger"))
}
