package crawl_dappradar_dapp

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	dto_dappradar "review-service/service/review/model/dto/dappradar"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var Start time.Time
var End time.Time
var Succes int
var Fail int
var Stage string

func Schedule() {
	var listDappradar dto_dappradar.ListDappradar

	Start = utils.TimeNowVietNam()

	err := listDappradar.GetAllDappradar("dappradar")
	if err != nil {
		log.Println(log.LogLevelError, `Schedule listDappradar.GetAllDappradar`, err.Error())
	}

	fmt.Println(len(listDappradar.DappList))

	listUrlDapp := []string{}

	for _, dapp := range listDappradar.DappList {
		path := strings.ReplaceAll(dapp.SourceUrl, "https://dappradar.com/", "")
		url := ""
		if len(dapp.Chains) > 1 {
			// fmt.Println("----------->", dapp.Chains)

			// fmt.Println(path)
			split := strings.Split(path, strings.Split(path, "/")[0])
			if len(split) > 1 {
				url = "/multichain" + split[1]
			}
		} else {
			url = "/" + path
		}
		listUrlDapp = append(listUrlDapp, url)
	}
	fmt.Println(len(listUrlDapp))

	Succes = 0
	Fail = 0

	listPathFail := []string{}
	Stage = "listUrlDapp"
	for _, dappPath := range listUrlDapp {
		timesleep := rand.Intn(2) + 3
		time.Sleep(time.Duration(timesleep) * time.Second)

		dom := utils.GetHtmlDomJsRenderByUrl(fmt.Sprintf("https://dappradar.com%s", dappPath))

		// net::ERR_INTERNET_DISCONNECTED
		if dom == nil {
			log.Println(log.LogLevelError, `Schedule dom nil`, "")
			Fail += 1
			listPathFail = append(listPathFail, dappPath)
		}

		dapp := dto_dappradar.DetailDapp{}

		CallDappDetailByHtmlDom(dom, &dapp)

		url := fmt.Sprintf("https://dappradar.com/v2/api/dapp%s/statistic/day?currency=USD", dappPath)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("authority", "dappradar.com")
		req.Header.Add("Cookie", "__cf_bm=y4UkETEnAvpt9ne1RA_Hg__gkJyFU53BjfDO3zBr03A-1673958997-0-AfWxXaQiT3lWaA2FX/88xeNub5Z3+yPz5Qbp6lJFEL6rPsITmfw0FaIBCdLAAUPWPKFAarDtyDCDYXm3Es6jBec=")
		req.Header.Add("scheme", "https")
		req.Header.Add("Cookie", "notification-session-identifier=e8c6dd4b-0164-4368-a457-c5cc12bb1ac3; _rdt_uuid=1665049402055.ca712abc-0100-4f8d-b266-39b78006bc5b; _cs_c=1; _omappvp=u65NuVRfQqLsvIfJEjBVmLCNf8O066pTSsFXtLFgitDyaPw8JAz2059MB040e8OSxo2mT6qWm4ZmH91pz57pNDOzKzGgNx9l; _ga_7R16E5X6VC=GS1.1.1669352638.4.1.1669352726.0.0.0; omSeen-faykkdsfyzcsrgeuzhxb=1672232914713; _ga_FHBCNF18GR=GS1.1.1674010844.5.0.1674010844.0.0.0; _gid=GA1.2.2047749284.1674749792; __cf_bm=RNiFd.EKEId2RZdGufYV20eJhk_3wGugExB79dMZU2c-1674804653-0-AQ9+oKONEfuPKzhYZp0zHCns1C+OsKxkRZhluWS6gQtvKKCj1KHywx7RkGm2T4BofPJdZWGR7O99n8mTVlT/cBmDO373srXE5ToyOz66Em9wI9ArQ+XHOK1wfd3lHBmmNWZsj51V70/V8rvRYZJ/IRecFGC6+z413FlKErFQfOLz+t8wPvNmeKFs3Qg7hSdctA==; _gat=1; _ga=GA1.1.1318637805.1665049402; _cs_id=20dc7e0a-61b1-aa42-d8df-61c0e125c748.1665049402.81.1674804736.1674804578.1.1699213402428; _cs_s=6.0.1.1674806536928; _ga_BTQFKMW6P9=GS1.1.1674804651.62.1.1674804757.0.0.0")

		req.Header.Add("accept", "/")
		// req.Header.Add("accept-encoding", "gzip, deflate, br")
		// req.Header.Add("accept-language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Add("sec-ch-ua", "Chromium;v=106, Google Chrome;v=106, Not;A=Brand;v=99")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sec-ch-ua-platform", "Linux")
		req.Header.Add("sec-fetch-dest", "empty")
		req.Header.Add("sec-fetch-mode", "cors")
		req.Header.Add("sec-fetch-site", "same-origin")
		req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(log.LogLevelError, `Schedule http.DefaultClient.Do(req)`, err.Error())
			Fail += 1
			listPathFail = append(listPathFail, dappPath)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(log.LogLevelError, `Schedule io.ReadAll(res.Body)`, err.Error())
			Fail += 1
			listPathFail = append(listPathFail, dappPath)
		}

		// fmt.Println("************************************************************************************************ 10")

		responseStatisticDayDapp := ResponseStatisticDayDapp{}

		err = json.Unmarshal(body, &responseStatisticDayDapp)
		if err != nil {
			log.Println(log.LogLevelError, `Schedule json.Unmarshal`, err.Error())
			Fail += 1
			listPathFail = append(listPathFail, dappPath)
			continue
		}

		if responseStatisticDayDapp.CoinName == "RADAR" {
			Fail += 1
			listPathFail = append(listPathFail, dappPath)
		} else {

			dapp := dto_dappradar.DetailDapp{
				User24h:   responseStatisticDayDapp.UserActivity,
				Volume24h: responseStatisticDayDapp.TotalVolumeInFiat,
				Balance:   responseStatisticDayDapp.TotalBalanceInFiat,
			}
			split := strings.Split(dappPath, "/")
			if len(split) > 3 {
				dapp.DAppCode = split[3]
			} else {
				Fail += 1
				listPathFail = append(listPathFail, dappPath)
				continue
			}

			err := dapp.UpdatedStatistic()
			if err != nil {
				log.Println(log.LogLevelError, `Schedule dapp.UpdatedStatistic()`, err.Error())
				Fail += 1
				listPathFail = append(listPathFail, dappPath)
			} else {
				Succes += 1
			}
		}
	}

	//todo: recall list dapp fail
	Stage = "listPathFail"
	for _, dappPath := range listPathFail {
		timesleep := rand.Intn(2) + 3
		time.Sleep(time.Duration(timesleep) * time.Second)

		dom := utils.GetHtmlDomJsRenderByUrl(fmt.Sprintf("https://dappradar.com%s", dappPath))

		// net::ERR_INTERNET_DISCONNECTED
		if dom == nil {
			log.Println(log.LogLevelError, `Schedule dom nil`, "")
			continue
		}

		dapp := dto_dappradar.DetailDapp{}

		CallDappDetailByHtmlDom(dom, &dapp)

		url := fmt.Sprintf("https://dappradar.com/v2/api/dapp%s/statistic/day?currency=USD", dappPath)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("authority", "dappradar.com")
		req.Header.Add("Cookie", "__cf_bm=y4UkETEnAvpt9ne1RA_Hg__gkJyFU53BjfDO3zBr03A-1673958997-0-AfWxXaQiT3lWaA2FX/88xeNub5Z3+yPz5Qbp6lJFEL6rPsITmfw0FaIBCdLAAUPWPKFAarDtyDCDYXm3Es6jBec=")
		req.Header.Add("scheme", "https")
		req.Header.Add("Cookie", "notification-session-identifier=e8c6dd4b-0164-4368-a457-c5cc12bb1ac3; _rdt_uuid=1665049402055.ca712abc-0100-4f8d-b266-39b78006bc5b; _cs_c=1; _omappvp=u65NuVRfQqLsvIfJEjBVmLCNf8O066pTSsFXtLFgitDyaPw8JAz2059MB040e8OSxo2mT6qWm4ZmH91pz57pNDOzKzGgNx9l; _ga_7R16E5X6VC=GS1.1.1669352638.4.1.1669352726.0.0.0; omSeen-faykkdsfyzcsrgeuzhxb=1672232914713; _ga_FHBCNF18GR=GS1.1.1674010844.5.0.1674010844.0.0.0; _gid=GA1.2.2047749284.1674749792; __cf_bm=RNiFd.EKEId2RZdGufYV20eJhk_3wGugExB79dMZU2c-1674804653-0-AQ9+oKONEfuPKzhYZp0zHCns1C+OsKxkRZhluWS6gQtvKKCj1KHywx7RkGm2T4BofPJdZWGR7O99n8mTVlT/cBmDO373srXE5ToyOz66Em9wI9ArQ+XHOK1wfd3lHBmmNWZsj51V70/V8rvRYZJ/IRecFGC6+z413FlKErFQfOLz+t8wPvNmeKFs3Qg7hSdctA==; _gat=1; _ga=GA1.1.1318637805.1665049402; _cs_id=20dc7e0a-61b1-aa42-d8df-61c0e125c748.1665049402.81.1674804736.1674804578.1.1699213402428; _cs_s=6.0.1.1674806536928; _ga_BTQFKMW6P9=GS1.1.1674804651.62.1.1674804757.0.0.0")

		req.Header.Add("accept", "/")
		// req.Header.Add("accept-encoding", "gzip, deflate, br")
		// req.Header.Add("accept-language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Add("sec-ch-ua", "Chromium;v=106, Google Chrome;v=106, Not;A=Brand;v=99")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sec-ch-ua-platform", "Linux")
		req.Header.Add("sec-fetch-dest", "empty")
		req.Header.Add("sec-fetch-mode", "cors")
		req.Header.Add("sec-fetch-site", "same-origin")
		req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(log.LogLevelError, `Schedule http.DefaultClient.Do(req)`, err.Error())
			continue
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(log.LogLevelError, `Schedule io.ReadAll(res.Body)`, err.Error())
			continue
		}

		// fmt.Println("************************************************************************************************ 10")

		responseStatisticDayDapp := ResponseStatisticDayDapp{}

		err = json.Unmarshal(body, &responseStatisticDayDapp)
		if err != nil {
			log.Println(log.LogLevelError, `Schedule json.Unmarshal`, err.Error())
			continue
		}

		if responseStatisticDayDapp.CoinName == "RADAR" {
			continue
		} else {

			dapp := dto_dappradar.DetailDapp{
				User24h:   responseStatisticDayDapp.UserActivity,
				Volume24h: responseStatisticDayDapp.TotalVolumeInFiat,
				Balance:   responseStatisticDayDapp.TotalBalanceInFiat,
			}
			split := strings.Split(dappPath, "/")
			if len(split) > 3 {
				dapp.DAppCode = split[3]
			} else {
				continue
			}

			err := dapp.UpdatedStatistic()
			if err != nil {
				log.Println(log.LogLevelError, `Schedule dapp.UpdatedStatistic()`, err.Error())
				continue
			} else {
				Succes += 1
				Fail -= 1
			}
		}
	}

	End = utils.TimeNowVietNam()
	log.Println(log.LogLevelInfo, `succes`, Succes)
	log.Println(log.LogLevelInfo, `fail`, Fail)
	fmt.Println("succes ", Succes)
	fmt.Println("fail ", Fail)

}

