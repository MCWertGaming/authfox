package endpoints

import (
	"github.com/PurotoApp/authfox/internal/ginHelper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetRoutes(router *gin.Engine, collUsers, collVerify, collSession, collVerifySession, collProfiles *mongo.Collection) {
	router.POST("/v1/register", registerUser(collUsers, collVerify, collSession, collVerifySession, collProfiles))
	router.POST("/v1/login", loginUser(collUsers, collSession, collVerifySession, collVerify, collProfiles))
	router.POST("/v1/verify", verifyUser(collVerifySession, collSession, collVerify, collProfiles, collUsers))
	router.POST("/v1/validate", validateSession(collVerifySession, collSession))
	router.POST("/v1/update", updatePassword(collVerifySession, collSession, collUsers))
	router.POST("/v1/remove", accountDeletion(collVerifySession, collSession, collUsers, collProfiles))
	// swagger docs
	router.Static("/swagger", "swagger/")
	// user redirects
	router.GET("/", ginHelper.Redirect("/swagger"))
	router.GET("/v1", ginHelper.Redirect("/swagger"))
}
