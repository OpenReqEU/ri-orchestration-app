package main

import (
	"github.com/robfig/cron"
)

var observableAppsGooglePlay = NewSet()
var reviews []AppReviewGooglePlay
var observer *cron.Cron

func startObsevation() {
	loadObservableApps()

	observer = cron.New()
	for packageName := range observableAppsGooglePlay.m {
		observerInterval := getObserverInterval(observableAppsGooglePlay.m[packageName])
		observer.AddFunc(observerInterval, func() {
			crawledAppReviews := crawlObservableApps(packageName)
			nonExistingAppReviews := RESTPostNonExistingAppReviewsGooglePlay(crawledAppReviews) // just consider app reviews that are not processed yet
			processedAppReviews := processObservableApps(nonExistingAppReviews)
			storeProcessedApps(processedAppReviews)
		})
	}
	observer.Start()
}

func stopObservation() {
	observer.Stop()
}

func loadObservableApps() {
	for _, observable := range RESTGetObservablesGooglePlay() {
		observableAppsGooglePlay.Add(observable.PackageName, observable.Interval)
	}
}

func getObserverInterval(interval string) string {
	switch interval {
	case "minutely":
		return "0 * * * * *"
	case "hourly":
		return "@hourly"
	case "daily":
		return "@daily"
	case "weekly":
		return "@weekly"
	case "monthly":
		return "@monthly"
	default:
		return interval // allows custom intervals to the cron job specification (https://godoc.org/github.com/robfig/cron) might thorw an error if the custom interval is wrong
	}
}

func crawlObservableApps(packageName string) []AppReviewGooglePlay {
	return RESTGetAppReviewsGooglePlay(packageName, 0)
}

func processObservableApps(appReviews []AppReviewGooglePlay) []AppReviewGooglePlay {
	return RESTPostProcessAppReviewsGooglePlay(appReviews)
}

func storeProcessedApps(processedAppReviews []AppReviewGooglePlay) {
	RESTPostStoreProcessedAppReviewsGooglePlay(processedAppReviews)
}

// RestartObservation stops the observation and starts it again
func RestartObservation() {
	if observer != nil {
		stopObservation()
	}
	startObsevation()
}
