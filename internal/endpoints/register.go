package endpoints

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PurotoApp/authfox/internal/security"
	"github.com/PurotoApp/authfox/internal/sessionHelper"
	"github.com/PurotoApp/libpuroto/logHelper"
	"github.com/PurotoApp/libpuroto/stringHelper"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrReceivedUserThatExists = errors.New("checkUserExists(): Received a user that already exists")

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

func registerUser(collUsers, collVerify, collSession, collVerifySession, collProfiles *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// only answer if content-type is set right
		if c.GetHeader("Content-Type") != "application/json" {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogEvent("authfox", "registerUser(): Received request with wrong Content-Type header")
			return
		}

		var sendUserStruct sendUserProfile

		// put the json into the struct
		if err := c.BindJSON(&sendUserStruct); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			logHelper.LogError("authfox", err)
			return
		}
		// make sure that the received values are legal
		if !checkSendUserProfile(&sendUserStruct) {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogEvent("authfox", "registerUser(): Received invalid or illegal registration data")
			return
		}

		// check if the given email or user name already exists
		exists, err := checkUserExists(sendUserStruct.NameFormat, sendUserStruct.Email, collVerify, collProfiles)
		if err == ErrReceivedUserThatExists || exists {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogError("authfox", err)
			return
		} else if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}

		// prepare saving of user data into verify DB
		var userData saveVerifyUser

		// hash the password
		hash, err := security.CreateHash(sendUserStruct.Password)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
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
			logHelper.LogError("authfox", err)
		}
		// create user ID
		userData.UserID = uuid.New().String()

		// store into DB
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		_, err = collVerify.InsertOne(ctx, userData)
		cancel()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}
		// create session
		session, err := sessionHelper.CreateSession(userData.UserID, collSession, collVerifySession, true)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
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
	// TODO: refuse if the name contains slurs / forbidden words
	// TODO: don't allow special characters:
	if strings.Count(profile.NameFormat, "") < 6 || strings.Count(profile.NameFormat, "") > 32 ||
		strings.Count(profile.NameFormat, "@") > 0 || strings.Count(profile.NameFormat, " ") > 0 {
		return false
	}

	// TODO: refuse if the email address is forbidden (trashmail etc)
	if profile.Email == "" || !stringHelper.CheckEmail(profile.Email) {
		return false
	}

	// TODO: refuse on weak passwords
	if strings.Count(profile.Password, "") < 9 || len(profile.Password) > 512 {
		return false
	}
	return true
}

// returns true if a user exists with the given name
func checkUserExists(name, email string, collVerify, collProfiles *mongo.Collection) (bool, error) {
	// count users with the given name
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	count, err := collVerify.CountDocuments(ctx, bson.M{"name_static": strings.ToLower(name)}, options.Count().SetLimit(1))
	cancel()
	if err != nil {
		return true, err
	}
	if count != 0 {
		return true, ErrReceivedUserThatExists
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
	count, err = collProfiles.CountDocuments(ctx, bson.M{"name_static": strings.ToLower(name)}, options.Count().SetLimit(1))
	cancel()
	if err != nil {
		return true, err
	}
	if count != 0 {
		return true, ErrReceivedUserThatExists
	}

	// count users with the given email
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
	count, err = collVerify.CountDocuments(ctx, bson.D{{Key: "email", Value: strings.ToLower(email)}}, options.Count().SetLimit(1))
	cancel()
	if err != nil {
		return true, err
	}
	if count != 0 {
		return true, ErrReceivedUserThatExists
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
	count, err = collProfiles.CountDocuments(ctx, bson.D{{Key: "email", Value: strings.ToLower(email)}}, options.Count().SetLimit(1))
	cancel()
	if err != nil {
		return true, err
	}
	if count != 0 {
		return true, ErrReceivedUserThatExists
	}

	// no account with the given values exists
	return false, nil
}
