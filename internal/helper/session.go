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
	"time"

	"github.com/PurotoApp/libpuroto/libpuroto"
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
	token, err := libpuroto.RandomString(512)
	if err != nil {
		return sessionPair{}, err
	}
	// select session type
	if verify {
		// creating a verify session, only one is allowed
		// so we'll just create a new secret and store it into redis
		// verify session is valid for 2 days
		// this will override the old session if neccessary
		if err = redisVerify.Set(*userID, token, time.Hour*48).Err(); err != nil {
			return sessionPair{}, err
		}
		return sessionPair{UserID: *userID, Token: token, VerifyOnly: verify}, nil
	} else {
		// creating a user session, only 5 are allowed
		// sessions are valid for 7 days
		// because redis can only store one key, we'll append a number to the UID
		// UID[session_number] : token
		if count, err := redisSession.Exists(*userID + "0").Result(); count == 0 {
			// no sessions, removing the next slot to keep the session count at 5
			redisSession.Del(*userID + "1")
			// creating one using this ID
			redisSession.Set(*userID+"0", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "0", VerifyOnly: verify}, nil
		} else if err != nil {
			return sessionPair{}, err
		} else if count, err := redisSession.Exists(*userID + "1").Result(); count == 0 {
			// no sessions, removing the next slot to keep the session count at 5
			redisSession.Del(*userID + "2")
			// creating one using this ID
			redisSession.Set(*userID+"1", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "1", VerifyOnly: verify}, nil
		} else if err != nil {
			return sessionPair{}, err
		} else if count, err := redisSession.Exists(*userID + "2").Result(); count == 0 {
			// no sessions, removing the next slot to keep the session count at 5
			redisSession.Del(*userID + "3")
			// creating one using this ID
			redisSession.Set(*userID+"2", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "2", VerifyOnly: verify}, nil
		} else if err != nil {
			return sessionPair{}, err
		} else if count, err := redisSession.Exists(*userID + "3").Result(); count == 0 {
			// no sessions, removing the next slot to keep the session count at 5
			redisSession.Del(*userID + "4")
			// creating one using this ID
			redisSession.Set(*userID+"3", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "3", VerifyOnly: verify}, nil
		} else if err != nil {
			return sessionPair{}, err
		} else if count, err := redisSession.Exists(*userID + "4").Result(); count == 0 {
			// no sessions, removing the next slot to keep the session count at 5
			redisSession.Del(*userID + "5")
			// creating one using this ID
			redisSession.Set(*userID+"4", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "4", VerifyOnly: verify}, nil
		} else if err != nil {
			return sessionPair{}, err
		} else {
			// no sessions, removing the next slot to keep the session count at 5
			redisSession.Del(*userID + "0")
			// create a 6th session because the first one is made free again
			redisSession.Set(*userID+"5", token, time.Hour*24*7)
			return sessionPair{Token: token, UserID: *userID + "5", VerifyOnly: verify}, nil
		}
	}
}
