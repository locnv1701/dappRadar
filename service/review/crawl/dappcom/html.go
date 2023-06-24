package dappcom

import (
	"encoding/json"
	"fmt"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"review-service/service/review/model/dao"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractDappComStatisticByHtmlText(dom *goquery.Document, codeDappCom string) (dao.DappComStat, bool) {
	dappComStat := dao.DappComStat{
		CodeDappCom: codeDappCom,
	}
	// fmt.Println(dom.Html())
	domText := dom.Text()
	indexStart := strings.Index(domText, "extract_stats")
	// indexStart = indexStart - 7

	indexEnd := strings.Index(domText, "community_gr")
	// indexEnd = indexEnd + 36

	if indexStart > 0 && indexEnd > 0 {

		dataJson := `{"` + domText[indexStart:indexEnd] + `community_gr": 0.0	}}}`
		// fmt.Println(dataJson)

		statDappDotCom := StatDappDotCom{}
		err := json.Unmarshal([]byte(dataJson), &statDappDotCom)
		if err != nil {
			log.Println(log.LogLevelInfo, "json.Unmarshal([]byte(dataJson), &statDappDotCom)", err)
		}
		// fmt.Println(len(statDappDotCom.ExtractStats))

		var (
			user24h   int
			user7d    int
			totaluser int

			volume24h   float64
			volume7d    float64
			totalVolume float64

			txs24h   int
			txs7d    int
			totalTxs int
		)

		for _, data := range statDappDotCom.ExtractStats {
			// fmt.Println(data.Data.Total, data.Data.Day, data.Data.Week)
			if strings.EqualFold(data.Name, "users") {
				if data.Data.Total != 0 {
					totaluser = int(data.Data.Total)
				}
				if data.Data.Week != 0 {
					user7d = int(data.Data.Week)
				}
				if data.Data.Day != 0 {
					user24h = int(data.Data.Day)
				} else {
					user24h = int(data.Data.Week / 7)
				}
			}

			if strings.EqualFold(data.Name, "volume") {
				if data.Data.Total != 0 {
					totalVolume = data.Data.Total
				}
				if data.Data.Week != 0 {
					volume7d = data.Data.Week
				}
				if data.Data.Day != 0 {
					volume24h = data.Data.Day
				} else {
					volume24h = data.Data.Week / 7
				}
			}

			if strings.EqualFold(data.Name, "transactions") {
				if data.Data.Total != 0 {
					totalTxs = int(data.Data.Total)
				}
				if data.Data.Week != 0 {
					txs7d = int(data.Data.Week)
				}
				if data.Data.Day != 0 {
					txs24h = int(data.Data.Day)
				} else {
					txs24h = int(data.Data.Week / 7)
				}
			}
		}
		dappComStat.User24h = user24h
		dappComStat.User7d = user7d
		dappComStat.Totaluser = totaluser

		dappComStat.Volume24h = volume24h
		dappComStat.Volume7d = volume7d
		dappComStat.TotalVolume = totalVolume

		dappComStat.Txs24h = txs24h
		dappComStat.Txs7d = txs7d
		dappComStat.TotalTxs = totalTxs

		return dappComStat, true

	}
	return dappComStat, false

}

func ExtractDappComByHtmlDom(dom *goquery.Document) bool {

	// fmt.Println(dom.Html())

	domKey := `section` + utils.ConvertClassesFormatFromBrowserToGoQuery(`intro-card-pc`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey := `img` + utils.ConvertClassesFormatFromBrowserToGoQuery(`app-logo`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			img, foundImg := s.Attr(`src`)
			if foundImg {
				fmt.Println("img", img)
			}
		})

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`extra-info`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`chain-item first`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {

				domKey = `a`
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {
					fmt.Println("chain", strings.Trim(s.Text(), " \t\n\x0B\f\r"), "<")
				})
			})
			domKey = `a`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				href, foundHref := s.Attr(`href`)
				if foundHref {
					if strings.Contains(href, "category") {
						fmt.Println("category", strings.Trim(s.Text(), " \t\n\x0B\f\r"), "<")

					}
				}
			})
		})

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`desc`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			fmt.Println("description", strings.Trim(s.Text(), " \t\n\x0B\f\r"), "<")
		})

		domKey = `a` + utils.ConvertClassesFormatFromBrowserToGoQuery(`visit-website`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			href, foundHref := s.Attr(`href`)
			if foundHref {
				fmt.Println("website", href)
			}

		})

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`dapp-socials`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {

			domKey = `a`
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				href, foundHref := s.Attr(`href`)
				if foundHref {
					fmt.Println("socail", i, href)
				}
			})

		})

	})

	domKey = `section` + utils.ConvertClassesFormatFromBrowserToGoQuery(`detail-page-outer`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`left-sec`)
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			fmt.Println(s.Html())

			domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`stats-card`)
			s.Find(domKey).Each(func(i int, s *goquery.Selection) {
				fmt.Println(s.Html())

				domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`item`)
				s.Find(domKey).Each(func(i int, s *goquery.Selection) {

					domKey = `span` + utils.ConvertClassesFormatFromBrowserToGoQuery(`main`)
					s.Find(domKey).Each(func(i int, s *goquery.Selection) {
						fmt.Println("statistics", i, s.Text())
					})

				})

			})

		})

	})

	// domString := dom.Text()
	// fmt.Println(domString)
	// // if err != nil {
	// // 	fmt.Println("error", err)
	// // }

	// indexOf7dUsers := strings.Index(domString, "3624")

	// fmt.Println("indexOf7dUsers", indexOf7dUsers)
	// fmt.Println("ne", domString[indexOf7dUsers-100:indexOf7dUsers+100])

	return strings.Contains(dom.Text(), "Users")
}
func ExtractDappListByHtmlText(dom *goquery.Document) (ListDappDotCom, bool) {

	// fmt.Println(dom.Html())
	indexResult := strings.Index(dom.Text(), "results")
	// fmt.Println("index results", indexResult)

	indexUpdateAt := strings.Index(dom.Text(), "update_at")
	// fmt.Println("index update_at", indexUpdateAt)
	indexUpdateAt = indexUpdateAt - 7
	indexResult = indexResult + 10
	listDappDotCom := ListDappDotCom{}

	if indexUpdateAt > 0 && indexResult > 0 {
		// fmt.Println(indexResult)
		// fmt.Println(indexUpdateAt)

		dataJson := dom.Text()[indexResult:indexUpdateAt]

		// fmt.Println(dataJson)

		err := json.Unmarshal([]byte(dataJson), &listDappDotCom.List)
		if err != nil {
			fmt.Println("unmarshal", err)
		}
		fmt.Print(" ", len(listDappDotCom.List), " ")
		// fmt.Println("Name", listDappDotCom.List[0].Dapp.Name)
		// fmt.Println("volume24h", listDappDotCom.List[0].Usd24H)
		// fmt.Println("user24h", listDappDotCom.List[0].User24H)
	} else {
		fmt.Println("-------------------------------------------------------")
		return listDappDotCom, false
	}

	return listDappDotCom, true
}
