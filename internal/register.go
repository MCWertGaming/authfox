package authfox

import (
	"net/http"
	"strings"
	"time"

	loghelper "github.com/PurotoApp/authfox/internal/logHelper"
	"github.com/PurotoApp/authfox/internal/security"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// This struct stores user informations send to the register API endpoint
type sendUserProfile struct {
	// the send user name
	NameFormat string `json:"user_name"`
	// Plain-text password send to the DB
	Password string `json:"password"`
	// account email
	Email string `json:"email"`
}

type saveVerifyUser struct {
	NameFormat   string    `bson:"name_format"`
	NameStatic   string    `bson:"name_static"`
	UserID       string    `bson:"uid"`
	Email        string    `bson:"email"`
	Password     string    `bson:"password"`
	RegisterIP   string    `bson:"register_ip"`
	RegisterTime time.Time `bson:"register_time"`
	VerifyCode   string    `bson:"verify_code"`
}

// This is like sessionPair but without the session type switch
type returnSession struct {
	UserID string `json:"uid"`
	Token  string `json:"token"`
}

func registerUser(collUsers, collVerify, collSession, collVerifySession *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// only answer if content-type is set right
		if c.GetHeader("Content-Type") != "application/json" {
			c.AbortWithStatus(http.StatusBadRequest)
			loghelper.LogEvent("authfox", "registerUser(): Received request with wrong Content-Type header")
			return
		}

		var sendUserStruct sendUserProfile

		// put the json into the struct
		if err := c.BindJSON(&sendUserStruct); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			loghelper.LogError("authfox", err)
			return
		}
		// make sure that the received values are legal
		if !checkSendUserProfile(&sendUserStruct) {
			c.AbortWithStatus(http.StatusBadRequest)
			loghelper.LogEvent("authfox", "registerUser(): Received invalid or illegal registration data")
			return
		}

		// prepare saving of user data into verify DB
		var userData saveVerifyUser

		// hash the password
		hash, err := security.CreateHash(sendUserStruct.Password)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			loghelper.LogError("authfox", err)
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
		if userData.VerifyCode, err = security.RandomString(32); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			loghelper.LogError("authfox", err)
		}
		// create user ID
		if userData.UserID, err = generateUserID(collUsers, collVerify); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			loghelper.LogError("authfox", err)
			return
		}
		// store into DB
		addVerifyUser(userData, collVerify)
		// create session
		session, err := createSession(userData.UserID, collSession, collVerifySession, true)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			loghelper.LogError("authfox", err)
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
func checkSendUserProfile(profile *sendUserProfile) bool {
	// TODO: refuse if the name is not between 3-32 characters
	// TODO: refuse if the name is already used
	// TODO: refuse if the name contains slurs / forbidden words
	// TODO: refuse if @ is used
	// TODO: don't allow special characters
	if profile.NameFormat == "" {
		return false
	}
	// TODO: refuse if the email is in invalid format
	// TODO: refuse if the email address is forbidden (trashmail etc)
	// TODO: refuse if localhost is used
	// TODO: refuse if the email is already used
	// TODO: check how long the password can be before it breaks the hash
	if profile.Email == "" {
		return false
	}
	// TODO: refuse if the password is under 8 chars
	// TODO: refuse on weak passwords
	if profile.Password == "" {
		return false
	}
	return true
}
