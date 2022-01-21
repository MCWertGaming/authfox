package ginHelper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getHealth(c *gin.Context) {
	c.Data(http.StatusOK, "application/json", []byte(`{"status":"ok"}`))
}
