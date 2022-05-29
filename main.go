package main

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	appId   = uuid.NewString()
	appName = "mesh-app"
)

type RequestDetails struct {
	AppId   string         `json:"appId"`
	AppName string         `json:"appName"`
	Url     string         `json:"url"`
	Host    string         `json:"host"`
	Method  string         `json:"method"`
	Proto   string         `json:"proto"`
	Headers http.Header    `json:"headers"`
	Body    string         `json:"body"`
	Cookies []*http.Cookie `json:"cookies"`
}

func main() {
	if providedName := os.Getenv("APP_NAME"); providedName != "" {
		appName = providedName
	}
	log.Infof("Running app %s with id %v on port 8080", appName, appId)

	http.HandleFunc("/", handleRequest)

	log.Panicf("%v", http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("Failed to read request body:\n %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to read request body due to error: " + err.Error()))
		return
	}

	reqDetails := RequestDetails{
		AppId:   appId,
		AppName: appName,
		Url:     req.URL.String(),
		Host:    req.Host,
		Method:  req.Method,
		Proto:   req.Proto,
		Headers: req.Header,
		Body:    string(body),
		Cookies: req.Cookies(),
	}
	log.Infof("App %s with id %v got request: %v", appName, appId, reqDetails)
	response, err := json.Marshal(&reqDetails)
	if err != nil {
		log.Errorf("Failed to marshall response json:\n %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to marshall response json due to error: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
