package crawl_dappradar_dapp

import (
	"fmt"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	dto_dappradar "review-service/service/review/model/dto/dappradar"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CrawlDappEndpointByHtmlDom(dom *goquery.Document) (bool, []dto_dappradar.EndpointDapp) {
	// fmt.Println(dom.Html())
	succes := false

	dappradarList := make([]dto_dappradar.EndpointDapp, 0)

	var classTable string
	var classTbody string
	var classTr string
	fmt.Println(classTr)

	domKey := `section` + utils.ConvertClassesFormatFromBrowserToGoQuery(`Container__content css-gncjt6`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		stringHTML, _ := s.Html()
		startIndexOfTable := strings.Index(stringHTML, "<table class=\"") + len("<table class=\"")
		fmt.Println("start index of table", startIndexOfTable)

		endIndexOfTable := 0

		for i := startIndexOfTable; ; i++ {
			if stringHTML[i] == '"' {
				endIndexOfTable = i
				break
			}
		}
		classTbody = stringHTML[startIndexOfTable:endIndexOfTable]
		fmt.Println("class table: ", classTbody)

		//***************************
		startIndexOfTbody := strings.Index(stringHTML, "<tbody class=\"") + len("<tbody class=\"")
		fmt.Println("start index of tbody", startIndexOfTbody)

		endIndexOfTbody := 0

		for i := startIndexOfTbody; ; i++ {
			if stringHTML[i] == '"' {
				endIndexOfTbody = i
				break
			}
		}
		classTbody = stringHTML[startIndexOfTbody:endIndexOfTbody]
		fmt.Println("class tbody: ", classTbody)

		//***************************

		startIndexOfTr := strings.Index(stringHTML[endIndexOfTbody:], "<tr class=\"")
		fmt.Println("???", stringHTML[startIndexOfTr:startIndexOfTr+50])
		endIndexOfTr := 0

		for i := startIndexOfTr; ; i++ {
			if stringHTML[i] == '"' {
				endIndexOfTr = i
				break
			}
		}
		classTr = stringHTML[startIndexOfTr:endIndexOfTr]
		fmt.Println("class tr ", classTr)

	})
	// return false, nil
	domKey = `table` + utils.ConvertClassesFormatFromBrowserToGoQuery(classTable)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		fmt.Println("table")
		// fmt.Println(s.Html())
		domKey := `tbody` + utils.ConvertClassesFormatFromBrowserToGoQuery(classTbody)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			fmt.Println("tbody")
			domKey = `tr` + utils.ConvertClassesFormatFromBrowserToGoQuery(classTr)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				fmt.Println("tr--------------", i)
				fmt.Println(s.Html())

				isAd := false
				domKey = `td` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-hiDMwi hGzArv`)
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					domKey = `span`
					s.Find(domKey).Each(func(i int, s *goquery.Selection) {

						_txtAdDisplayHtml := `Ad`
						if s.Text() == _txtAdDisplayHtml {
							isAd = true
						}

					})

				})

				if !isAd {
					exist := false
					dtoEndpointDapp := dto_dappradar.EndpointDapp{}
					dtoDetailDapp := dto_dappradar.DetailDapp{}
					dtoEndpointDapp.DetailDapp = &dtoDetailDapp

					//Id, name product
					domKey = `td` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-eJDSGI cJTUWM`)
					s.Find(domKey).Each(func(i int, s *goquery.Selection) {
						succes = true
						domKey := `a`
						s.Find(domKey).Each(func(i int, s *goquery.Selection) {
							endpointDetailUrl, foundEndpointDetailUrl := s.Attr(`href`)
							if foundEndpointDetailUrl {
								fmt.Println(endpointDetailUrl)
								dtoEndpointDapp.Endpoint = endpointDetailUrl

								urlParts := strings.Split(endpointDetailUrl, `/`)

								if len(urlParts) > 3 {
									productId := urlParts[3]
									dtoDetailDapp.DAppCode = productId

									productCategory := urlParts[2]
									dtoDetailDapp.Category = productCategory

									productBlockchainId := urlParts[1]

									fmt.Println("==>", dtoDetailDapp.DAppCode, " ", dtoDetailDapp.Category, " ", productBlockchainId)

									// //todo check exist Dapp by dappCode

									// fmt.Println("check exist", dtoDetailDapp.DAppCode)

									dapp, ok := MapDappradar[dtoDetailDapp.DAppCode]
									// fmt.Println(dapp)
									// If the key exists
									if ok {
										// fmt.Println("ton tai")

										newChainName := productBlockchainId

										if productBlockchainId == "binance-smart-chain" {
											newChainName = "binance"
										}
										if productBlockchainId == "sxnetwork" {
											newChainName = "sx"
										}

										err := dapp.UpdateNewChain(newChainName, MapChainList[newChainName])
										if err != nil {
											log.Println(log.LogLevelError, "CrawlDappEndpointByHtmlDom dtocDetailDapp.UpdateNewChain(productBlockchainId)", err.Error())
										}

										MapDappradar[dtoDetailDapp.DAppCode] = dapp
										// fmt.Println(MapDappradar[dtoDetailDapp.DAppCode])

										exist = true
									} else {

										dtoDetailDapp.Chains = make(map[string]any)
										dtoDetailDapp.Chains[productBlockchainId] = MapChainList[productBlockchainId]

										MapDappradar[dtoDetailDapp.DAppCode] = dtoDetailDapp
									}

									dtoDetailDapp.Chains = make(map[string]any)

									dtoDetailDapp.Chains[productBlockchainId] = MapChainList[productBlockchainId]

								}

							}
							dtoDetailDapp.DAppName = s.Text()
						})
					})

					//Balance, UAW, Volume
					domKey = `td` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-eJDSGI ljPWzr`)
					s.Find(domKey).Each(func(i int, s *goquery.Selection) {
						if i == 0 {
							floatValue, err := HandleValueStringToFloat(s.Text())
							if err == nil {
								dtoDetailDapp.Balance = floatValue
							} else {
								log.Println(log.LogLevelError, "CrawlDappEndpointByHtmlDom error parse string value to float: "+s.Text(), err.Error())
							}
						}
						if i == 1 {
							floatValue, err := HandleValueStringToFloat(s.Text())
							if err == nil {
								dtoDetailDapp.User24h = floatValue
							} else {
								log.Println(log.LogLevelError, "CrawlDappEndpointByHtmlDom error parse string value to float: "+s.Text(), err.Error())
							}
						}
						if i == 2 {
							floatValue, err := HandleValueStringToFloat(s.Text())
							if err == nil {
								dtoDetailDapp.Volume24h = floatValue
							} else {
								log.Println(log.LogLevelError, "CrawlDappEndpointByHtmlDom error parse string value to float: "+s.Text(), err.Error())
							}
						}
					})

					//Image product
					domKey = `td` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-eJDSGI lpiXsG`)
					s.Find(domKey).Each(func(i int, s *goquery.Selection) {
						domKey := `img`
						s.Find(domKey).Each(func(i int, s *goquery.Selection) {
							imgUrl, foundImgUrl := s.Attr(`src`)
							if foundImgUrl {
								dtoDetailDapp.Image = imgUrl
							}
						})

					})

					// fmt.Println("=>", dtoDetailDapp, "<=")

					if !exist {
						dappradarList = append(dappradarList, dtoEndpointDapp)
					}
				}

			})

		})

	})
	return succes, dappradarList
}

func HandleValueStringToFloat(value string) (float64, error) {

	value = strings.ToLower(value)
	value = strings.Replace(value, " ", "", -1)

	var valueFloat float64

	if strings.Contains(value, "-") {
		value = strings.Split(value, "-"+strings.Split(value, "-")[len(strings.Split(value, "-"))-1])[0]
	}
	if strings.Contains(value, "+") {
		value = strings.Split(value, "+"+strings.Split(value, "+")[len(strings.Split(value, "+"))-1])[0]
	}
	if strings.Contains(value, "$") {
		value = strings.Split(value, "$")[1]
	}
	if strings.Contains(value, "k") {
		value = strings.Split(value, "k")[0]
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0.0, err
		} else {
			return valueFloat * 1000, nil
		}
	}
	if strings.Contains(value, "m") {
		value = strings.Split(value, "m")[0]
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0.0, err
		} else {
			return valueFloat * 1000000, nil
		}
	}
	if strings.Contains(value, "b") {
		value = strings.Split(value, "b")[0]
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0.0, err
		} else {
			return valueFloat * 1000000000, nil
		}
	}
	if valueFloat == 0 {
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0.0, err
		} else {
			return valueFloat, nil
		}
	}

	return valueFloat, nil
}
