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
	"crypto/subtle"
	"net/http"

	"github.com/PurotoApp/libpuroto/libpuroto"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type sendVerify struct {
	UserID     string `json:"uid"`
	Token      string `json:"token"`
	VerifyCode string `json:"verify_code"`
}

func verifyUser(pg_conn *gorm.DB, redisVerify, redisSession *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// only answer if content-type is set right
		if libpuroto.JsonRequested(c) {
			return
		}

		var sendVerifyStruct sendVerify

		// put the json into the struct
		if err := c.BindJSON(&sendVerifyStruct); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			libpuroto.LogError("authfox", err)
			return
		}

		// check if the send values are valid
		if !checkVerifyStruct(&sendVerifyStruct) {
			c.AbortWithStatus(http.StatusBadRequest)
			libpuroto.LogEvent("authfox", "verifyUser(): Recived invalid data")
			return
		}

		valid, err := libpuroto.SessionValid(&sendVerifyStruct.UserID, &sendVerifyStruct.Token, redisVerify)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", err)
			return
		} else if !valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			libpuroto.LogEvent("authfox", "verifyUser(): Received verification with invalid session")
			return
		}

		// retrieve user data
		var verifyData Verify
		if err := pg_conn.Where("user_id = ?", sendVerifyStruct.UserID).Take(&verifyData).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", err)
			return
		}

		// securely check if the verify token is valid
		if subtle.ConstantTimeCompare([]byte(sendVerifyStruct.VerifyCode), []byte(verifyData.VerifyCode)) != 1 {
			c.AbortWithStatus(http.StatusUnauthorized)
			libpuroto.LogEvent("authfox", "verifyUser(): Received verification with invalid Verify-Code")
			return
		}

		// create initial user profile
		var userProfile Profile
		userProfile.NamePretty = verifyData.NameFormat
		userProfile.NameFormat = verifyData.NameFormat
		userProfile.NameStatic = verifyData.NameStatic
		userProfile.UserID = verifyData.UserID
		userProfile.Email = verifyData.Email
		// Giving user the beta tester badge
		userProfile.BadgeBetaTester = true
		userProfile.BadgeAlphaTester = true
		userProfile.BadgeStaff = false
		// save into DB
		if err = pg_conn.Create(&userProfile).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", err)
			return
		}

		// create initial user data
		var userData User
		userData.UserID = verifyData.UserID
		userData.Password = verifyData.Password
		userData.RegisterIP = verifyData.RegisterIP
		userData.RegisterTime = verifyData.RegisterTime
		// save into DB
		if err = pg_conn.Create(&userData).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", err)
			return
		}

		// delete old data
		if err = pg_conn.Delete(&verifyData).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", err)
			return
		}

		// delete old session
		if err = redisVerify.Del(sendVerifyStruct.UserID).Err(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", err)
			return
		}

		c.Status(http.StatusAccepted)
	}
}

// returns false if the struct holds empty values
func checkVerifyStruct(verifyStruct *sendVerify) bool {
	if verifyStruct.Token == "" || verifyStruct.UserID == "" || verifyStruct.VerifyCode == "" {
		return false
	}
	return true
}
