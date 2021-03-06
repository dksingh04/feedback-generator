package main

import (
	"context"
	"encoding/json"
	"feedback-generator/internal/config"
	f "feedback-generator/pkg/api/v1/feedbackreqpb"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

var serverAddr = "localhost:9090"
var tpl *template.Template
var cc *clientConfig

type clientConfig struct {
	logger     *logrus.Logger
	clientConn *grpc.ClientConn
	client     f.FeedbackServiceClient
}

var pwd = ""
var fb = &f.Feedback{TechSkills: []*f.TechSkill{&f.TechSkill{Topics: []*f.Topic{&f.Topic{}}}}}

//TopicData for adding topic dynamically
type TopicData struct {
	SkillIndex int
	Topics     []*f.Topic
}

func init() {
	var err error
	pwd, err = os.Getwd()
	if err != nil {
		logrus.Fatal("Unable to get Working Directory.. ", err)
	}
	tpl = template.Must(template.New("feedback-template").ParseGlob(pwd + "/public/templates/*.html"))
}

func main() {
	logger, err := config.CreateDefaultLogConfiguration()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"filename": "logger",
			"status":   500,
			"Error":    err,
		}).Fatal("Unable to create the default logger configuration!")
	}
	// Read config information
	conf, _ := config.ReadConfig()

	//Connect to grpc server
	conn := connectToGRPCServer(logger, conf)
	defer conn.Close()

	cc = initializeClientConfig(logger, conn)

	mux := http.NewServeMux()
	//router.Ha
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.FromSlash(pwd+"/public/css/")))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(filepath.FromSlash(pwd+"/public/js/")))))
	imgServer := http.FileServer(http.Dir(filepath.FromSlash(pwd + "/public/img/")))
	mux.Handle("/img/", http.StripPrefix("/img/", imgServer))
	mux.HandleFunc("/addskill/", addSkill)
	mux.HandleFunc("/addtopic/", addTopic)
	mux.HandleFunc("/", index)
	logger.Info("Started client server..")
	log.Fatal(http.ListenAndServe(":"+conf.ClientPort, mux))

}
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.Method == http.MethodGet {
		fmt.Println("Called Index Get method!!")
		tpl.ExecuteTemplate(w, "index.html", fb)
	} else {
		//It is POST method
		var t *f.Feedback
		json.NewDecoder(r.Body).Decode(&t)
		//fmt.Println(t.TechSkills[0].Topics[0])
		fRes, err := generateFeedbackResponseFromRequest(t)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"status":  http.StatusInternalServerError,
				"message": "Unable to process the feedback request and unable to generate response",
				"Error":   err,
			}).Fatalf("Unable to process the message error: %v", err)
		}
		//fReport, err := json.Marshal(fRes)

		//Write the report to response writer
		//w.Write(fReport)
		tpl.ExecuteTemplate(w, "report.html", fRes)
		//fmt.Println(string(fReport))
	}

}

func addSkill(w http.ResponseWriter, r *http.Request) {
	var techSkill = &f.TechSkill{Topics: []*f.Topic{&f.Topic{}}}
	var t *f.Feedback

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":  http.StatusInternalServerError,
			"message": "Unable to decode request received in the request body",
			"Error":   err,
		}).Errorf("Unable to process the message error: %v", err)
	}
	fb = t
	fb.TechSkills = append(fb.TechSkills, techSkill)
	fmt.Println(len(fb.TechSkills))
	tpl.ExecuteTemplate(w, "skill.html", fb)
}

func addTopic(w http.ResponseWriter, r *http.Request) {
	var topic = &f.Topic{}
	var t *f.Feedback
	//Always parse request body before retrieving any form or URL parameter value.
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":  http.StatusInternalServerError,
			"message": "Unable to decode request received in the request body",
			"Error":   err,
		}).Errorf("Unable to process the message error: %v", err)
	}

	r.ParseForm()
	index, err := strconv.Atoi(r.Form.Get("index"))
	fmt.Println(index)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status": 500,
			"Error":  err,
		}).Error("Parameters is missing!")
	}

	fb = t
	fb.TechSkills[index].Topics = append(fb.TechSkills[index].Topics, topic)
	tData := &TopicData{SkillIndex: index, Topics: fb.TechSkills[index].Topics}
	fmt.Println(len(fb.TechSkills[index].Topics))
	tpl.ExecuteTemplate(w, "topic.html", tData)
}

func generateFeedbackResponseFromRequest(feedback *f.Feedback) (*f.GeneratedFeedbackResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fReq := &f.GenerateFeedbackRequest{Api: "v1", FeedbackReq: feedback, SummaryNote: ""}
	fRes, err := cc.client.GenerateFeedbackFromFormData(ctx, fReq)
	/*fmt.Println("Professional Summary:")
	fmt.Println("")
	fmt.Printf(fRes.SummaryText)
	fmt.Println("")
	fmt.Printf("%s:\n\n", fRes.SkillFeedback[0].Skill)
	fmt.Printf(fRes.SkillFeedback[0].FeedbackText)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Error:")
	fmt.Println(err)
	*/
	fb = feedback
	return fRes, err
}

func connectToGRPCServer(logger *logrus.Logger, conf *config.Config) *grpc.ClientConn {
	conn, err := grpc.Dial(conf.GRPCServer+":"+conf.GRPCPort, grpc.WithInsecure())
	if err != nil {
		logger.WithFields(logrus.Fields{
			"server": "localhost:9090",
			"status": 500,
			"Error":  err,
		}).Fatalln("Unable to connect to grpc server")
	}
	return conn
}
func initializeClientConfig(logger *logrus.Logger, conn *grpc.ClientConn) *clientConfig {
	client := f.NewFeedbackServiceClient(conn)

	cc := &clientConfig{
		logger:     logger,
		clientConn: conn,
		client:     client,
	}

	return cc
}
