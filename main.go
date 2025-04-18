package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	httpserver "github.com/baohuamap/zchat-api/api/http"
	"github.com/baohuamap/zchat-api/api/ws"
	"github.com/baohuamap/zchat-api/pkg/aws"
	"github.com/baohuamap/zchat-api/pkg/gorm"
	"github.com/baohuamap/zchat-api/repository"
	"github.com/baohuamap/zchat-api/router"
	"github.com/baohuamap/zchat-api/service"
)

func init() {
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(h))
}

func main() {
	slog.Info("Starting Chat Server")

	if err := godotenv.Overload(); err != nil {
		slog.Error("Cannot load .env file: ", "error", err.Error())
	}

	db := gorm.NewDB()

	dsn := strings.Join([]string{
		"host=", os.Getenv("PG_HOST"),
		" port=", os.Getenv("PG_PORT"),
		" user=", os.Getenv("PG_USER"),
		" password=", os.Getenv("PG_PASSWORD"),
		" database=", os.Getenv("PG_NAME"),
		" sslmode=", os.Getenv("PG_SSLMODE"),
	}, "")
	if err := db.Connect(dsn); err != nil {
		slog.Error("Creating connection to DB: ", slog.String("error", err.Error()))
	}

	r := gin.Default()

	server := &http.Server{
		Addr:         "0.0.0.0:5050",
		Handler:      r.Handler(),
		ReadTimeout:  time.Duration(20) * time.Second,
		WriteTimeout: time.Duration(20) * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(20)*time.Second)
	defer cancel()

	userRepo := repository.NewUserRepository(db.Gormer())
	friendshipRepo := repository.NewFriendshipRepository(db.Gormer())
	participantRepo := repository.NewParticipantRepository(db.Gormer())
	conversationRepo := repository.NewConversationRepository(db.Gormer())
	messageRepo := repository.NewMessageRepository(db.Gormer())

	s3Client, err := aws.NewS3Client(ctx)
	if err != nil {
		slog.Error("Creating S3 client: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	u := service.NewUserService(userRepo, friendshipRepo, s3Client)
	m := service.NewMessageService(conversationRepo, messageRepo, participantRepo)
	httpHandler := httpserver.NewHandler(u, m)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub, conversationRepo, participantRepo, messageRepo)
	go hub.Run()

	router.SetupRoutes(r, httpHandler, wsHandler)

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to listen")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server...")

	if err := server.Shutdown(ctx); err != nil {
		slog.Info("Failed to shutdown server: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Server exiting...")
}
