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

	r.GET("/user/:userId/findUsers", httpHandler.FindUsers)
	r.GET("/user/:userId", httpHandler.GetUser)
	r.POST("/user/:userId/uploadAvatar", httpHandler.UploadAvatar)

	r.POST("/addFriend/:userId/:friendId", httpHandler.AddFriend)
	r.DELETE("/cancelFriendRequest/:userId/:friendId", httpHandler.CancelFriendRequest)
	r.PUT("/acceptFriend/:friendId/:userId", httpHandler.AcceptFriend)
	r.PUT("/rejectFriend/:friendId/:userId", httpHandler.RejectFriend)
	r.GET("/sentFriendRequests/:userId", httpHandler.GetSentFriendRequests)
	r.GET("/receivedFriendRequests/:friendId", httpHandler.GetReceivedFriendRequests)
	r.GET("user/:userId/friends", httpHandler.GetFriends)

	r.GET("/user/:userId/conversations", httpHandler.LoadConversations)
	r.GET("/messages/:conversationId", httpHandler.LoadMessages)
	r.POST("/seenMessages/:conversationId", httpHandler.SeenMessages)

	// ws
	r.POST("/ws/createConversation", wsHandler.CreateConversation)
	r.GET("/ws/joinConversation/:conversationId", wsHandler.JoinConversation)
	r.GET("/ws/getClients/:conversationId", wsHandler.GetClients)
}
