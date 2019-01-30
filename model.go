package main

// ObservableGooglePlay model
type ObservableGooglePlay struct {
	PackageName string `json:"package_name" bson:"package_name"`
	Interval    string `json:"interval" bson:"interval"`
}

// AppPageGooglePlay model
type AppPageGooglePlay struct {
	Name                    string             `json:"name" bson:"name"`
	PackageName             string             `json:"package_name" bson:"package_name"`
	DateCrawled             int64              `json:"date_crawled" bson:"date_crawled"`
	Category                string             `json:"category" bson:"category"`
	USK                     string             `json:"usk" bson:"usk"`
	Price                   string             `json:"price" bson:"price"`
	PriceValue              float64            `json:"price_value" bson:"price_value"`
	PriceCurrency           string             `json:"price_currency" bson:"price_currency"`
	Description             string             `json:"description" bson:"description"`
	WhatsNew                []string           `json:"whats_new" bson:"whats_new"`
	Rating                  float64            `json:"rating" bson:"rating"`
	StarsCount              int64              `json:"stars_count" bson:"stars_count"`
	CountPerRating          StarCountPerRating `json:"count_per_rating" bson:"count_per_rating"`
	EstimatedDownloadNumber int64              `json:"estimated_download_number" bson:"estimated_download_number"`
	DeveloperName           string             `json:"developer" bson:"developer"`
	TopDeveloper            bool               `json:"top_developer" bson:"top_developer"`
	ContainsAds             bool               `json:"contains_ads" bson:"contains_ads"`
	InAppPurchases          bool               `json:"in_app_purchase" bson:"in_app_purchase"`
	LastUpdate              int64              `json:"last_update" bson:"last_update"`
	Os                      string             `json:"os" bson:"os"`
	RequiresOsVersion       string             `json:"requires_os_version" bson:"requires_os_version"`
	CurrentSoftwareVersion  string             `json:"current_software_version" bson:"current_software_version"`
	SimilarApps             []string           `json:"similar_apps" bson:"similar_apps"`
}

// StarCountPerRating model
type StarCountPerRating struct {
	Five  int `json:"5"`
	Four  int `json:"4"`
	Three int `json:"3"`
	Two   int `json:"2"`
	One   int `json:"1"`
}

// AppReviewGooglePlay model
type AppReviewGooglePlay struct {
	ReviewID       string `json:"review_id" bson:"review_id"`
	PackageName    string `json:"package_name" bson:"package_name"`
	Author         string `json:"author" bson:"author"`
	Date           int64  `json:"date_posted" bson:"date_posted"`
	Rating         int    `json:"rating" bson:"rating"`
	Title          string `json:"title" bson:"title"`
	Body           string `json:"body" bson:"body"`
	PermaLink      string `json:"perma_link" bson:"perma_link"`
	FeatureRequest bool   `json:"cluster_is_feature_request" bson:"cluster_is_feature_request"`
	BugReport      bool   `json:"cluster_is_bug_report" bson:"cluster_is_bug_report"`
}

// Response model
type Response struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}
