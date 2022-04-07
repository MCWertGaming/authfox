package endpoints

import (
	"errors"
	"net/http"

	"github.com/PurotoApp/authfox/internal/security"
	"github.com/PurotoApp/authfox/internal/sessionHelper"
	"github.com/PurotoApp/libpuroto/logHelper"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type sendUpdateData struct {
	UserID      string `json:"uid"`
	Token       string `json:"token"`
	PasswordOld string `json:"password_old"`
	PasswordNew string `json:"password_new"`
}

func updatePassword(pg_conn *gorm.DB, redisVerify, redisSession *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// about on incorrect request-header
		if c.GetHeader("Content-Type") != "application/json" {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogEvent("authfox", "registerUser(): Received request with wrong Content-Type header")
			return
		}

		var sendDataStruct sendUpdateData

		// put the json into the struct
		if err := c.BindJSON(&sendDataStruct); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			logHelper.LogError("authfox", err)
			return
		}

		// validate session
		valid, err := sessionHelper.SessionValid(&sendDataStruct.UserID, &sendDataStruct.Token, redisVerify, redisSession, false)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		} else if !valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "Received invalid session")
			return
		}

		// validate old password
		// get the hashed password
		localPass, err := findUserPassword(pg_conn, sendDataStruct.UserID)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}
		// compare passwords
		match, err := security.ComparePasswordAndHash(sendDataStruct.PasswordOld, localPass)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}
		if !match {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "Invalid password received")
			return
		}

		// update password
		// TODO recycle hash
		newPassHash, err := security.CreateHash(sendDataStruct.PasswordNew)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}
		// save new pass
		pg_conn.Model(&User{UserID: sendDataStruct.UserID}).Update("password", newPassHash)

		c.Status(http.StatusAccepted)
	}
}
func findUserPassword(pg_conn *gorm.DB, userID string) (string, error) {
	var localUser User
	res := pg_conn.Where("user_id = ?", userID).Take(&localUser)
	if res.Error != nil {
		return "", res.Error
	} else if res.RowsAffected != 1 {
		return "", errors.New("invalid numbers of rows found while searching for user password")
	}
	return localUser.Password, nil
}
