package helper

import (
	"github.com/PurotoApp/libpuroto/logHelper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	// build connection URI
	// TODO: remove
	uri := "host=localhost user=user password=pass dbname=authfox port=5432 sslmode=disable TimeZone=Europe/Berlin"
	//uri := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
	//	os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"),
	//	os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_SSLMODE"),
	//	os.Getenv("POSTGRES_TIMEZONE"))

	// connect to the DB
	// TODO: does 'warn' loglevel contain sensitive data?
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		logHelper.ErrorFatal("authfox", err)
	}
	return db
}
