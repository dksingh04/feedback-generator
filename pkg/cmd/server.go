package cmd

import (
	"context"
	"feedback-generator/internal/config"
	"feedback-generator/pkg/protocol/grpc"
	v1 "feedback-generator/pkg/service/v1"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RunServer will start grpc Server
func RunServer(ignoreDB bool) error {
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

	//Connect to MongoDB database
	var db *mongo.Database = nil
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if !ignoreDB {
		clientOptions := options.Client().ApplyURI(c.DBUri)
		clientOptions.SetMaxPoolSize(12)

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"database": c.DBName,
				"host":     c.DBHost,
				"Error":    err,
				"status":   500,
			}).Errorf("Unable to connect to MongoDB database: %s and host: %s", c.DBName, c.DBHost)

			return err
		}
		db = client.Database(c.DBName)
		defer client.Disconnect(ctx)
	}
	//Initialize service with mongoClient
	v1fService := v1.NewFeedbackServiceServer(db, logger)

	return grpc.RunServer(ctx, v1fService, c.GRPCPort)
}
