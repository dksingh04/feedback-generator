package main

import (
	"context"
	"feedback-generator/internal/config"
	f "feedback-generator/pkg/api/v1/feedbackreqpb"
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

var serverAddr = "localhost:9090"

type clientConfig struct {
	logger     *logrus.Logger
	clientConn *grpc.ClientConn
	client     f.FeedbackServiceClient
}

func main() {
	logger, err := config.CreateDefaultLogConfiguration()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"filename": "logger",
			"status":   500,
			"Error":    err,
		}).Fatal("Unable to read the Config file given!")
	}

	//Connect to grpc server
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"server": "localhost:9090",
			"status": 500,
			"Error":  err,
		}).Fatalln("Unable to connect to grpc server")
	}

	client := f.NewFeedbackServiceClient(conn)

	cc := &clientConfig{
		logger:     logger,
		clientConn: conn,
		client:     client,
	}

	app := cli.NewApp()

	app.EnableBashCompletion = true
	app.Name = "feedback-generator"
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "Deepak Singh",
			Email: "deepaksingh04@gmail.com",
		},
	}
	app.Copyright = "(c) 2020 feedback-generator"
	app.Commands = []*cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create Feedback Request",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "-c [requestID]",
				},
			},
			Action: func(c *cli.Context) error {
				cc.createFeedbackRequest()
				return nil
			},
		},
		{
			Name:    "read",
			Aliases: []string{"r"},
			Usage:   "Read Feedback Request",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "read",
					Aliases: []string{"r"},
					Usage:   "-r [requestID]",
				},
			},
			Action: func(c *cli.Context) error {
				requestID := c.String("read")
				cc.readFeedbackRequest(requestID)
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete created request",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "-d [requestID]",
				},
			},
			Action: func(c *cli.Context) error {
				requestID := c.String("delete")
				cc.deleteFeedbackRequest(requestID)
				return nil
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "gen",
					Aliases: []string{"g"},
					Usage:   "-g [requestID]",
				},
			},
			Usage: "Generate feedback for the created request",
			Action: func(c *cli.Context) error {
				requestID := c.String("gen")
				cc.generateFeedbackResponse(requestID)
				return nil
			},
		},
	}

	errs := app.Run(os.Args)

	if errs != nil {
		logger.Fatalf("Error in initiating commands %s", errs)
	}

}

func (cc *clientConfig) createFeedbackRequest() {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	t := time.Now().In(time.UTC)
	createDate, _ := ptypes.TimestampProto(t)

	fReq := f.FeedbackRequest{
		Api: apiVersion,
		FeedbackReq: &f.Feedback{
			CandidateName:           "Deepak Singh",
			RecruiterName:           "Amanda",
			CreatDate:               createDate,
			UpdateDate:              createDate,
			IsCodingRequired:        true,
			IsAbleToWritePseudoCode: true,
			IsAlgoEfficient:         true,
			IsCodeCompiled:          false,
			IsIdRequired:            false,
			IsProxy:                 false,
			IsWhiteboardingRequired: false,
			IsWhiteboardDone:        false,
			JobType:                 "Full-Stack Java Developer",
			MyComments:              "Good Candidate",
			TechSkills: []*f.TechSkill{
				&f.TechSkill{
					SkillName:            "Java",
					ExperienceRating:     3,
					SkillRating:          3,
					IsHandsOn:            true,
					InDepthUnderstanding: true,
					QuestionsAsked:       "OOPs, Mutithreading, Exception Handling, Collection classes etc.",
					Topics: []*f.Topic{
						&f.Topic{
							TopicName:                 "OOPs",
							HaveTheroreticalKnowledge: true,
							InDepthUnderstanding:      true,
							TheoryQuestion:            "Encapsulation, Inheritance, Polymorphism etc.",
							IsAbleToExaplain:          true,
							IsAbleToExplainScenario:   true,
							IsScenarioCovered:         true,
							IsHandsOn:                 true,
							PartiallyExplained:        false,
							WhatSceanrioQuestion:      "How OOPs concept used in his project, explain with example",
						},
					},
				},
			},
		},
	}

	//fRes, err := client.CreateSimpleRequest(ctx, &sReq)
	fRes, err := cc.client.Create(ctx, &fReq)
	fmt.Println(fRes)
	//TODO this can be moved to common method
	if err != nil {
		marshaler := jsonpb.Marshaler{}
		jsonReq, errMarshaling := marshaler.MarshalToString(&fReq)
		if errMarshaling != nil {
			cc.logger.WithFields(logrus.Fields{
				"reqMessage": fReq.FeedbackReq,
				"status":     500,
				"Error":      err,
			}).Fatalln("Error in Transforming Request message to Json String!")
		}
		cc.logger.WithFields(logrus.Fields{
			"request": jsonReq,
			"status":  500,
			"Error":   err,
		}).Fatalln("Unable to connect to grpc server")
	}
}

func (cc *clientConfig) readFeedbackRequest(requestID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	rf := &f.ReadFeedbackRequest{
		Api:       "v1",
		RequestId: requestID,
	}
	fRes, err := cc.client.Read(ctx, rf)
	fmt.Println(fRes)
	fmt.Println(err)
}

func (cc *clientConfig) deleteFeedbackRequest(requestID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	drf := &f.DeleteFeedbackRequest{
		Api:       "v1",
		RequestId: requestID,
	}
	dfRes, err := cc.client.Delete(ctx, drf)

	fmt.Println(dfRes)
	fmt.Println(err)
}

func (cc *clientConfig) generateFeedbackResponse(requestID string) {
	fmt.Println(requestID)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	feedback := &f.Feedback{Id: requestID}
	fReq := &f.GenerateFeedbackRequest{Api: "v1", FeedbackReq: feedback, SummaryNote: "Summary Notes"}
	fRes, err := cc.client.GenerateFeedbackForRequest(ctx, fReq)
	fmt.Println("Professional Summary:")
	fmt.Println("")
	fmt.Printf(fRes.SummaryText)
	fmt.Println("")
	fmt.Printf("%s:\n\n", fRes.SkillFeedback[0].Skill)
	fmt.Printf(fRes.SkillFeedback[0].FeedbackText)
	fmt.Println("")
	fmt.Println("Error:")
	fmt.Println(err)
}
