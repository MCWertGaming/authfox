package endpoints

import (
	"time"

	"github.com/PurotoApp/libpuroto/logHelper"
	"gorm.io/gorm"
)

// verify DB struct
type Verify struct {
	UserID       string    `gorm:"unique;not null;primaryKey"`
	NameFormat   string    `gorm:"unique;not null"`
	NameStatic   string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	Password     string    `gorm:"unique;not null"`
	RegisterIP   string    `gorm:"not null"`
	RegisterTime time.Time `gorm:"not null"`
	VerifyCode   string    `gorm:"unique;not null"`
}

// user DB struct
type User struct {
	UserID       string    `gorm:"unique;not null;primaryKey"`
	Password     string    `gorm:"unique;not null"`
	RegisterIP   string    `gorm:"not null"`
	RegisterTime time.Time `gorm:"not null"`
}

// profile DB struct, will be moved to Meltdown
type Profile struct {
	UserID     string `gorm:"unique;not null;primaryKey"`
	NameFormat string `gorm:"unique;not null"`
	NameStatic string `gorm:"unique;not null"`
	NamePretty string `gorm:"unique;not null"`
	Email      string `gorm:"unique;not null"`
	// visual
	BadgeBetaTester  bool
	BadgeAlphaTester bool
	BadgeStaff       bool
}

func AutoMigrateAuthfox(pg_conn *gorm.DB) {
	if err := pg_conn.AutoMigrate(&Verify{}, &User{}, &Profile{}); err != nil {
		logHelper.ErrorPanic(err)
	}
}
