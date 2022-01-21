package authfox

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetRoutes(router *gin.Engine, collUsers *mongo.Collection, collVerify *mongo.Collection, collSession *mongo.Collection) {
	router.POST("/v1/register", registerUser(collUsers, collVerify, collSession))
}
