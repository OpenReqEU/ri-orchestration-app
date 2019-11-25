package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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

type endpoint struct {
	method string
	url    string
}

func (e endpoint) withVars(vs ...interface{}) endpoint {
	e.url = fmt.Sprintf(e.url, vs...)
	return e
}

func (e endpoint) executeRequest(payload interface{}) (error, *httptest.ResponseRecorder) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(payload)
	if err != nil {
		return err, nil
	}

	req, err := http.NewRequest(e.method, e.url, body)
	if err != nil {
		return err, nil
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return nil, rr
}

func (e endpoint) mustExecuteRequest(payload interface{}) *httptest.ResponseRecorder {
	err, rr := e.executeRequest(payload)
	if err != nil {
		panic(errors.Wrap(err, `Could not execute request`))
	}

	return rr
}

func isSuccess(code int) bool {
	return code >= 200 && code < 300
}

func assertSuccess(t *testing.T, rr *httptest.ResponseRecorder) {
	if !isSuccess(rr.Code) {
		t.Errorf("Status code differs. Expected success.\n Got status %d (%s) instead", rr.Code, http.StatusText(rr.Code))
	}
}
func assertFailure(t *testing.T, rr *httptest.ResponseRecorder) {
	if isSuccess(rr.Code) {
		t.Errorf("Status code differs. Expected failure.\n Got status %d (%s) instead", rr.Code, http.StatusText(rr.Code))
	}
}

func TestPostObserveAppGooglePlay(t *testing.T) {
	ep := endpoint{method: "POST", url: "/hitec/orchestration/app/observe/google-play/package-name/%s/interval/%s"}
	assertFailure(t, ep.withVars("", "fail").mustExecuteRequest(nil))
	assertSuccess(t, ep.withVars("eu.openreq", "monthly").mustExecuteRequest(nil))
}

func TestPostProcessTweets(t *testing.T) {
	ep := endpoint{method: "POST", url: "/hitec/orchestration/app/process/google-play/package-name/%s"}
	assertFailure(t, ep.withVars("").mustExecuteRequest(nil))
	assertSuccess(t, ep.withVars("eu.openreq").mustExecuteRequest(nil))
}