func CallDappDetailByHtmlDom(dom *goquery.Document, detailDapp *dto_dappradar.DetailDapp) bool {

	// fmt.Println("=====>", detailDapp.DAppName)
	// fmt.Println(dom.Html())

	//Descripttion
	domKey := `h2#dappradar-full-description`
	description := ``
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		s.Parent().Children().Each(func(i int, s *goquery.Selection) {
			unspaceStr := strings.TrimSpace(s.Text())
			if unspaceStr != `` && unspaceStr != `Back to top` {
				description += fmt.Sprintf("%s\n", s.Text())
			}
		})

	})
	detailDapp.Description = description

	return detailDapp.Description != ``

	// domString, err := dom.Html()
	// fmt.Println("err", err)
	// fmt.Println("find", strings.Contains(domString, "UAW"))

	// return strings.Contains(domString, "UAW")
}

type ResponseStatisticDayDapp struct {
	UpdatedAt                time.Time `json:"updated_at"`
	Balance                  float64   `json:"balance"`
	BalanceInFiat            float64   `json:"balanceInFiat"`
	TotalBalanceInFiat       float64   `json:"totalBalanceInFiat"`
	TransactionsChart        string    `json:"transactionsChart"`
	ActiveUsersChart         string    `json:"activeUsersChart"`
	VolumeInFiatChart        string    `json:"volumeInFiatChart"`
	BalanceFiatChart         string    `json:"balanceFiatChart"`
	ExchangeRate             float64   `json:"exchangeRate"`
	CurrencyName             string    `json:"currencyName"`
	BalanceChangeFiat        float64   `json:"balanceChangeFiat"`
	TransactionCount         float64   `json:"transactionCount"`
	TransactionChange        float64   `json:"transactionChange"`
	UserActivity             float64   `json:"userActivity"`
	UserActivityChange       float64   `json:"userActivityChange"`
	Volume                   float64   `json:"volume"`
	VolumeInFiat             float64   `json:"volumeInFiat"`
	VolumeChangeInFiat       float64   `json:"volumeChangeInFiat"`
	TotalVolumeInFiat        float64   `json:"totalVolumeInFiat"`
	TotalVolumeChangeInFiat  float64   `json:"totalVolumeChangeInFiat"`
	TotalBalanceChangeInFiat float64   `json:"totalBalanceChangeInFiat"`
	CoinName                 string    `json:"coinName"`
}
