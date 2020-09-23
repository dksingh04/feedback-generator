package main

import (
	"feedback-generator/internal/config"
	"feedback-generator/pkg/cmd"
	"flag"
	"log"
	"os"
)

func main() {
	logger, err := config.CreateDefaultLogConfiguration()
	var db string
	//Added flag to control for connecting to database or not, using MongoDB which connects to local db at port 27017
	flag.StringVar(&db, "db", "ignore", "Usage: -db=connect to connect to db")
	flag.Parse()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	var ignoreDB bool
	if db == "ignore" {
		ignoreDB = true
	} else {
		ignoreDB = false
	}
	//db := flag.String("word", "foo", "a string")
	if err := cmd.RunServer(ignoreDB); err != nil {
		logger.WithError(err).Fatalf("%v\n", err)
		os.Exit(1)
	}
}
