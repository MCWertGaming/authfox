package authfox

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetRoutes(router *gin.Engine, client *mongo.Client) {
	router.POST("/v1/register", registerUser(client))
}
