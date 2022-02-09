package endpoints

import (
	"context"
	"net/http"

	"github.com/PurotoApp/authfox/internal/logHelper"
	"github.com/PurotoApp/authfox/internal/security"
	"github.com/PurotoApp/authfox/internal/sessionHelper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type sendRemoveData struct {
	UserID   string `json:"uid"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

func accountDeletion(collVerifySession, collSession, collUser, collProfile *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// about on incorrect request-header
		if c.GetHeader("Content-Type") != "application/json" {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogEvent("authfox", "registerUser(): Received request with wrong Content-Type header")
			return
		}

		var sendDataStruct sendRemoveData

		// put the json into the struct
		if err := c.BindJSON(&sendDataStruct); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			logHelper.LogError("authfox", err)
			return
		}

		// validate session
		valid, err := sessionHelper.SessionValid(&sendDataStruct.UserID, &sendDataStruct.Token, collVerifySession, collSession, false)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		} else if !valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "Received invalid session")
			return
		}

		// validate password
		// get the data we need
		userData := collUser.FindOne(context.TODO(), bson.D{{Key: "uid", Value: sendDataStruct.UserID}})
		if userData.Err() != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", userData.Err())
			return
		}
		// decode data
		var passwordLocal passwordData
		if err := userData.Decode(&passwordLocal); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", userData.Err())
			return
		}
		// compare passwords
		match, err := security.ComparePasswordAndHash(sendDataStruct.Password, passwordLocal.Password)
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

		// remove sessions
		_, err = collSession.DeleteMany(context.TODO(), bson.D{{Key: "uid", Value: sendDataStruct.UserID}})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogEvent("authfox", "Invalid password received")
			return
		}
		// remove user
		_, err = collUser.DeleteOne(context.TODO(), bson.D{{Key: "uid", Value: sendDataStruct.UserID}})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogEvent("authfox", "Invalid password received")
			return
		}
		// remove profile
		_, err = collProfile.DeleteOne(context.TODO(), bson.D{{Key: "uid", Value: sendDataStruct.UserID}})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "Invalid password received")
			return
		}
		// TODO: remove posts

		c.Status(http.StatusAccepted)
	}
}
