package authfox

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: move inline
func addVerifyUser(userStruct saveVerifyUser, collVerify *mongo.Collection) error {
	_, err := collVerify.InsertOne(context.TODO(), userStruct)
	return err
}
