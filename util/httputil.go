package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//APITester is a warapper round http client
type APITester struct {
	baseURL  string
	userID   string
	password string
	jwtToken string
	client   *http.Client
	verbose  bool
}

//NewAPITester builds a new instance of APITester
func NewAPITester(baseURL, userID, password string, verbose bool) *APITester {
	at := new(APITester)
	at.Init(baseURL, userID, password, verbose)
	return at
}

//Init initialized the struct
func (at *APITester) Init(baseURL, userID, password string, verbose bool) error {
	at.baseURL = baseURL
	at.userID = userID
	at.password = password
	at.client = &http.Client{
		Timeout: time.Second * 10,
	}
	at.verbose = verbose
	return nil
}

//PerformHTTPCall peforms a json posting
func (at *APITester) PerformHTTPCall(method, url string, data interface{}) (int, map[string]interface{}) {
	var request *http.Request
	if data != nil || method != "GET" {
		jsonBytes, _ := json.MarshalIndent(data, "", " ")
		if at.verbose {
			log.Printf("Request object \n%s\n", string(jsonBytes))
		}
		request, _ = http.NewRequest(method, fmt.Sprintf("%s%s", at.baseURL, url), bytes.NewBuffer(jsonBytes))
	} else {
		request, _ = http.NewRequest(method, fmt.Sprintf("%s%s", at.baseURL, url), nil)
	}
	request.Header.Set("Content-type", "application/json")
	resp, err := at.client.Do(request)
	if err != nil {
		log.Printf("Response error %v", err)
		return -1, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Response body read error %v", err)
		return -1, nil
	}
	respnMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respnMap)
	if err != nil && len(body) > 0 {
		respnMap["payload"] = string(body)
	}
	return resp.StatusCode, respnMap
}
