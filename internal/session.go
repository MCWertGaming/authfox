package authfox

import (
	"context"
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
