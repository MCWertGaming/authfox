/* <AuthFox - a simple authentication and session server for Puroto>
   Copyright (C) 2022  PurotoApp

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package endpoints

import (
	"time"

	"github.com/PurotoApp/libpuroto/libpuroto"
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
		libpuroto.ErrorPanic(err)
	}
}
