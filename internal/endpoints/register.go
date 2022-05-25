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
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PurotoApp/authfox/internal/helper"
	"github.com/PurotoApp/libpuroto/libpuroto"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// This struct stores user informations send to the register API endpoint
type sendUserProfile struct {
	// the send user name
	NameFormat string `json:"user_name"`
	// Plain-text password received from client
	Password string `json:"password"`
	// account email
	Email string `json:"email"`
}

// This is like sessionPair but without the session type switch
type returnSession struct {
	UserID string `json:"uid"`
	Token  string `json:"token"`
}

func registerUser(pg_conn *gorm.DB, redisVerify, redisSession *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// only answer if content-type is set right
		if helper.JsonRequested(c) {
			return
		}

		var sendUserStruct sendUserProfile

		// put the json into the struct
		if err := c.BindJSON(&sendUserStruct); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			libpuroto.LogError("authfox", err)
			return
		}
		// make sure that the received values are legal
		if !checkSendUserProfile(&sendUserStruct) {
			c.AbortWithStatus(http.StatusBadRequest)
			libpuroto.LogEvent("authfox", "registerUser(): Received invalid or illegal registration data")
			return
		}

		// check if the given email or user name already exists
		result := pg_conn.Where("name_static = ?", strings.ToLower(sendUserStruct.NameFormat)).Where("email = ?", strings.ToLower(sendUserStruct.Email)).Take(&Verify{})
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", result.Error)
			return
		} else if result.RowsAffected > 0 {
			c.AbortWithStatus(http.StatusBadRequest)
			libpuroto.LogEvent("authfox", "Received user that already exists")
			return
		}

		// prepare saving of user data into verify DB
		var userData Verify

		// hash the password
		hash, err := helper.CreateHash(&sendUserStruct.Password)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", err)
			return
		}
		// safe the hashed password
		userData.Password = hash
		// remove the old password from memory
		sendUserStruct.Password = ""

		// fill other user data
		userData.NameFormat = sendUserStruct.NameFormat
		userData.NameStatic = strings.ToLower(sendUserStruct.NameFormat)
		userData.Email = strings.ToLower(sendUserStruct.Email)
		userData.RegisterIP = c.ClientIP()
		userData.RegisterTime = time.Now()
		if userData.VerifyCode, err = helper.RandomString(32); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", err)
		}
		// create user ID
		userData.UserID = uuid.New().String()

		// store into DB
		if pg_conn.Create(&userData).Error != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", err)
			return
		}

		// create session
		session, err := helper.CreateSession(&userData.UserID, redisVerify, redisSession, true)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			libpuroto.LogError("authfox", err)
			return
		}

		// remove the session type as it is always true
		var basicSession returnSession
		basicSession.UserID = session.UserID
		basicSession.Token = session.Token

		c.JSON(http.StatusAccepted, basicSession)
	}
}

// check the send user data for correctnes and forbidden values
// TODO: move into Guardian service
func checkSendUserProfile(profile *sendUserProfile) bool {
	// TODO: refuse if the name contains slurs / forbidden words
	// TODO: don't allow special characters:
	if strings.Count(profile.NameFormat, "") < 6 || strings.Count(profile.NameFormat, "") > 32 ||
		strings.Count(profile.NameFormat, "@") > 0 || strings.Count(profile.NameFormat, " ") > 0 {
		return false
	}

	// TODO: refuse if the email address is forbidden (trashmail etc)
	if profile.Email == "" || !libpuroto.CheckEmail(profile.Email) {
		return false
	}

	// TODO: refuse on weak passwords
	if strings.Count(profile.Password, "") < 9 || len(profile.Password) > 512 {
		return false
	}
	return true
}
