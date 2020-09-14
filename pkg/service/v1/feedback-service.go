package v1

import (
	"context"
	v1 "feedback-generator/pkg/api/v1/feedbackreqpb"
	"fmt"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type feedbackServiceServer struct {
	client *mongo.Client
	db     *mongo.Database
	logger *logrus.Logger
}

//Initializing feedback mapping comments, this information we can have in DB, currently having it in code.
//TODO move this mapping in DB
var feedbackMapping = map[string]string{
	"CodeCompiled":  "written compilable and executable code",
	"PseudoCode":    "able to write pseudo code",
	"AlgoEfficient": "it was efficient candidate has considered space and Time complexity, while implementing the solution",
	"NotEfficient":  "it was not efficient, didn't considered or think of Time and space complexity",
	//This will go in summary and skill comments
	"Proxy":                    "using proxy and someone else was giving Interview on behalf of him, since it was proxy hence I have done some basic discussion on each skill sets.",
	"Whiteboard_explained":     "in white-boarding session, candidate has performed very well, explained the solution with proper diagram and flow.",
	"Whiteboard_partial":       "in white-boarding session, candidate was partially able to explain the solution",
	"Whiteboard_not_explained": "in white-boarding session, candidate was unable perform better, not able to solve the given problem at all",
	"Coding_Standards":         "well-versed with coding standards and followed the same while writing the code",
	"s-1":                      "Candidate needs substantial development and have to work a lot, candidate was missing fundamentals",
	"e-1":                      "Candidate has no experience in this skill and was unable to demonstrate his experience",
	"s-2":                      "Candidate needs some training to bring competency up to standards, have some basic understanding but missing some other fundamentals",
	"e-2":                      "Candidate has limited experience, close supervision will be needed for him",
	"s-3":                      "Candidate is competent and can perform his task, no additional training is required at this time",
	"e-3":                      "Candidate is competent and can complete assignments with reasonable supervision",
	"s-4":                      "Candidate is above average, and competent in this skill, no training required",
	"e-4":                      "Candidate has considerable experience and can perform his tasks with very minimal supervision",
	"s-5":                      "Candidate is expert in this skill and can teach and mentor others in the team",
	"e-5":                      "Candidate has extensive experience and can work independently",
	"HaveTheoretical":          "theoretically clear and explained the concepts of '%v' very well",
	"NoTheoretical":            "not clear with theoretical part of '%v', unable to explain '%v'",
	"InDepthUnderstanding":     "deep understanding of the technology and explained all the concepts discussed (for e.g. %s)",
	"AbleToExplain":            "theoretically clear and explained the concepts of '%v' very well with example",
	"PartiallyExplained":       "theoretically fine and partially able to explain the concepts, missing in-depth understanding of the concept, this will cause challenge in debugging/troubleshooting the problems",
	"Hands-On":                 "hands-on with the skill",
	"Hands-On-Topic":           "hands-on in this topic (%v)",
	"ScenarioQuestioned":       "I have covered some scenarion questions",
	"ScenarioExplained":        "explained the scenario question (%v) very well and how to solve the problem in such cases",
	"ScenarioNotExplained":     "unable to explain the scenario question (%v), seems to me not much hands-on in this skill",
}

var candidate = "Candidate"
var was = "was"
var has = "has"
var needs = "needs"
var is = "is"
var can = "can"
var topicsCovered = "We have covered following topics:"

// NewFeedbackServiceServer creates FeedbackServiceServer
func NewFeedbackServiceServer(client *mongo.Client, db *mongo.Database, logger *logrus.Logger) v1.FeedbackServiceServer {
	return &feedbackServiceServer{client: client, db: db, logger: logger}
}

func (fs *feedbackServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

func (fs *feedbackServiceServer) Create(ctx context.Context, req *v1.FeedbackRequest) (*v1.FeedbackResponse, error) {

	if err := fs.checkAPI(req.Api); err != nil {
		fmt.Println(err)
		return new(v1.FeedbackResponse), err
	}
	fReq := req.GetFeedbackReq()
	// Create new request Id
	fReq.Id = primitive.NewObjectID().Hex()

	coll := fs.db.Collection("feedback_request")
	//Insert created FeedbackRequest
	result, err := coll.InsertOne(ctx, fReq)

	if result.InsertedID == nil && err != nil {
		fs.logger.WithFields(logrus.Fields{
			"request": fReq,
			"status":  http.StatusInternalServerError,
			"Error":   err,
		}).Error("Unable Insert the Document!")

		fRes := &v1.FeedbackResponse{
			Api:         "v1",
			StatusCode:  http.StatusInternalServerError,
			RequestId:   "",
			Message:     fmt.Sprintln(err),
			FeedbackRes: nil,
		}
		return fRes, nil
	}

	fResult := coll.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: fReq.Id}})
	feedbackRes := v1.Feedback{}

	if err := fResult.Decode(&feedbackRes); err != nil {
		logrus.Errorf("Unable to read document for request id: %v", fReq.Id)
		return nil, nil
	}

	fRes := &v1.FeedbackResponse{
		Api:         "v1",
		StatusCode:  http.StatusCreated,
		RequestId:   fReq.Id,
		Message:     "Document Inserted Successfuly",
		FeedbackRes: &feedbackRes,
	}
	return fRes, nil
}

