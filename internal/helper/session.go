/* <AuthFox - a simple authentication and session server for Puroto>
   Copyright (C) 2022  PurotoApp

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package helper

import (
	"crypto/subtle"
	"time"

	"github.com/go-redis/redis"
)

// session information for sending to the client
type sessionPair struct {
	UserID     string `json:"uid"`
	Token      string `json:"token"`
	VerifyOnly bool   `json:"verify_only"`
}

// TODO: use string pointer for UID
func CreateSession(userID *string, redisVerify, redisSession *redis.Client, verify bool) (sessionPair, error) {
	// session token
	token, err := RandomString(512)
	if err != nil {
		return sessionPair{}, err
	}
	// select session type
	if verify {
		// creating a verify session, only one is allowed
		// so we'll just create a new secret and store it into redis
		// verify session is valid for 2 days
		// this will override the old session if neccessary
		if redisVerify.Set(*userID, token, time.Hour*48).Err() != nil {
			return sessionPair{}, err
		}
		return sessionPair{UserID: *userID, Token: token, VerifyOnly: verify}, nil
	} else {
		// creating a user session, only 4 are allowed
		// sessions are valid for 2 days
		// because redis can only store one key, we'll append a number to the UID
		// UID[session_number] : token
		if count, err := redisSession.Exists(*userID + "0").Result(); count == 0 {
			// no sessions, creating one using this ID
			redisSession.Set(*userID+"0", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "0", VerifyOnly: verify}, nil
		} else if err != nil {
			return sessionPair{}, err
		} else if count, err := redisSession.Exists(*userID + "1").Result(); count == 0 {
			// no sessions, creating one using this ID
			redisSession.Set(*userID+"1", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "1", VerifyOnly: verify}, nil
		} else if err != nil {
			return sessionPair{}, err
		} else if count, err := redisSession.Exists(*userID + "2").Result(); count == 0 {
			// no sessions, creating one using this ID
			redisSession.Set(*userID+"2", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "2", VerifyOnly: verify}, nil
		} else if err != nil {
			return sessionPair{}, err
		} else if count, err := redisSession.Exists(*userID + "3").Result(); count == 0 {
			// no sessions, creating one using this ID
			redisSession.Set(*userID+"3", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "3", VerifyOnly: verify}, nil
		} else if err != nil {
			return sessionPair{}, err
		} else if count, err := redisSession.Exists(*userID + "4").Result(); count == 0 {
			// no sessions, creating one using this ID
			redisSession.Set(*userID+"4", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "4", VerifyOnly: verify}, nil
		} else if err != nil {
			return sessionPair{}, err
		} else {
			// overwrite the first session since the session limit is reached
			redisSession.Set(*userID+"0", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "0", VerifyOnly: verify}, nil
		}
	}
}

// returns true if the session of the given redis DB is valid
func SessionValid(uid, token *string, redisClient *redis.Client) (bool, error) {
	var res string
	var err error

	// the UUID session extension is part of the session, so no work is needed
	res, err = redisClient.Get(*uid).Result()

	if err != nil {
		return false, err
		// } else if res != *token {
	} else if subtle.ConstantTimeCompare([]byte(res), []byte(*token)) != 1 {
		// TODO: Use secure matching function
		// session and token don't match
		return false, nil
	}
	// the session seems valid
	return true, nil
}
