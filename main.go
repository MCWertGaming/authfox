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

package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/PurotoApp/libpuroto/libpuroto"
)

// test

type TestDB1 struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	Test      string
}

func createID(id uint, create_time time.Time) string {
	var bits uint64 = uint64(create_time.Unix() << 16)
	var bits_new uint64 = uint64(bits | uint64(id))
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(bits_new))
	return base64.RawStdEncoding.EncodeToString(b)
}
func getRowNum(ID string, create_time time.Time) uint {
	decoded, _ := base64.RawStdEncoding.DecodeString(ID)
	return uint(binary.LittleEndian.Uint64(decoded) & uint64(create_time.Unix()))
}

func main() {
	// connect to the PostgreSQL
	pg_conn := libpuroto.ConnectDB()

	if err := pg_conn.AutoMigrate(TestDB1{}); err != nil {
		libpuroto.ErrorPanic(err)
	}
	var testRow TestDB1
	var created_id string
	for i := 0; i < 9000; i++ {
		testRow = TestDB1{Test: "Hello"}
		pg_conn.Create(&testRow)
		created_id = createID(testRow.ID, testRow.CreatedAt)
		fmt.Println(created_id)
		if testRow.ID != getRowNum(created_id, testRow.CreatedAt) {
			fmt.Println("Error!")
		}
	}
}

/*
type TestDB struct {
	ID        string
	CreatedAt time.Time
	Test      string
}

func main() {
	// connect to the PostgreSQL
	pg_conn := libpuroto.ConnectDB()

	if err := pg_conn.AutoMigrate(TestDB{}); err != nil {
		libpuroto.ErrorPanic(err)
	}
	var testRow TestDB

	for i := 0; i < 9000; i++ {
		testRow = TestDB{Test: "Hello", ID: string(uuid.New().String())}
		pg_conn.Create(&testRow)
		fmt.Println(testRow.ID)
	}
}
*/
// ------------------------
/*
func main() {
	// connect to the PostgreSQL
	pg_conn := libpuroto.ConnectDB()
	// Connect to Redis
	redisVerify := libpuroto.Connect(1)
	redisSession := libpuroto.Connect(2)

	// migrate all tables
	endpoints.AutoMigrateAuthfox(pg_conn)

	// create router
	router := gin.Default()

	// configure gin
	libpuroto.ConfigRouter(router)

	// set routes
	endpoints.SetRoutes(router, pg_conn, redisVerify, redisSession)

	// start
	router.Run("0.0.0.0:3621")
}
*/
