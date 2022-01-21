package mongoHelper

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetDBUri() string {
	// TODO: allow setting parameters through cli flags
	db_user_name := os.Getenv("MONGO_DB_USER")
	db_password := os.Getenv("MONGO_DB_PASSWD")
	db_host := os.Getenv("MONGO_DB_HOST")
	db_database := os.Getenv("MONGO_DB_DATABASE")
	db_authsource := os.Getenv("CASCADE_DB_AUTHSRC")

	// if no auth database is specified, the admin db will be used
	if db_authsource == "" {
		db_authsource = "admin"
	}

	return "mongodb://" + db_user_name + ":" + db_password + "@" + db_host + "/" + db_database + "?authSource=" + db_authsource
}
func ConnectDB(URI string) (*mongo.Client, error) {
	return mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
}
func TestDBConnection(client *mongo.Client) error {
	return client.Ping(context.TODO(), readpref.Primary())
}
func DisconnectDB(client *mongo.Client) error {
	return client.Disconnect(context.TODO())
}
