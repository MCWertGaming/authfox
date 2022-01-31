package authfox

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// session information used for create a new session
type newSession struct {
	UserID       string    `bson:"uid"`
	Token        string    `bson:"token"`
	CreationTime time.Time `bson:"creation_time"`
}

// session information for sending to the client
type sessionPair struct {
	UserID     string `json:"uid"`
	Token      string `json:"token"`
	VerifyOnly bool   `json:"verify_only"`
}

func createSession(userID string, collSession, collVerifySession *mongo.Collection, verify bool) (sessionPair, error) {
	// genrate new token
	token, err := generateSessionToken(collSession)
	if err != nil {
		return sessionPair{}, err
	}

	// save session
	if verify {
		// check and remove session
		_, err := collVerifySession.DeleteOne(context.TODO(), bson.M{"uid": userID})
		if err != nil {
			return sessionPair{}, err
		}

		// add session to the verify DB
		_, err = collVerifySession.InsertOne(context.TODO(), newSession{Token: token, UserID: userID, CreationTime: time.Now()})
		if err != nil {
			return sessionPair{}, err
		}

	} else {
		// add session to the session DB
		// TODO: delete the last session if 5 is reached
		_, err = collSession.InsertOne(context.TODO(), newSession{Token: token, UserID: userID, CreationTime: time.Now()})
	}

	if err != nil {
		return sessionPair{}, err
	}

	return sessionPair{Token: token, UserID: userID, VerifyOnly: verify}, nil
}

// returns true if the session is valid
func sessionValid(uid, token *string, collVerifySession, collSession *mongo.Collection, verify bool) (bool, error) {
	var sessionDataRaw *mongo.SingleResult

	// search for the session
	// TODO: limit to 50ms
	if verify {
		sessionDataRaw = collVerifySession.FindOne(context.TODO(), bson.D{{Key: "uid", Value: uid}, {Key: "token", Value: token}})
	} else {
		sessionDataRaw = collSession.FindOne(context.TODO(), bson.D{{Key: "uid", Value: uid}, {Key: "token", Value: token}})
	}

	// check error
	if sessionDataRaw.Err() != nil {
		return false, sessionDataRaw.Err()
	}

	// decode DB data
	var localsessionData newSession
	if err := sessionDataRaw.Decode(&localsessionData); err != nil {
		return false, err
	}
	// check if the session is older than 7 days
	if !localsessionData.CreationTime.Add(time.Hour * 60 * 7).After(time.Now()) {
		return false, errors.New("sessionValid(): session is outdated")
	}

	return true, nil
}
