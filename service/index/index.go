package index

import (
	"fmt"
	"net/http"
	"review-service/pkg/router"
	"review-service/service/review/crawl/dappcom"
	crawl_dappradar_dapp "review-service/service/review/crawl/dappradar_dapp"
)

// GetIndex Function to Show API Information
func GetIndex(w http.ResponseWriter, r *http.Request) {
	router.ResponseSuccess(w, "200", "Go Framework is running")
}

// GetHealth Function to Show Health Check Status
func GetHealth(w http.ResponseWriter, r *http.Request) {
	router.HealthCheck(w)
}

type ScheduleInfo struct {
	Succes  int
	Fail    int
	Start   string
	TimeDif float64
	End     string
	Stage   string
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	info := dappcom.Tracking{
		LastUpdateTime: dappcom.LastUpdated,
		Update:         dappcom.UpdateCount,
		Insert:         dappcom.InsertCount,
	}
	router.ResponseCreatedWithData(w, "200", "Dapp Crawl Schedule Info", info)
}
func UpdateWebsiteDapp(w http.ResponseWriter, r *http.Request) {
	router.ResponseSuccess(w, fmt.Sprintf("succes: %d, fail: %d", crawl_dappradar_dapp.CountUpdateWebsite, crawl_dappradar_dapp.CountUpdateWebsiteFail), "Go Framework is running")
}
