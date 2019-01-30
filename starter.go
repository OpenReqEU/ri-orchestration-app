package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/mux"
)

func main() {
	log.SetOutput(os.Stdout)

	router := mux.NewRouter()
	router.HandleFunc("/hitec/orchestration/app/observe/google-play/package-name/{package_name}/interval/{interval}", postObserveAppGooglePlay).Methods("POST")
	router.HandleFunc("/hitec/orchestration/app/process/google-play/package-name/{package_name}", postProcessAppGooglePlay).Methods("POST")

	log.Fatal(http.ListenAndServe(":9702", router))
}

/*
* This method calls for each step the reponsible MS
*
* Steps:
*  1. store app to observe
*  2. notify the observer (crawler)
*  2.1 notify the processing layer to classify the newly addded reviews
 */
func postObserveAppGooglePlay(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	packageName := params["package_name"]
	interval := params["interval"] // possible intervals: minutely, hourly, daily, monthly

	// 1. store app to observe
	ok := RESTPostStoreObserveAppGooglePlay(packageName, interval)
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: false, Message: "storage layer unreachable"})
		return
	}

	// 2. notify the observer (crawler)
	RestartObservation()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: true, Message: "observation successfully initiated"})
}

/*
* This method calls for each step the reponsible MS
*
* Steps:
*  1. crawl app page
*  2. store app page
*  2. crawl app reviews
*  3. process reviews
*  4. store processed app reviews
 */
func postProcessAppGooglePlay(w http.ResponseWriter, r *http.Request) {
	fmt.Println("postProcessAppGooglePlay called")
	params := mux.Vars(r)
	packageName := params["package_name"]

	//  1. crawl app page
	fmt.Println("1. crawl app page")
	appPage := RESTGetAppPageGooglePlay(packageName)

	//  2. store app page
	fmt.Println("2. store app page")
	ok := RESTPostStoreAppPageGooglePlay(appPage)
	fmt.Println("Could store app page", ok)

	//  3. crawl app reviews
	fmt.Println("3. crawl app reviews")
	crawledAppReviews := RESTGetAppReviewsGooglePlay(packageName, 0)
	nonExistingAppReviews := RESTPostNonExistingAppReviewsGooglePlay(crawledAppReviews) // just consider app reviews that are not processed yet

	//  4. process reviews
	fmt.Println("4. process reviews")
	processedAppReviess := RESTPostProcessAppReviewsGooglePlay(nonExistingAppReviews)

	//  5. store processed app reviews
	fmt.Println("5. store processed app reviews")
	RESTPostStoreProcessedAppReviewsGooglePlay(processedAppReviess)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: true, Message: "crawled, processed, and stored app reviews"})
}
