package endpoints

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PurotoApp/authfox/internal/security"
	"github.com/PurotoApp/authfox/internal/sessionHelper"
	"github.com/go-redis/redis"
	"gorm.io/gorm"

	"github.com/PurotoApp/libpuroto/logHelper"
	"github.com/PurotoApp/libpuroto/stringHelper"
	"github.com/gin-gonic/gin"
)

var (
	ErrAccountNotExisting = errors.New("findUserData(): Given account does not exist")
)

type sendLogin struct {
	LoginName string `json:"login"`
	Password  string `json:"password"`
}

func loginUser(pg_conn *gorm.DB, redisVerify, redisSession *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// only answer if content-type is set right
		if c.GetHeader("Content-Type") != "application/json" {
			c.AbortWithStatus(http.StatusNotAcceptable)
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
		if !checkLoginData(&sendLoginStruct) {
			c.AbortWithStatus(http.StatusBadRequest)
			logHelper.LogEvent("authfox", "loginUser(): Invalid data recieved")
			return
		}
		// find user
		localPassword, localUserID, verify, err := findUserData(pg_conn, &sendLoginStruct.LoginName)
		if err == ErrAccountNotExisting {
			// account does not exist
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "Received login for non existent account")
			return
		} else if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		}

		// check if the password matches the stored one
		match, err := security.ComparePasswordAndHash(&sendLoginStruct.Password, &localPassword)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logHelper.LogError("authfox", err)
			return
		} else if !match {
			c.AbortWithStatus(http.StatusUnauthorized)
			logHelper.LogEvent("authfox", "loginUser(): Invalid password received")
			return
		}

		// create session
		session, err := sessionHelper.CreateSession(&localUserID, redisVerify, redisSession, verify)
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
func checkLoginData(loginData *sendLogin) bool {
	if loginData.LoginName == "" {
		return false
	}
	if loginData.Password == "" {
		return false
	}
	return true
}

func findUserData(pg_conn *gorm.DB, login *string) (password, UserID string, verify bool, err error) {
	// we'll send verify as true on failture as they are limited to a single use case

	var localProfile Profile
	var res *gorm.DB
	// switch wether is an email or account name
	if stringHelper.CheckEmail(strings.ToLower(*login)) {
		// get account by email
		res = pg_conn.Where("email = ?", strings.ToLower(*login)).Take(&localProfile)
	} else {
		// get account by user name
		res = pg_conn.Where("name_static = ?", strings.ToLower(*login)).Take(&localProfile)
	}

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return "", "", true, res.Error
	} else if res.RowsAffected > 1 {
		// illegal value, two accounts???
		return "", "", true, errors.New("findUserData(): DB returned multiple accounts for a single search in user table")
	} else if res.RowsAffected == 0 {
		// no account was found
		// searching for one in the verify table
		var localVerify Verify
		// switch search method
		if stringHelper.CheckEmail(strings.ToLower(*login)) {
			// get account by email
			res = pg_conn.Where("email = ?", strings.ToLower(*login)).Take(&localVerify)
		} else {
			// get account by user name
			res = pg_conn.Where("name_static = ?", strings.ToLower(*login)).Take(&localVerify)
		}

		if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
			return "", "", true, res.Error
		} else if res.RowsAffected > 1 {
			return "", "", true, errors.New("findUserData(): DB returned multiple accounts for a single search in verify table")
		} else if res.RowsAffected == 0 {
			// no account was found
			// returning special error
			return "", "", true, ErrAccountNotExisting
		} else if res.RowsAffected == 1 {
			// account exists! Everything is good
			return localVerify.Password, localVerify.UserID, true, nil
		} else {
			// illegal edge case
			return "", "", true, errors.New("findUserData(): entered illegal edge case on verify DB checking")
		}
	} else if res.RowsAffected == 1 {
		// account exists! Everything is good

		// fetch password
		var localUser User
		res := pg_conn.Where("user_id = ?", localProfile.UserID).Take(&localUser)
		if res.Error != nil {
			return "", "", true, res.Error
		} else if res.RowsAffected != 1 {
			return "", "", true, errors.New("findUserData(): Invalid number of rows found on getting password from user DB")
		}
		// return the found password
		return localUser.Password, localUser.UserID, false, nil
	} else {
		// illegal edge case
		return "", "", true, errors.New("findUserData(): entered illegal edge case on user DB checking")
	}
}
