package model

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//FeedbackRequest will be used to store in MongoDB
type FeedbackRequest struct {
	ID                      primitive.ObjectID   `bson:"_id,omitempty"`
	API                     string               `bson:"api,omitempty"`
	CandidateName           string               `bson:"candidate_name,omitempty"`
	RecruiterName           string               `bson:"recruiter_name,omitempty"`
	JobType                 string               `bson:"job_type,omitempty"`
	IsProxy                 bool                 `bson:"is_proxy,omitempty"`
	IsCodingRequired        bool                 `bson:"is_coding_required,omitempty"`
	IsWhiteBoardingRequired bool                 `bson:"whtbrd_required,omitempty"`
	IsIDProofRequired       bool                 `bson:"id_required,omitempty"`
	IsCodeCompiled          bool                 `bson:"code_compiled,omitempty"`
	IsAbleToWritePseudoCode bool                 `bson:"is_able_to_write_pseudoCode,omitempty"`
	IsCompletedWhiteBoard   bool                 `bson:"completed_whiteboard,omitempty"`
	IsAlgoEfficient         bool                 `bson:"algo_efficient,omitempty"`
	TechSkills              []*TechSkill         `bson:"techskills,omitempty"`
	Notes                   string               `bson:"notes,omitempty"`
	CreatDate               *timestamp.Timestamp `bson:"createDate,omitempty"`
	UpdateDate              *timestamp.Timestamp `bson:"updateDate,omitempty"`
}

//TechSkill covered during the discussion
type TechSkill struct {
	SkillName        string    `bson:"skill_name,omitempty"`
	SkillRating      int32     `bson:"skill_rating,omitempty"`
	ExperienceRating int32     `bson:"experience_rating,omitempty"`
	Topics           []*Topics `bson:"topics,omitempty"`
}

//Topics to capture different topics covered in the skillset
type Topics struct {
	TopicName                 string `bson:"topic_name,omitempty"`
	IsScenarioCovered         bool   `bson:"is_scenario_covered,omitempty"`
	WhatSceanrioQuestion      string `bson:"sceanrio_question,omitempty"`
	IsAbleToExplainScenario   bool   `bson:"able_to_explain_scenario,omitempty"`
	IsAbleToExaplain          bool   `bson:"is_able_to_exaplain,omitempty"`
	InDepthUnderstanding      bool   `bson:"in_depth_understanding,omitempty"`
	PartiallyExplained        bool   `bson:"partially_explained,omitempty"`
	IsHandsOn                 bool   `bson:"is_hands_on,omitempty"`
	HaveTheroreticalKnowledge bool   `bson:"have_theroretical_knowledge,omitempty"`
}
