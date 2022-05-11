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
	"net/http"

	"github.com/PurotoApp/libpuroto/logHelper"
	"github.com/gin-gonic/gin"
)

// returns true, if the client requested json format, also sets the response to 406, if not
func JsonRequested(c *gin.Context) bool {
	if c.GetHeader("Content-Type") != "application/json" {
		c.AbortWithStatus(http.StatusNotAcceptable)
		logHelper.LogEvent("authfox", "Received request with wrong Content-Type header")
		return false
	}
	return true
}
