package authfox

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetRoutes(router *gin.Engine, collUsers, collVerify, collSession, collVerifySession, collProfiles *mongo.Collection) {
	router.POST("/v1/register", registerUser(collUsers, collVerify, collSession, collVerifySession))
	router.POST("/v1/login", loginUser(collUsers, collSession, collVerifySession, collVerify))
	router.POST("/v1/verify", verifyUser(collVerifySession, collSession, collVerify, collProfiles, collUsers))
}
