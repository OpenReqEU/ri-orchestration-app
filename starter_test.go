package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/robfig/cron"

	"github.com/gorilla/mux"
)

var router *mux.Router

func TestMain(m *testing.M) {
	fmt.Println("--- Start Tests")
	setup()

	// run the test cases defined in this file
	retCode := m.Run()

	tearDown()

	// call with result of m.Run()
	os.Exit(retCode)
}

func setup() {
	fmt.Println("--- --- setup")
	setupRouter()
}

func setupRouter() {
	router = mux.NewRouter()
	router.HandleFunc("/hitec/orchestration/app/observe/google-play/package-name/{package_name}/interval/{interval}", MockPostObserveAppGooglePlay).Methods("POST")
	router.HandleFunc("/hitec/orchestration/app/process/google-play/package-name/{package_name}", MockPostProcessAppGooglePlay).Methods("POST")
}

func tearDown() {
	fmt.Println("--- --- tear down")
}

func buildRequest(method, endpoint string, payload io.Reader, t *testing.T) *http.Request {
	req, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	return req
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func TestPostObserveAppGooglePlay(t *testing.T) {
	fmt.Println("start TestPostObserveAppGooglePlay")
	var method = "POST"
	var endpoint = "/hitec/orchestration/app/observe/google-play/package-name/%s/interval/%s"

	/*
	 * test for faillure
	 */
	endpointFail := fmt.Sprintf(endpoint, "", "fail")
	req := buildRequest(method, endpointFail, nil, t)
	rr := executeRequest(req)

	//Confirm the response has the right status code
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusMovedPermanently, status)
	}

	/*
	 * test for success
	 */
	endpointSuccess := fmt.Sprintf(endpoint, "eu.openreq", "monthly")
	req = buildRequest(method, endpointSuccess, nil, t)
	rr = executeRequest(req)

	//Confirm the response has the right status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
}

func MockPostObserveAppGooglePlay(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	packageName := params["package_name"]
	interval := params["interval"] // possible intervals: minutely, hourly, daily, monthly

	var ok bool
	allowedSpcialIntervals := map[string]bool{
		"minutely": true,
		"hourly":   true,
		"daily":    true,
		"weekly":   true,
		"monthly":  true,
	}
	_, err := cron.Parse(interval)
	if packageName == "" || (err != nil && !allowedSpcialIntervals[interval]) {
		ok = false
	} else {
		ok = true
	}
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: false, Message: "storage layer unreachable"})
		return
	}

	fmt.Printf("1.1 restart observation \n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: true, Message: "observation successfully initiated"})
}

func TestPostProcessTweets(t *testing.T) {
	fmt.Println("start TestPostProcessTweets")
	var method = "POST"
	var endpoint = "/hitec/orchestration/app/process/google-play/package-name/%s"
	/*
	 * test for faillure
	 */
	endpointFail := fmt.Sprintf(endpoint, "")
	req := buildRequest(method, endpointFail, nil, t)
	rr := executeRequest(req)

	//Confirm the response has the right status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusNotFound, status)
	}

	/*
	 * test for success
	 */
	endpointSuccess := fmt.Sprintf(endpoint, "eu.openreq")
	req = buildRequest(method, endpointSuccess, nil, t)
	rr = executeRequest(req)

	//Confirm the response has the right status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
}

func MockPostProcessAppGooglePlay(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	packageName := params["package_name"]

	var ok = packageName != ""
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: false, Message: "account name or language are empty"})
		return
	}

	fmt.Printf("1.1 restart observation \n")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: true, Message: "tweets successfully processed"})
}
