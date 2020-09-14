package grpc

import (
	"context"
	"feedback-generator/internal/config"
	f "feedback-generator/pkg/api/v1/feedbackreqpb"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

// RunServer runs gRPC service to publish Feedback service
func RunServer(ctx context.Context, v1fService f.FeedbackServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	logger := config.GetDefaultLogger()
	logger.Info("Registering FeedbackService to gRPC Server")
	// register service
	server := grpc.NewServer()
	f.RegisterFeedbackServiceServer(server, v1fService)
	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Info("shutting down gRPC server...")
			//log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Info("Starting gRPC server...")
	return server.Serve(listen)
}
