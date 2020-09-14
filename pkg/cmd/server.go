package cmd

import (
	"context"
	"feedback-generator/internal/config"
	"feedback-generator/pkg/protocol/grpc"
	v1 "feedback-generator/pkg/service/v1"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RunServer will start grpc Server
func RunServer() error {
	//Read server configuration
	c, err := config.ReadConfig()
	logger := config.GetDefaultLogger()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"filename": "config",
			"status":   500,
			"Error":    err,
		}).Fatal("Unable to read the Config file given!")
	}

	fmt.Println(c)

	//Connect to MongoDB database
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(c.DBUri)
	clientOptions.SetMaxPoolSize(12)

	client, err := mongo.Connect(ctx, clientOptions)
	db := client.Database(c.DBName)
	defer client.Disconnect(ctx)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"database": c.DBName,
			"host":     c.DBHost,
			"Error":    err,
			"status":   500,
		}).Errorf("Unable to connect to MongoDB database: %s and host: %s", c.DBName, c.DBHost)

		return err
	}

	//Initialize service with mongoClient
	v1fService := v1.NewFeedbackServiceServer(client, db, logger)

	return grpc.RunServer(ctx, v1fService, c.GRPCPort)
}
