package crawl_dappradar_dapp

import (
	"fmt"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"review-service/service/constant"
	"review-service/service/review/model/dao"
	dto_dappradar "review-service/service/review/model/dto/dappradar"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//crawl javascript render
//ref: https://www.devdungeon.com/content/web-scraping-go
//lib: https://github.com/geziyor/geziyor

var MapChainList map[string]string

var MapDappradar map[string]dto_dappradar.DetailDapp

var CountUpdateWebsite int
var CountUpdateWebsiteFail int

func CrawlUpdateDapp() {

	MapDappradar = make(map[string]dto_dappradar.DetailDapp)

	var listDappradar dto_dappradar.ListDappradar
	err := listDappradar.GetAllDappradar("dappradar")
	if err != nil {
		log.Println(log.LogLevelError, `listDappradar.GetAllDappradar(dappradar)`, err.Error())
		return
	}

	for _, dapp := range listDappradar.DappList {
		MapDappradar[dapp.DAppCode] = dapp
	}

	log.Println(log.LogLevelInfo, `Get All Dappradar `, len(MapDappradar))

	chainList := dao.ChainList{}

	err = chainList.GetChainList()
	if err != nil {
		log.Println(log.LogLevelError, `chainList.GetChainList()`, err.Error())
	} else {
		MapChainList = make(map[string]string, 0)

		for _, chain := range chainList.List {
			MapChainList[chain.Chainname] = chain.ChainId
		}
	}

	if len(MapChainList) == 0 {
		return
	}

	log.Println(log.LogLevelInfo, `Get chainlist succes`, len(MapChainList))

	//init chainlist, map exist dapp

	url := fmt.Sprintf(`%s%s`, constant.BASE_URL_DAPPRADAR, constant.ENDPOINT_CHAIN_ALL_DAPPRADAR)
CallAPIBlockchain:
	fmt.Println("1")
	dom := utils.GetHtmlDomJsRenderByUrl(url)
	fmt.Println("2")
	if dom == nil {
		time.Sleep(5 * time.Second)
		goto CallAPIBlockchain
	}
	endpointBlockChains := CrawlEndpointBlockchainsByHtmlDom(dom)
	if len(endpointBlockChains) == 0 {
		time.Sleep(5 * time.Second)
		goto CallAPIBlockchain
	}

	fmt.Println("endpointBlockChains", len(endpointBlockChains))

	for _, endpoint := range endpointBlockChains {
		fmt.Print(endpoint.BlockchainId, endpoint.BlockchainName, endpoint.Endpoint, "===")
	}

	for index, endpointBlockChain := range endpointBlockChains {
		for pageIdx := 1; ; pageIdx++ {
			url := fmt.Sprintf(`%s%v/%d`, constant.BASE_URL_DAPPRADAR, endpointBlockChain.Endpoint, pageIdx)
			fmt.Println("url next", url)
		CallAPIListPagination:
			dom := utils.GetHtmlDomJsRenderByUrl(url)
			if dom == nil {
				log.Println(log.LogLevelDebug, `service/review/crawl/dappradar_dapp/Crawl/getDomJsLoad`, `dom get by js loading is nil`)
				pageIdx--
				time.Sleep(5 * time.Second)
				continue
			}
			if IsEndPage(dom) {
				break
			}
			fmt.Println("done get dom")
			succes, endpointDappList := CrawlDappEndpointByHtmlDom(dom)
			//Response without data
			if !succes {
				time.Sleep(5 * time.Second)
				goto CallAPIListPagination
			}

			fmt.Println("endpointDappList", len(endpointDappList))

			for _, endpointDapp := range endpointDappList {
				fmt.Println(endpointDapp.DetailDapp.DAppName)
				//todo: check exist to update or insert

				countFail := 0
			CallAPIDetailPage:
				url := fmt.Sprintf(`%s%s`, constant.BASE_URL_DAPPRADAR, endpointDapp.Endpoint)
				dom := utils.GetHtmlDomJsRenderByUrl(url)

				//net::ERR_INTERNET_DISCONNECTED
				if dom == nil {
					fmt.Println(log.LogLevelError, `utils.GetHtmlDomJsRenderByUrl(url)`+url, err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
				isPasss := CrawlDappDetailByHtmlDom(dom, endpointDapp.DetailDapp)
				if !isPasss {
					countFail += 1
					time.Sleep(1 * time.Second)
					if countFail == 10 {
						log.Println(log.LogLevelError, `dom fail `+url, countFail)
					} else {
						goto CallAPIDetailPage
					}
				}

				endpointDapp.DetailDapp.DAppId = "gear5_dapp_" + endpointDapp.DetailDapp.DAppCode + "_dappradar"
				endpointDapp.DetailDapp.DAppSrc = "dappradar"
				endpointDapp.DetailDapp.SourceUrl = url
				endpointDapp.DetailDapp.SourceName = "dappradar"

				err := endpointDapp.InsertDB()
				if err != nil {
					fmt.Println(log.LogLevelError, `Insert fail endpointDapp`+url, err.Error())
					continue
				}
				// fmt.Println("i", i)
			}
			log.Println(log.LogLevelInfo, fmt.Sprintf("Done crawl %d dapp chain %s page %d: %d/ %d(blockchain)", len(endpointDappList), endpointBlockChain.BlockchainId, pageIdx, index, len(endpointBlockChains)), len(MapDappradar))

		}
	}
}

func Crawl() {

	MapDappradar = make(map[string]dto_dappradar.DetailDapp)

	var listDappradar dto_dappradar.ListDappradar
	err := listDappradar.GetAllDappradar("dappradar")
	if err != nil {
		log.Println(log.LogLevelError, `listDappradar.GetAllDappradar(dappradar)`, err.Error())
		return
	}

	for _, dapp := range listDappradar.DappList {
		MapDappradar[dapp.DAppCode] = dapp
	}

	log.Println(log.LogLevelInfo, `Get All Dappradar `, len(MapDappradar))

	chainList := dao.ChainList{}

	err = chainList.GetChainList()
	if err != nil {
		log.Println(log.LogLevelError, `chainList.GetChainList()`, err.Error())
	} else {
		MapChainList = make(map[string]string, 0)

		for _, chain := range chainList.List {
			MapChainList[chain.Chainname] = chain.ChainId
		}
	}

	if len(MapChainList) == 0 {
		return
	}

	log.Println(log.LogLevelInfo, `Get chainlist succes`, len(MapChainList))

	url := fmt.Sprintf(`%s%s`, constant.BASE_URL_DAPPRADAR, constant.ENDPOINT_CHAIN_ALL_DAPPRADAR)
CallAPIBlockchain:
	dom := utils.GetHtmlDomJsRenderByUrl(url)

	if dom == nil {
		time.Sleep(5 * time.Second)
		goto CallAPIBlockchain
	}
	endpointBlockChains := CrawlEndpointBlockchainsByHtmlDom(dom)
	if len(endpointBlockChains) == 0 {
		time.Sleep(5 * time.Second)
		goto CallAPIBlockchain
	}

	fmt.Println("endpointBlockChains", len(endpointBlockChains))

	// for _, endpoint := range endpointBlockChains {
	// 	fmt.Print(endpoint.BlockchainId, endpoint.BlockchainName, endpoint.Endpoint, "===")
	// }

	dtoEndpointBlockchainRepo := dto_dappradar.EndpointBlockchainRepo{}
	dtoEndpointBlockchainRepo.EndpointBlockchains = endpointBlockChains
	daoBlockchainRepo := &dao.BlockchainRepo{}
	dtoEndpointBlockchainRepo.ConvertTo(daoBlockchainRepo)
	daoBlockchainRepo.InsertDB()

	for index, endpointBlockChain := range endpointBlockChains {
		for pageIdx := 1; ; pageIdx++ {
			url := fmt.Sprintf(`%s%v/%d`, constant.BASE_URL_DAPPRADAR, endpointBlockChain.Endpoint, pageIdx)
			// fmt.Println("url next", url)
		CallAPIListPagination:
			dom := utils.GetHtmlDomJsRenderByUrl(url)
			if dom == nil {
				log.Println(log.LogLevelDebug, `service/review/crawl/dappradar_dapp/Crawl/getDomJsLoad`, `dom get by js loading is nil`)
				pageIdx--
				time.Sleep(5 * time.Second)
				continue
			}
			if IsEndPage(dom) {
				break
			}

			succes, endpointDappList := CrawlDappEndpointByHtmlDom(dom)
			//Response without data
			if !succes {
				time.Sleep(5 * time.Second)
				goto CallAPIListPagination
			}

			// fmt.Println("endpointDappList", len(endpointDappList))

			for _, endpointDapp := range endpointDappList {
				//todo: check exist to update or insert

				countFail := 0
			CallAPIDetailPage:
				url := fmt.Sprintf(`%s%s`, constant.BASE_URL_DAPPRADAR, endpointDapp.Endpoint)
				dom := utils.GetHtmlDomJsRenderByUrl(url)

				//net::ERR_INTERNET_DISCONNECTED
				if dom == nil {
					fmt.Println(log.LogLevelError, `utils.GetHtmlDomJsRenderByUrl(url)`+url, err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
				isPasss := CrawlDappDetailByHtmlDom(dom, endpointDapp.DetailDapp)
				if !isPasss {
					countFail += 1
					time.Sleep(1 * time.Second)
					if countFail == 10 {
						log.Println(log.LogLevelError, `dom fail `+url, countFail)
					} else {
						goto CallAPIDetailPage
					}
				}

				endpointDapp.DetailDapp.DAppId = "gear5_dapp_" + endpointDapp.DetailDapp.DAppCode + "_dappradar"
				endpointDapp.DetailDapp.DAppSrc = "dappradar"
				endpointDapp.DetailDapp.SourceUrl = url
				endpointDapp.DetailDapp.SourceName = "dappradar"

				err := endpointDapp.InsertDB()
				if err != nil {
					fmt.Println(log.LogLevelError, `Insert fail endpointDapp`+url, err.Error())
					continue
				}
				// fmt.Println("i", i)
			}
			log.Println(log.LogLevelInfo, fmt.Sprintf("Done crawl %d dapp chain %s page %d: %d/ %d(blockchain)", len(endpointDappList), endpointBlockChain.BlockchainId, pageIdx, index, len(endpointBlockChains)), len(MapDappradar))

		}
	}
}

func IsEndPage(dom *goquery.Document) bool {
	isEndPage := false
	domKey := `h2`
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		_txtNotifyEndPageDappRadar := `Please change the filters to explore more`
		if s.Text() == _txtNotifyEndPageDappRadar {
			isEndPage = true
		}
	})
	return isEndPage
}

func UpdateWebsiteDapp() {

	CountUpdateWebsite = 0
	CountUpdateWebsiteFail = 0

	listDapp := dto_dappradar.ListDappradar{}
	err := listDapp.GetSourceurlAllDapp()
	if err != nil {
		fmt.Println(log.LogLevelError, `listDapp.GetSourceurlAllDapp()`, err.Error())
	}

	fmt.Println(len(listDapp.DappList))

	for index, dapp := range listDapp.DappList {
		countFail := 0
	CallWebsiteDapp:
		dom := utils.GetHtmlDomJsRenderByUrl(dapp.SourceUrl)
		if dom == nil {
			log.Println(log.LogLevelDebug, `utils.GetHtmlDomJsRenderByUrl(dapp.SourceUrl)`, `dom get by js loading is nil`)
			index--
			time.Sleep(5 * time.Second)
			continue
		}
		if IsEndPage(dom) {
			break
		}

		succes := CrawlDappWebsiteByHtmlDom(dom, &dapp)
		//Response without data
		if !succes {
			countFail += 1
			time.Sleep(1 * time.Second)
			if countFail == 10 {
				log.Println(log.LogLevelError, `dom fail `+dapp.SourceUrl, countFail)
				CountUpdateWebsiteFail += 1
				continue
			} else {
				goto CallWebsiteDapp
			}
		}
		// fmt.Println(dapp)

		err := dapp.UpdatedWebsite()
		if err != nil {
			log.Println(log.LogLevelError, `dapp.UpdatedWebsite()`+dapp.SourceUrl, err.Error())
			CountUpdateWebsiteFail += 1
		} else {
			CountUpdateWebsite += 1
		}
	}
}
