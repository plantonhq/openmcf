package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/plantonhq/openmcf/app/backend/internal/database"
	"github.com/plantonhq/openmcf/app/backend/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	// Parse command line flags
	port := flag.String("port", getEnv("SERVER_PORT", "50051"), "Server port")
	mongoURI := flag.String("mongo-uri", getEnv("MONGODB_URI", "mongodb://localhost:27017"), "MongoDB connection URI")
	mongoDatabase := flag.String("mongo-database", getEnv("MONGODB_DATABASE", "openmcf"), "MongoDB database name")
	flag.Parse()

	// Set up logging
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	ctx := context.Background()

	// Connect to MongoDB
	logrus.WithFields(logrus.Fields{
		"uri":      *mongoURI,
		"database": *mongoDatabase,
	}).Info("Connecting to MongoDB")

	mongo, err := database.Connect(ctx, *mongoURI, *mongoDatabase)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to MongoDB")
	}
	defer func() {
		if err := mongo.Disconnect(context.Background()); err != nil {
			logrus.WithError(err).Error("Failed to disconnect from MongoDB")
		}
	}()

	// Create and start server
	cfg := &server.Config{
		Port:    *port,
		MongoDB: mongo,
	}

	srv := server.NewServer(cfg)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			logrus.WithError(err).Fatal("Server failed")
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	logrus.Info("Received shutdown signal")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logrus.WithError(err).Fatal("Failed to shutdown server gracefully")
	}

	logrus.Info("Server stopped")

	os.Exit(0)
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
