package sessionHelper

import (
	"context"
	"errors"
	"time"

	"github.com/PurotoApp/authfox/internal/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// generates a secure session token
func generateSessionToken(collSession, collVerifySession *mongo.Collection) (string, error) {
	var token string
	var err error
	var count int64

	for i := 0; i < 20; {
		token, err = security.RandomString(512)
		if err != nil {
			return "", err
		}
		// check if the token already exists
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		count, err = collSession.CountDocuments(ctx, bson.D{{Key: "token", Value: token}}, options.Count().SetLimit(1))
		cancel()
		if err != nil {
			return "", err
		}
		// doesn't exist, so continue
		if count == 0 {
			// check if it exists as verify session
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
			count, err = collVerifySession.CountDocuments(ctx, bson.D{{Key: "token", Value: token}}, options.Count().SetLimit(1))
			cancel()
			if count == 0 {
				return token, nil
			}
		}
	}
	return "", errors.New("generateSessionToken(): failed to generate token after 20 tries")
}
