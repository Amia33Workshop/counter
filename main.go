package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	if err := godotenv.Load(); err != nil {
		LogInfo("No .env file found")
	}
	InitLogger()
	if os.Getenv("LOG_LEVEL") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	loadThemes()
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI).SetConnectTimeout(10 * time.Second)
	var err error
	client, err = mongo.Connect(opts)
	if err != nil {
		LogError(err)
		return
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			LogError(err)
		}
	}()
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		LogError(err)
		return
	}
	LogInfo("Pinged your deployment. You successfully connected to MongoDB!")
	collection = client.Database("counter").Collection("counters")
	dbIntervalStr := os.Getenv("DB_INTERVAL")
	dbInterval := 60
	if dbIntervalStr != "" {
		if val, err := strconv.Atoi(dbIntervalStr); err == nil {
			dbInterval = val
		}
	}
	if dbInterval > 0 {
		ticker := time.NewTicker(time.Duration(dbInterval) * time.Second)
		go func() {
			for range ticker.C {
				pushCacheToDB()
			}
		}()
	}
	router := setupRouter()
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			LogErrorf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	LogInfo("Shutting down server...")
	LogInfo("Flushing cache to database...")
	pushCacheToDB()
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		LogError("Server forced to shutdown:", err)
	}
	LogInfo("Server exiting")
}
