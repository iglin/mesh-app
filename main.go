package main

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	appId   = uuid.NewString()
	appName = "mesh-app"
)

type TargetSvcRequest struct {
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
}

type TargetSvcResponse struct {
	RequestDetails
	ResponseHeaders http.Header `json:"responseHeaders"`
}

type RequestDetails struct {
	AppId             string             `json:"appId"`
	AppName           string             `json:"appName"`
	Url               string             `json:"url"`
	Host              string             `json:"host"`
	Method            string             `json:"method"`
	Proto             string             `json:"proto"`
	Headers           http.Header        `json:"headers"`
	Body              string             `json:"body"`
	Cookies           []*http.Cookie     `json:"cookies"`
	TargetSvcResponse *TargetSvcResponse `json:"targetSvcResponse"`
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
		_, _ = w.Write([]byte("Failed to read request body due to error: " + err.Error()))
		return
	}
	bodyStr := string(body)

	targetSvcResponse := sendRequestIfNeeded(bodyStr)

	reqDetails := RequestDetails{
		AppId:             appId,
		AppName:           appName,
		Url:               req.URL.String(),
		Host:              req.Host,
		Method:            req.Method,
		Proto:             req.Proto,
		Headers:           req.Header,
		Body:              bodyStr,
		Cookies:           req.Cookies(),
		TargetSvcResponse: targetSvcResponse,
	}
	log.Infof("App %s with id %v got request: %v", appName, appId, reqDetails)
	response, err := json.Marshal(&reqDetails)
	if err != nil {
		log.Errorf("Failed to marshall response json:\n %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to marshall response json due to error: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func sendRequestIfNeeded(body string) *TargetSvcResponse {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil
	}
	var reqSpec TargetSvcRequest
	if err := json.Unmarshal([]byte(body), &reqSpec); err != nil {
		log.Errorf("Failed to unmarshall TargetSvcRequest from json:\n %v", err)
		return nil
	}
	targetUrl, err := url.Parse(reqSpec.URL)
	if err != nil {
		log.Errorf("Failed to parse target URL:\n %v", err)
		return nil
	}
	req := &http.Request{
		Method: reqSpec.Method,
		URL:    targetUrl,
		Header: reqSpec.Headers,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Request to target service failed:\n %v", err)
		return nil
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to read target service response body:\n %v", err)
	}
	var responseBody TargetSvcResponse
	if err := json.Unmarshal(respBody, &responseBody); err != nil {
		log.Errorf("Failed to unmarshall TargetSvcResponse from json:\n %v", err)
		return nil
	}
	return &responseBody
}
