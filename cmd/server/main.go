package main

import (
	"feedback-generator/internal/config"
	"feedback-generator/pkg/cmd"
	"log"
	"os"
)

func main() {
	logger, err := config.CreateDefaultLogConfiguration()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	if err := cmd.RunServer(); err != nil {
		logger.WithError(err).Fatalf("%v\n", err)
		os.Exit(1)
	}
}
