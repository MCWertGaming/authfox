package sessionHelper

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	VerifyOnly bool   `json:"verify_only,omitempty"`
}

func CreateSession(userID string, collSession, collVerifySession *mongo.Collection, verify bool) (sessionPair, error) {
	// genrate new token
	token, err := generateSessionToken(collSession, collVerifySession)
	if err != nil {
		return sessionPair{}, err
	}

	// save session
	if verify {
		// check and remove session
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		_, err := collVerifySession.DeleteOne(ctx, bson.M{"uid": userID})
		cancel()
		if err != nil {
			return sessionPair{}, err
		}

		// add session to the verify DB
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
		_, err = collVerifySession.InsertOne(ctx, newSession{Token: token, UserID: userID, CreationTime: time.Now()})
		cancel()
		if err != nil {
			return sessionPair{}, err
		}
	} else {
		// check how many sessions are open
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		count, err := collSession.CountDocuments(ctx, bson.M{"uid": userID}, options.Count().SetLimit(5))
		cancel()
		if err != nil {
			return sessionPair{}, err
		}
		if count > 4 {
			// the user has 5 or more sessions, let's remove one
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
			_, err := collSession.DeleteOne(ctx, bson.M{"uid": userID})
			cancel()
			if err != nil {
				return sessionPair{}, err
			}
		}

		// add session to the session DB
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
		_, err = collSession.InsertOne(ctx, newSession{Token: token, UserID: userID, CreationTime: time.Now()})
		cancel()
		if err != nil {
			return sessionPair{}, err
		}
	}
	return sessionPair{Token: token, UserID: userID, VerifyOnly: verify}, nil
}

// returns true if the session is valid
func SessionValid(uid, token *string, collVerifySession, collSession *mongo.Collection, verify bool) (bool, error) {
	var sessionDataRaw *mongo.SingleResult

	// search for the session
	if verify {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		sessionDataRaw = collVerifySession.FindOne(ctx, bson.D{{Key: "uid", Value: uid}, {Key: "token", Value: token}})
		cancel()
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		sessionDataRaw = collSession.FindOne(ctx, bson.D{{Key: "uid", Value: uid}, {Key: "token", Value: token}})
		cancel()
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
