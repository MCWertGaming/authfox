package authfox

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"github.com/PurotoApp/authfox/internal/security"
	"github.com/PurotoApp/libpuroto/logHelper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrAccountNotExisting = errors.New("findUserData(): Given account does not exist")
)

type sendLogin struct {
	LoginName string `json:"login"`
	Password  string `json:"password"`
}
type savedUserData struct {
	NameStatic string `bson:"name_static"`
	UserID     string `bson:"uid"`
	Email      string `bson:"email"`
	Password   string `bson:"password"`
}
type userID struct {
	UserID string `bson:"uid"`
}

func loginUser(collUser, collSession, collVerifySession, collVerify, collProfiles *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// only answer if content-type is set right
		if c.GetHeader("Content-Type") != "application/json" {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogEvent("authfox", "loginUser(): Received request with wrong Content-Type header")
			return
		}
		var sendLoginStruct sendLogin

		if err := c.BindJSON(&sendLoginStruct); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogError("authfox", err)
			return
		}

		// check the data for validity
		if !checkLoginData(sendLoginStruct) {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogEvent("authfox", "loginUser(): Invalid data recieved")
			return
		}
		// find user
		userData, verify, err := findUserData(collUser, collVerify, collProfiles, sendLoginStruct.LoginName)
		// check if the given user not existed
		if err == ErrAccountNotExisting {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "loginUser(): Received login for non existing user")
			return
			// check for internal error
		} else if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}

		// decode DB data
		var localUserData savedUserData
		if err := userData.Decode(&localUserData); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}

		// check if the password matches the stored one
		match, err := security.ComparePasswordAndHash(sendLoginStruct.Password, localUserData.Password)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}
		if !match {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "loginUser(): Invalid password received")
			return
		}

		// create session
		session, err := createSession(localUserData.UserID, collSession, collVerifySession, verify)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}

		// return session
		c.JSON(http.StatusAccepted, session)
	}
}

// returns false if the login struct includes valid data
func checkLoginData(loginData sendLogin) bool {
	if loginData.LoginName == "" {
		return false
	}
	if loginData.Password == "" {
		return false
	}
	return true
}

// returns true if the given string is an email
func checkEmail(value string) bool {
	_, err := mail.ParseAddress(value)
	return err == nil
}

func findUserData(collUser, collVerify, collProfiles *mongo.Collection, login string) (userData *mongo.SingleResult, verify bool, err error) {
	// set the search parameter
	var loginType string
	if checkEmail(strings.ToLower(login)) {
		loginType = "email"
	} else {
		loginType = "name_static"
	}

	// find user profile
	userProfile := collProfiles.FindOne(context.TODO(), bson.D{{Key: loginType, Value: strings.ToLower(login)}})
	// check if a profile was found
	if userProfile.Err() == mongo.ErrNoDocuments {
		// user was not found in user DB, check the verify DB
		userData = collVerify.FindOne(context.TODO(), bson.D{{Key: loginType, Value: strings.ToLower(login)}})
		// check if a user was found this time
		if userData.Err() == mongo.ErrNoDocuments {
			// user does not excist
			return &mongo.SingleResult{}, true, ErrAccountNotExisting
		} else if userData.Err() != nil {
			return &mongo.SingleResult{}, true, userData.Err()
		}
		// valid data was found as it seems!
		return userData, true, nil
	} else if userProfile.Err() != nil {
		return &mongo.SingleResult{}, true, userProfile.Err()
	}
	// A profile was found, encoding it to get the UserID
	var userIDStruct userID
	if err := userProfile.Decode(&userIDStruct); err != nil {
		return &mongo.SingleResult{}, true, err
	}
	// get the data we need
	userData = collUser.FindOne(context.TODO(), bson.D{{Key: "uid", Value: userIDStruct.UserID}})
	if userData.Err() == mongo.ErrNoDocuments {
		// user does not excist
		return &mongo.SingleResult{}, true, ErrAccountNotExisting
	} else if userData.Err() != nil {
		return &mongo.SingleResult{}, true, userData.Err()
	}
	// valid data was found as it seems!
	return userData, false, nil
}
