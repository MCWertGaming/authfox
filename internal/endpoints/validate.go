package endpoints

import (
	"net/http"

	"github.com/PurotoApp/authfox/internal/logHelper"
	"github.com/PurotoApp/authfox/internal/sessionHelper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type sendSession struct {
	UserID string `json:"uid"`
	Token  string `json:"token"`
}

func validateSession(collVerifySession, collSession *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// about on incorrect request-header
		if c.GetHeader("Content-Type") != "application/json" {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogEvent("authfox", "registerUser(): Received request with wrong Content-Type header")
			return
		}

		var sendSessionStruct sendSession

		// put the json into the struct
		if err := c.BindJSON(&sendSessionStruct); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			logHelper.LogError("authfox", err)
			return
		}

		valid, err := sessionHelper.SessionValid(&sendSessionStruct.UserID, &sendSessionStruct.Token, collVerifySession, collSession, false)
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
