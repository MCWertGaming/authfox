/*
<AuthFox - a simple authentication and session server for Puroto>
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

package helper

import (
	"fmt"
	"os"

	"github.com/PurotoApp/libpuroto/logHelper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	// build connection URI
	// TODO: remove
	// uri := "host=localhost user=user password=pass dbname=authfox port=5432 sslmode=disable TimeZone=Europe/Berlin"
	uri := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"),
		os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_SSLMODE"),
		os.Getenv("POSTGRES_TIMEZONE"))

	// connect to the DB
	// TODO: does 'warn' loglevel contain sensitive data?
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		logHelper.ErrorFatal("authfox", err)
	}
	return db
}
