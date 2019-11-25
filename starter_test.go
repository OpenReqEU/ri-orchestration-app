package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
	router = makeRouter()
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
