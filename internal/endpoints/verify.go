package endpoints

import (
	"context"
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/PurotoApp/authfox/internal/logHelper"
	"github.com/PurotoApp/authfox/internal/sessionHelper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type sendVerify struct {
	UserID     string `json:"uid"`
	Token      string `json:"token"`
	VerifyCode string `json:"verify_code"`
}

type saveInitProfile struct {
	NamePretty      string `bson:"name_pretty"`
	NameFormat      string `bson:"name_format"`
	NameStatic      string `bson:"name_static"`
	UserID          string `bson:"uid"`
	EMail           string `bson:"email"`
	BadgeBetaTester bool   `bson:"badge_beta_tester"`
}
type saveUserData struct {
	UserID       string    `bson:"uid"`
	Password     string    `bson:"password"`
	RegisterIP   string    `bson:"register_ip"`
	RegisterTime time.Time `bson:"register_time"`
}

func verifyUser(collVerifySession, collSession, collVerify, collProfiles, collUsers *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// only answer if content-type is set right
		if c.GetHeader("Content-Type") != "application/json" {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogEvent("authfox", "registerUser(): Received request with wrong Content-Type header")
			return
		}

		var sendVerifyStruct sendVerify

		// put the json into the struct
		if err := c.BindJSON(&sendVerifyStruct); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			logHelper.LogError("authfox", err)
			return
		}

		// check if the send values are valid
		if !checkVerifyStruct(&sendVerifyStruct) {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogEvent("authfox", "verifyUser(): Recived invalid data")
			return
		}

		valid, err := sessionHelper.SessionValid(&sendVerifyStruct.UserID, &sendVerifyStruct.Token, collVerifySession, collSession, true)

		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "verifyUser(): Received verification with non existent session")
			return
		} else if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		} else if !valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "verifyUser(): Received verification with invalid session")
			return
		}

		// retrieve user data
		verifyUserRaw := collVerify.FindOne(context.TODO(), bson.D{{Key: "uid", Value: sendVerifyStruct.UserID}})

		if verifyUserRaw.Err() != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", verifyUserRaw.Err())
			return
		}

		// decode data
		var localVerifyUser saveVerifyUser
		if err := verifyUserRaw.Decode(&localVerifyUser); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", verifyUserRaw.Err())
			return
		}

		// securely check if the verify roken is valid
		if subtle.ConstantTimeCompare([]byte(sendVerifyStruct.VerifyCode), []byte(localVerifyUser.VerifyCode)) != 1 {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "verifyUser(): Received verification with invalid Verify-Code")
			return
		}

		// create initial user profile
		var saveUserProfile saveInitProfile
		saveUserProfile.NamePretty = localVerifyUser.NameFormat
		saveUserProfile.NameFormat = localVerifyUser.NameFormat
		saveUserProfile.NameStatic = localVerifyUser.NameStatic
		saveUserProfile.UserID = localVerifyUser.UserID
		saveUserProfile.EMail = localVerifyUser.Email
		// Giving user the beta tester badge
		saveUserProfile.BadgeBetaTester = true
		// save into DB
		_, err = collProfiles.InsertOne(context.TODO(), saveUserProfile)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}

		// create initial user data
		var saveUserDataStruct saveUserData
		saveUserDataStruct.UserID = localVerifyUser.UserID
		saveUserDataStruct.Password = localVerifyUser.Password
		saveUserDataStruct.RegisterIP = localVerifyUser.RegisterIP
		saveUserDataStruct.RegisterTime = localVerifyUser.RegisterTime
		// save into DB
		_, err = collUsers.InsertOne(context.TODO(), saveUserDataStruct)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}

		// delete old data
		_, err = collVerify.DeleteOne(context.TODO(), bson.D{{Key: "uid", Value: sendVerifyStruct.UserID}})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", verifyUserRaw.Err())
			return
		}
		// delete old session
		_, err = collVerifySession.DeleteOne(context.TODO(), bson.D{{Key: "uid", Value: sendVerifyStruct.UserID}})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", verifyUserRaw.Err())
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
