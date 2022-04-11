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
		if helper.JsonRequested(c) {
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
