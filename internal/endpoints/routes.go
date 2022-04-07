package endpoints

import (
	"github.com/PurotoApp/libpuroto/ginHelper"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func SetRoutes(router *gin.Engine, pg_conn *gorm.DB, redisVerify, redisSession *redis.Client) {
	router.POST("/v1/user", registerUser(pg_conn, redisVerify, redisSession))
	router.POST("/v1/user/login", loginUser(pg_conn, redisVerify, redisSession))
	router.POST("/v1/user/verify", verifyUser(pg_conn, redisVerify, redisSession))
	// router.POST("/v1/user/validate", validateSession(collVerifySession, collSession))
	// router.PATCH("/v1/user", updatePassword(collVerifySession, collSession, collUsers))
	// router.POST("/v1/user/delete", accountDeletion(collVerifySession, collSession, collUsers, collProfiles))
	// swagger docs
	router.Static("/swagger", "swagger/")
	// user redirects
	router.GET("/", ginHelper.Redirect("/swagger"))
	router.GET("/v1", ginHelper.Redirect("/swagger"))
}
