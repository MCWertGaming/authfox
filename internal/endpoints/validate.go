/*
<AuthFox - a simple authentication and session server for Puroto>
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

	"github.com/PurotoApp/authfox/internal/helper"
	"github.com/PurotoApp/libpuroto/logHelper"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type sendSession struct {
	UserID string `json:"uid"`
	Token  string `json:"token"`
}

func validateSession(redisVerify, redisSession *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// only answer if content-type is set right
		if !helper.JsonRequested(c) {
			return
		}

		var sendSessionStruct sendSession

		// put the json into the struct
		if err := c.BindJSON(&sendSessionStruct); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			logHelper.LogError("authfox", err)
			return
		}

		valid, err := helper.SessionValid(&sendSessionStruct.UserID, &sendSessionStruct.Token, redisSession)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		} else if !valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "Received invalid session")
			return
		}
		c.Status(http.StatusOK)
	}
}
