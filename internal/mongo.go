package authfox

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func addVerifyUser(userStruct saveVerifyUser, collVerify *mongo.Collection) error {
	_, err := collVerify.InsertOne(context.TODO(), userStruct)
	return err
}
