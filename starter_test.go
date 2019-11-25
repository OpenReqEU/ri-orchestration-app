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
var stopTestServer func()

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
	setupMockClient()
}

func setupMockClient() {
	fmt.Println("Mocking client")
	handler := makeMockHandler()
	s := httptest.NewServer(handler)
	stopTestServer = s.Close
	baseURL = s.URL
}

func makeMockHandler() http.Handler {
	r := mux.NewRouter()
	mockAnalyticsClassificationGooglePlayReview(r)
	mockCollectionExplicitFeedbackGooglePlayReview(r)
	mockCollectionExplicitFeedbackGooglePlayPage(r)
	mockStorageApp(r)
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(errors.Errorf("Service method not mocked: %s", r.URL))
		w.WriteHeader(http.StatusNotFound)
	})
	return r
}

func mockAnalyticsClassificationGooglePlayReview(r *mux.Router) {
	// endpointPostClassifyAppReviews = "/ri-analytics-classification-google-play-review/hitec/classify/domain/google-play-reviews/"
	r.HandleFunc("/ri-analytics-classification-google-play-review/hitec/classify/domain/google-play-reviews/", func(w http.ResponseWriter, request *http.Request) {
		respond(w, http.StatusOK, `[]`)
	})
}

func mockCollectionExplicitFeedbackGooglePlayReview(r *mux.Router) {
	// endpointPostCrawlAppReviewsGooglePlay = "/ri-collection-explicit-feedback-google-play-review/hitec/crawl/app-reviews/google-play/%s/limit/%d"
	r.HandleFunc("/ri-collection-explicit-feedback-google-play-review/hitec/crawl/app-reviews/google-play/{package_name}/limit/{limit}", func(w http.ResponseWriter, request *http.Request) {
		respond(w, http.StatusOK, `[]`)
	})
}

func mockCollectionExplicitFeedbackGooglePlayPage(r *mux.Router) {
	// endpointPostCrawlAppPageGooglePlay = "/ri-collection-explicit-feedback-google-play-page/hitec/crawl/app-page/google-play/%s"
	r.HandleFunc("/ri-collection-explicit-feedback-google-play-page/hitec/crawl/app-page/google-play/{package_name}", func(w http.ResponseWriter, request *http.Request) {
		respond(w, http.StatusOK, `[]`)
	})
}

func mockStorageApp(r *mux.Router) {
	// endpointPostObserveAppGooglePlay = "/ri-storage-app/hitec/repository/app/observe/app/google-play/package-name/%s/interval/%s"
	r.HandleFunc("/ri-storage-app/hitec/repository/app/observe/app/google-play/package-name/{package_name}/interval/{interval}", func(w http.ResponseWriter, request *http.Request) {
		respond(w, http.StatusOK, nil)
	})

	// endpointGetObservablesGooglePlay = "/ri-storage-app/hitec/repository/app/observable/google-play"
	r.HandleFunc("/ri-storage-app/hitec/repository/app/observable/google-play", func(w http.ResponseWriter, request *http.Request) {
		respond(w, http.StatusOK, []interface{}{map[string]string{
			"package_name": "eu.openreq",
			"interval":     "midnight",
		}})
	})

	// endpointPostAppReviewGooglePlay = "/ri-storage-app/hitec/repository/app/store/app-review/google-play/"
	r.HandleFunc("/ri-storage-app/hitec/repository/app/store/app-review/google-play/", func(w http.ResponseWriter, request *http.Request) {
		respond(w, http.StatusOK, nil)
	})

	// endpointPostAppPageGooglePlay = "/ri-storage-app/hitec/repository/app/store/app-page/google-play/"
	r.HandleFunc("/ri-storage-app/hitec/repository/app/store/app-page/google-play/", func(w http.ResponseWriter, request *http.Request) {
		respond(w, http.StatusOK, nil)
	})

	// endpointPosNonExistingtAppReviewsGooglePlay = "/ri-storage-app/hitec/repository/app/non-existing/app-review/google-play"
	r.HandleFunc("/ri-storage-app/hitec/repository/app/non-existing/app-review/google-play", func(w http.ResponseWriter, request *http.Request) {
		respond(w, http.StatusOK, `[]`)
	})
}

func respond(writer http.ResponseWriter, statusCode int, body interface{}) {
	var bodyData []byte
	var err error
	if body == nil {
		bodyData = make([]byte, 0)
	} else {
		switch body.(type) {
		case string:
			bodyData = []byte(body.(string))
		case []byte:
			bodyData = body.([]byte)
		default:
			bodyData, err = json.Marshal(body)
			if err != nil {
				panic(err)
			}
		}
	}

	writer.WriteHeader(statusCode)

	if _, err = writer.Write(bodyData); err != nil {
		panic(err)
	}
}

func tearDown() {
	fmt.Println("--- --- tear down")
	stopTestServer()
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
