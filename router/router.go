package router

import (
	"github.com/baohuamap/zchat-api/api/http"
	"github.com/baohuamap/zchat-api/api/ws"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, httpHandler http.Handler, wsHandler ws.Handler) {

	// router.Use(middleware.ErrorHandler())

	// health check
	r.GET("/healthz", httpHandler.ServerStatus)

	// http
	r.POST("/signup", httpHandler.CreateUser)
	r.POST("/login", httpHandler.Login)
	r.GET("/logout", httpHandler.Logout)
	r.GET("/users/:userID/profile-image", httpHandler.GetProfileImage)

	// ws
	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/getClients/:roomId", wsHandler.GetClients)
}
