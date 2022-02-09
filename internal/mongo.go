package authfox

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: move inline
func addVerifyUser(userStruct saveVerifyUser, collVerify *mongo.Collection) error {
	ctx,cancel := context.WithTimeout(context.Background(), time.Millisecond * 20)
	defer cancel()
	_, err := collVerify.InsertOne(ctx, userStruct)
	return err
}
