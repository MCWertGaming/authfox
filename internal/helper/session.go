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
	"strconv"
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
		var sessionNumber uint8

		// find the first session that can be used
		for sessionNumber = 0; sessionNumber < 6; sessionNumber++ {
			// check if a session with that number already exists
			if count, err := redisSession.Exists(*userID + strconv.Itoa(int(sessionNumber))).Result(); count == 0 {
				break
			} else if err != nil {
				return sessionPair{}, err
			}
		}
		// removing the next slot to keep the session count at 5
		// the 5th slot removes the first
		if sessionNumber == 5 {
			if err := redisSession.Del(*userID + "0").Err(); err != nil {
				return sessionPair{}, err
			}
		} else {
			// remove the next slot
			if err := redisSession.Del(*userID + strconv.Itoa(int(sessionNumber+1))).Err(); err != nil {
				return sessionPair{}, err
			}
		}
		// creating one using this ID
		if err := redisSession.Set(*userID+strconv.Itoa(int(sessionNumber)), token, time.Hour*24*7).Err(); err != nil {
			return sessionPair{}, err
		}
		return sessionPair{Token: token, UserID: *userID + strconv.Itoa(int(sessionNumber)), VerifyOnly: verify}, nil
	}
}
