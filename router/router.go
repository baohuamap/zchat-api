package router

import (
	"github.com/baohuamap/zchat-api/api/http"
	"github.com/baohuamap/zchat-api/api/ws"
	"github.com/baohuamap/zchat-api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, httpHandler http.Handler, wsHandler ws.Handler) {

	r.Use(middleware.CORSMiddleware())

	// health check
	r.GET("/healthz", httpHandler.ServerStatus)

	// http
	r.POST("/signup", httpHandler.CreateUser)
	r.POST("/login", httpHandler.Login)
	r.GET("/logout", httpHandler.Logout)

	r.POST("/addFriend/:userId/:friendId", httpHandler.AddFriend)
	r.PUT("/acceptFriend/:userId/:friendId", httpHandler.AcceptFriend)
	r.PUT("/rejectFriend/:userId/:friendId", httpHandler.RejectFriend)
	r.GET("/getFriendRequests/:userId", httpHandler.GetFriendRequests)
	r.GET("/getFriends/:userId", httpHandler.GetFriends)

	// ws
	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/getClients/:roomId", wsHandler.GetClients)
}
