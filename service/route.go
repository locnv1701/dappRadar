package service

import (
	"fmt"
	"review-service/pkg/router"
	"review-service/service/index"
	"review-service/service/review"
	"review-service/service/review/crawl/dappcom"
	// "review-service/service/review/crawl/dappcom"
	//2
	//3
	// crawl_ico_drop "review-service/service/review/crawl/ico_drop"
)

func init() {
	go func() {
		//1
		// dappcom.Crawl()

		//2
		// fmt.Println("start")
		// crawl_dappradar_dapp.CrawlUpdateDapp()

		//3
		// fmt.Println("start")
		// crawl_ico_drop.GetListIcoDrop()

		// //************************************************************************************************

		// crawl_dappradar_dapp.Schedule()
		// dappcom.Crawl()
		// dappcom.UpdateDappRadarByDappCom()
		dappcom.ScheduleUpdateDappRadarByDappCom()

	}()
}

// LoadRoutes to Load Routes to Router
func LoadRoutes() {
	fmt.Println("routes")
	// Set Endpoint for admin
	router.Router.Get(router.RouterBasePath+"/", index.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", index.GetHealth)

	router.Router.Get(router.RouterBasePath+"/updateWebsiteDapp", index.UpdateWebsiteDapp)
	router.Router.Get(router.RouterBasePath+"/info", index.GetInfo)

	router.Router.Mount(router.RouterBasePath+"/review", review.ExchangeInfoServiceSubRoute)
}
