package endpoints

import (
	"context"
	"net/http"
	"time"

	"github.com/PurotoApp/authfox/internal/security"
	"github.com/PurotoApp/authfox/internal/sessionHelper"
	"github.com/PurotoApp/libpuroto/logHelper"
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		userData := collUser.FindOne(ctx, bson.D{{Key: "uid", Value: sendDataStruct.UserID}})
		cancel()
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
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
		_, err = collSession.DeleteMany(ctx, bson.D{{Key: "uid", Value: sendDataStruct.UserID}})
		cancel()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogEvent("authfox", "Invalid password received")
			return
		}
		// remove user
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
		_, err = collUser.DeleteOne(ctx, bson.D{{Key: "uid", Value: sendDataStruct.UserID}})
		cancel()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogEvent("authfox", "Invalid password received")
			return
		}
		// remove profile
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
		_, err = collProfile.DeleteOne(ctx, bson.D{{Key: "uid", Value: sendDataStruct.UserID}})
		cancel()
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "Invalid password received")
			return
		}
		// TODO: remove posts

		c.Status(http.StatusAccepted)
	}
}
