package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"log"
)

var baseURL = os.Getenv("BASE_URL")

const (
	// analytics layer
	endpointPostClassifyAppReviews = "/ri-analytics-classification-google-play-review/hitec/classify/domain/google-play-reviews/"

	// collection layer
	endpointPostCrawlAppReviewsGooglePlay = "/ri-collection-explicit-feedback-google-play-review/hitec/crawl/app-reviews/google-play/%s/limit/%d"
	// collection layer
	endpointPostCrawlAppPageGooglePlay = "/ri-collection-explicit-feedback-google-play-page/hitec/crawl/app-page/google-play/%s"

	// storage layer
	endpointPostObserveAppGooglePlay            = "/ri-storage-app/hitec/repository/app/observe/app/google-play/package-name/%s/interval/%s"
	endpointGetObservablesGooglePlay            = "/ri-storage-app/hitec/repository/app/observable/google-play"
	endpointPostAppReviewGooglePlay             = "/ri-storage-app/hitec/repository/app/store/app-review/google-play/"
	endpointPostAppPageGooglePlay               = "/ri-storage-app/hitec/repository/app/store/app-page/google-play/"
	endpointPosNonExistingtAppReviewsGooglePlay = "/ri-storage-app/hitec/repository/app/non-existing/app-review/google-play"
)

var client = getHTTPClient()

func getHTTPClient() *http.Client {
	pwd, _ := os.Getwd()
	caCert, err := ioutil.ReadFile(pwd + "/ca_chain.crt")
	timeout := time.Duration(2 * time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
		Timeout: timeout,
	}

	return client
}

// RESTPostStoreObserveAppGooglePlay returns ok
func RESTPostStoreObserveAppGooglePlay(packageName string, interval string) bool {
	for connectionTries := 3; connectionTries > 0; connectionTries-- {
		endpoint := fmt.Sprintf(endpointPostObserveAppGooglePlay, packageName, interval)
		url := baseURL + endpoint
		res, err := client.Post(url, "custom", nil)
		if err != nil {
			log.Printf("ERR %v\n", err)
			continue
		}
		if res.StatusCode == 200 {
			return true
		}
	}

	return false
}

// RESTGetObservablesGooglePlay retrieve all observables from the storage layer
func RESTGetObservablesGooglePlay() []ObservableGooglePlay {
	var obserables []ObservableGooglePlay

	res, err := client.Get(baseURL + endpointGetObservablesGooglePlay)
	if err != nil {
		fmt.Println("ERR", err)
		return obserables
	}

	err = json.NewDecoder(res.Body).Decode(&obserables)
	if err != nil {
		fmt.Println("ERR", err)
		return obserables
	}

	return obserables
}

// RESTGetAppPageGooglePlay retrieve all reviews from the collection layer
func RESTGetAppPageGooglePlay(packageName string) AppPageGooglePlay {
	var appPage AppPageGooglePlay

	endpoint := fmt.Sprintf(endpointPostCrawlAppPageGooglePlay, packageName)
	res, err := client.Get(baseURL + endpoint)
	if err != nil {
		fmt.Println("ERR", err)
		return appPage
	}

	err = json.NewDecoder(res.Body).Decode(&appPage)
	if err != nil {
		fmt.Println("ERR", err)
		return appPage
	}

	return appPage
}

// RESTGetAppReviewsGooglePlay retrieve all reviews from the collection layer
func RESTGetAppReviewsGooglePlay(packageName string, limit int) []AppReviewGooglePlay {
	var reviews []AppReviewGooglePlay

	endpoint := fmt.Sprintf(endpointPostCrawlAppReviewsGooglePlay, packageName, limit)
	res, err := client.Get(baseURL + endpoint)
	if err != nil {
		fmt.Println("ERR", err)
		return reviews
	}

	err = json.NewDecoder(res.Body).Decode(&reviews)
	if err != nil {
		fmt.Println("ERR", err)
		return reviews
	}

	return reviews
}

// RESTPostProcessAppReviewsGooglePlay sends the crawled reviews to the processing layer and retrieves app reviews including their ml classes
func RESTPostProcessAppReviewsGooglePlay(reviews []AppReviewGooglePlay) []AppReviewGooglePlay {
	var appReviews []AppReviewGooglePlay

	js := new(bytes.Buffer)
	json.NewEncoder(js).Encode(reviews)
	res, err := client.Post(baseURL+endpointPostClassifyAppReviews, "application/json; charset=utf-8", js)
	if err != nil {
		fmt.Println("ERR", err)
		return appReviews
	}

	err = json.NewDecoder(res.Body).Decode(&appReviews)
	if err != nil {
		fmt.Println("ERR", err)
		return appReviews
	}

	return appReviews
}

// RESTPostStoreProcessedAppReviewsGooglePlay sends the processed app reviews to the storage layer. Returns ok MS could be reached
func RESTPostStoreProcessedAppReviewsGooglePlay(appReviews []AppReviewGooglePlay) bool {
	js := new(bytes.Buffer)
	json.NewEncoder(js).Encode(reviews)
	_, err := client.Post(baseURL+endpointPostAppReviewGooglePlay, "application/json; charset=utf-8", js)
	if err != nil {
		fmt.Println("ERR", err)
		return false
	}

	return true
}

// RESTPostStoreAppPageGooglePlay sends the crawled app page to the storage layer. Returns ok MS could be reached
func RESTPostStoreAppPageGooglePlay(appPage AppPageGooglePlay) bool {
	js := new(bytes.Buffer)
	json.NewEncoder(js).Encode(appPage)
	_, err := client.Post(baseURL+endpointPostAppPageGooglePlay, "application/json; charset=utf-8", js)
	if err != nil {
		fmt.Println("ERR", err)
		return false
	}

	return true
}

// RESTPostNonExistingAppReviewsGooglePlay sends the crawled app reviews and gets a list of app reviews in return that do not yet exist in the db.
func RESTPostNonExistingAppReviewsGooglePlay(appReviews []AppReviewGooglePlay) []AppReviewGooglePlay {
	var nonExistingAppReviews []AppReviewGooglePlay

	js := new(bytes.Buffer)
	json.NewEncoder(js).Encode(appReviews)
	res, err := client.Post(baseURL+endpointPosNonExistingtAppReviewsGooglePlay, "application/json; charset=utf-8", js)
	if err != nil {
		fmt.Println("ERR", err)
		return nonExistingAppReviews
	}

	err = json.NewDecoder(res.Body).Decode(&nonExistingAppReviews)
	if err != nil {
		fmt.Println("ERR", err)
		return nonExistingAppReviews
	}

	return nonExistingAppReviews
}