func (fs *feedbackServiceServer) Read(ctx context.Context, req *v1.ReadFeedbackRequest) (*v1.FeedbackResponse, error) {
	coll := fs.db.Collection("feedback_request")
	fResult := coll.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: req.RequestId}})
	feedbackRes := v1.Feedback{}
	fRes := &v1.FeedbackResponse{}
	fRes.Api = "v1"
	fRes.RequestId = req.RequestId

	if err := fResult.Decode(&feedbackRes); err != nil {
		logrus.Errorf("Unable to read document for request id: %v", req.RequestId)
		fRes.StatusCode = http.StatusNotFound
		fRes.Message = fmt.Sprint(fResult.Decode(bson.M{}))
		return fRes, nil
	}

	fRes.StatusCode = http.StatusCreated
	fRes.Message = "Document Inserted Successfuly"
	fRes.FeedbackRes = &feedbackRes

	return fRes, nil
}
func (fs *feedbackServiceServer) GenerateFeedbackForRequest(ctx context.Context, req *v1.GenerateFeedbackRequest) (*v1.GeneratedFeedbackResponse, error) {
	coll := fs.db.Collection("feedback_request")
	fResult := coll.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: req.FeedbackReq.Id}})
	gfRes := &v1.GeneratedFeedbackResponse{}
	gfRes.Api = "v1"

	fRes := v1.Feedback{}
	if err := fResult.Decode(&fRes); err != nil {
		logrus.Errorf("Unable to read document for request id: %v", req.FeedbackReq.Id)
	}
	var summaryText = req.SummaryNote + "\n\n"

	if fRes.IsProxy {
		summaryText += fmt.Sprintf("%s %s %s", candidate, was, feedbackMapping["Proxy"])
	} else {
		// Build feedback when coding is required.
		// TODO refactored and move it into separate methods
		if fRes.IsCodingRequired {
			if fRes.IsCodeCompiled && fRes.IsAlgoEfficient {
				summaryText += fmt.Sprintf("%s %s %s and %s.\n", candidate, has, feedbackMapping["CodeCompiled"], feedbackMapping["AlgoEfficient"])
			} else if fRes.IsCodeCompiled && !fRes.IsAlgoEfficient {
				summaryText += fmt.Sprintf("%s %s %s and %s.\n", candidate, has, feedbackMapping["CodeCompiled"], feedbackMapping["NotEfficient"])
			} else {
				if fRes.IsAbleToWritePseudoCode && fRes.IsAlgoEfficient {
					summaryText += fmt.Sprintf("%s %s %s and %s.\n", candidate, was, feedbackMapping["PseudoCode"], feedbackMapping["AlgoEfficient"])
				} else if fRes.IsAbleToWritePseudoCode && !fRes.IsAlgoEfficient {
					summaryText += fmt.Sprintf("%s %s %s and %s.\n", candidate, has, feedbackMapping["PseudoCode"], feedbackMapping["NotEfficient"])
				}
			}

			if fRes.FollowedCodingStandards {
				summaryText += fmt.Sprintf("%s %s %s.\n", candidate, was, feedbackMapping["Coding_Standards"])
			}
			if fRes.AnyCodingComment != "" {
				summaryText += fmt.Sprintf("%s.\n", fRes.AnyCodingComment)
			}

		}
		// Build feedback when Whiteboarding is required
		// TODO refactored and move it into separate methods
		if fRes.IsWhiteboardingRequired && fRes.IsWhiteboardQuestionAsked {
			if fRes.WhiteboardExplained {
				summaryText += fmt.Sprintf("%s %s.\n", candidate, feedbackMapping["Whiteboard_explained"])
			} else if fRes.WhiteboardPartial {
				summaryText += fmt.Sprintf("%s %s.\n", candidate, feedbackMapping["Whiteboard_explained"])
			} else {
				summaryText += fmt.Sprintf("%s %s.\n", candidate, feedbackMapping["Whiteboard_not_explained"])
			}
		}

		if fRes.MyComments != "" {
			summaryText += fmt.Sprintf("%s\n", fRes.MyComments)
		}

		// Build skill feedback
		// TODO refactored and move it into separate methods
		var sFeedbackSlice = []*v1.SkillFeedback{}
		for _, tech := range fRes.TechSkills {
			fmt.Println(tech.SkillName)
			if fRes.IsProxy {
				sFeedbackSlice = append(sFeedbackSlice, &v1.SkillFeedback{
					Skill:        tech.SkillName,
					FeedbackText: fmt.Sprintf("%s %s %s.\n", candidate, was, feedbackMapping["Proxy"]),
				})
			} else {
				sFeedback := &v1.SkillFeedback{
					Skill:        tech.SkillName,
					FeedbackText: fmt.Sprintf("%s. %s.\n", feedbackMapping["s-"+strconv.FormatInt(int64(tech.SkillRating), 10)], feedbackMapping["e-"+strconv.FormatInt(int64(tech.ExperienceRating), 10)]),
				}
				sFeedbackSlice = append(sFeedbackSlice, sFeedback)
				sFeedback.FeedbackText += fmt.Sprintf("\n%s\n", topicsCovered)
				//Build topics feedback
				for _, topic := range tech.Topics {
					sFeedback.FeedbackText += fmt.Sprintf("\n%s:\n", topic.TopicName)
					if topic.IsAbleToExaplain {
						sFeedback.FeedbackText += fmt.Sprintf("\n%s %s %s.\n", candidate, was, fmt.Sprintf(feedbackMapping["AbleToExplain"], topic.TopicName))
					} else if topic.PartiallyExplained {
						sFeedback.FeedbackText += fmt.Sprintf("%s %s %s.\n", candidate, was, feedbackMapping["PartiallyExplained"])
					}
					if topic.IsScenarioCovered {
						if topic.IsAbleToExplainScenario {
							sFeedback.FeedbackText += fmt.Sprintf("%s %s %s.\n", candidate, has, fmt.Sprintf(feedbackMapping["ScenarioExplained"], topic.WhatSceanrioQuestion))
						} else {
							sFeedback.FeedbackText += fmt.Sprintf("%s %s %s.\n", candidate, has, fmt.Sprintf(feedbackMapping["ScenarioNotExplained"], topic.WhatSceanrioQuestion))
						}
					}
					if topic.InDepthUnderstanding {
						sFeedback.FeedbackText += fmt.Sprintf("%s %s %s. \n", candidate, has, fmt.Sprintf(feedbackMapping["InDepthUnderstanding"], topic.TheoryQuestion))
					}

					if topic.IsHandsOn {
						sFeedback.FeedbackText += fmt.Sprintf("%s %s %s.\n", candidate, is, fmt.Sprintf(feedbackMapping["Hands-On-Topic"], topic.TopicName))
					}
				}

				if tech.InDepthUnderstanding {
					//fmt.Println(tech.QuestionsAsked)
					sFeedback.FeedbackText += fmt.Sprintf("\n%s %s %s. \n", candidate, has, fmt.Sprintf(feedbackMapping["InDepthUnderstanding"], tech.QuestionsAsked))
				}

				if tech.IsHandsOn {
					sFeedback.FeedbackText += fmt.Sprintf("\n%s %s %s.\n", candidate, is, feedbackMapping["Hands-On"])
				}

			}
		}
		gfRes.SkillFeedback = sFeedbackSlice
	}

	gfRes.SummaryText = summaryText
	gfRes.StatusCode = http.StatusOK

	gfRes.Message = "Generated feedback successfully"

	return gfRes, nil
}

func (fs *feedbackServiceServer) Delete(ctx context.Context, req *v1.DeleteFeedbackRequest) (*v1.DeleteFeedbackResponse, error) {
	coll := fs.db.Collection("feedback_request")
	result, err := coll.DeleteOne(ctx, bson.D{primitive.E{Key: "id", Value: req.RequestId}})

	if err != nil {
		fs.logger.Errorf("Unable to delete document for request id: %v Deleted Count: %v", req.RequestId, result.DeletedCount)
		return nil, nil
	}
	var fRes = &v1.DeleteFeedbackResponse{}
	fRes.Api = "v1"
	if result.DeletedCount == 0 {
		fRes.StatusCode = http.StatusNotFound
	} else {
		fRes.StatusCode = http.StatusOK
	}

	return fRes, nil
}
