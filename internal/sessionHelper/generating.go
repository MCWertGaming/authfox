package sessionHelper

import (
	"context"
	"errors"

	"github.com/PurotoApp/authfox/internal/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// generates a secure user ID
func GenerateUserID(collUsers *mongo.Collection, collVerify *mongo.Collection) (string, error) {
	var userID string
	var count_users, count_verify int64
	var err error

	for i := 0; i < 20; {
		userID, err = security.RandomString(100)
		if err != nil {
			return "", err
		}
		// check if the token already exists
		// TODO: Set timeout to 100ms
		count_users, err = collUsers.CountDocuments(context.TODO(), bson.D{{Key: "token", Value: userID}}, options.Count().SetLimit(1))
		if err != nil {
			return "", err
		}
		// TODO: set timeout to 100ms
		count_verify, err = collVerify.CountDocuments(context.TODO(), bson.D{{Key: "token", Value: userID}}, options.Count().SetLimit(1))
		if err != nil {
			return "", err
		}
		// doesn't exist, so continue
		if count_users == 0 || count_verify == 0 {
			return userID, nil
		}
	}
	return "", errors.New("generateUserID(): Failed to create token after 20 tries")
}

// generates a secure session token
func generateSessionToken(collSession *mongo.Collection) (string, error) {
	var token string
	var err error
	var count int64

	for i := 0; i < 20; {
		token, err = security.RandomString(512)
		if err != nil {
			return "", err
		}
		// check if the token already exists
		// TODO: Set timeout to 100ms
		count, err = collSession.CountDocuments(context.TODO(), bson.D{{Key: "token", Value: token}}, options.Count().SetLimit(1))
		// TODO: Also search the verify session
		if err != nil {
			return "", err
		}
		// doesn't exist, so continue
		if count == 0 {
			return token, nil
		}
	}
	return "", errors.New("generateSessionToken(): failed to generate token after 20 tries")
}
