package endpoints

import (
	"github.com/PurotoApp/libpuroto/ginHelper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetRoutes(router *gin.Engine, collUsers, collVerify, collSession, collVerifySession, collProfiles *mongo.Collection) {
	router.POST("/v1/user", registerUser(collUsers, collVerify, collSession, collVerifySession, collProfiles))
	router.POST("/v1/user/login", loginUser(collUsers, collSession, collVerifySession, collVerify, collProfiles))
	router.POST("/v1/user/verify", verifyUser(collVerifySession, collSession, collVerify, collProfiles, collUsers))
	router.POST("/v1/user/validate", validateSession(collVerifySession, collSession))
	router.PATCH("/v1/user", updatePassword(collVerifySession, collSession, collUsers))
	router.POST("/v1/user/delete", accountDeletion(collVerifySession, collSession, collUsers, collProfiles))
	// swagger docs
	router.Static("/swagger", "swagger/")
	// user redirects
	router.GET("/", ginHelper.Redirect("/swagger"))
	router.GET("/v1", ginHelper.Redirect("/swagger"))
}
