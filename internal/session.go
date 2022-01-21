package authfox

import (
	"context"
	"time"

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
	UserID string `json:"uid"`
	Token  string `json:"token"`
}

func createSession(userID string, collSession *mongo.Collection) (sessionPair, error) {
	// genrate new token
	token, err := generateSessionToken(collSession)
	if err != nil {
		return sessionPair{}, err
	}

	// save session
	_, err = collSession.InsertOne(context.TODO(), newSession{Token: token, UserID: userID, CreationTime: time.Now()})
	if err != nil {
		return sessionPair{}, err
	}

	return sessionPair{Token: token, UserID: userID}, nil
}
