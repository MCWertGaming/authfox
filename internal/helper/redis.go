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
	"os"

	"github.com/go-redis/redis"
)

func Connect(dbNumber int) *redis.Client {
	// return redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: dbNumber})
	return redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_HOST"), Password: os.Getenv("REDIS_PASS"), DB: dbNumber})
}
