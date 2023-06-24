package main

import (
	"encoding/json"
	"fmt"
	"review-service/pkg/utils"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// for {
	// 	dom := utils.GetHtmlDomJsRenderByUrl("https://www.dapp.com/api/ranking/dapp/?chain=0&page=2&sort=usd_24h")
	// 	if dom == nil {
	// 		fmt.Println(`dom nil`)
	// 	} else {
	// 		succes := ExtractDappStatisticByHtmlText(dom)
	// 		fmt.Println(succes)
	// 		if succes {
	// 			break
	// 		} else {
	// 			timesleep := rand.Intn(2) + 3
	// 			time.Sleep(time.Duration(timesleep) * time.Second)
	// 			break
	// 		}
	// 	}
	// }

	// for {
	// 	dom := utils.GetHtmlDomJsRenderByUrl("https://www.dapp.com/app/aave-protocol")
	// 	if dom == nil {
	// 		fmt.Println(`dom nil`)
	// 	} else {
	// 		succes := ExtractDappComByHtmlDom(dom)
	// 		fmt.Println(succes)
	// 		if succes {
	// 			break
	// 		} else {
	// 			timesleep := rand.Intn(2) + 3
	// 			time.Sleep(time.Duration(timesleep) * time.Second)
	// 			break
	// 		}
	// 	}
	// }

	// for {
	// 	dom := utils.GetHtmlDomJsRenderByUrl("https://www.dapp.com/api/v3/apps/opensea/related_info/")
	// 	if dom == nil {
	// 		fmt.Println(`dom nil`)
	// 	} else {
	// 		succes := ExtractDappComStatisticByHtmlText(dom)
	// 		fmt.Println(succes)
	// 		if succes {
	// 			break
	// 		} else {
	// 			timesleep := rand.Intn(2) + 3
	// 			time.Sleep(time.Duration(timesleep) * time.Second)
	// 			break
	// 		}
	// 	}
	// }

}

func ExtractDappComStatisticByHtmlText(dom *goquery.Document) bool {

	// fmt.Println(dom.Html())
	domText := dom.Text()

	indexStart := strings.Index(domText, "extract_stats")
	indexStart = indexStart - 7

	indexEnd := strings.Index(domText, "community_gr")
	indexEnd = indexEnd + 36

	dataJson := domText[indexStart:indexEnd]
	fmt.Println(dataJson)

	statDappDotCom := StatDappDotCom{}
	err := json.Unmarshal([]byte(dataJson), &statDappDotCom)
	if err != nil {
		fmt.Println("unmarshal", err)
	}
	fmt.Println(len(statDappDotCom.ExtractStats))

	for _, data := range statDappDotCom.ExtractStats {
		fmt.Println(data.Name, data.Data.Day, data.Data.Week, data.Data.Month, data.Data.ThreeMonth)
	}

	return strings.Contains(dom.Text(), `281481481481482`)
}

type StatDappDotCom struct {
	ExtractStats []struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		HasDollarSign bool   `json:"has_dollar_sign"`
		IsCustom      bool   `json:"is_custom"`
		NoTotal       bool   `json:"no_total"`
		Data          struct {
			Day             float64 `json:"24h"`
			DayGr           float64 `json:"24h_gr"`
			Week            float64 `json:"7d"`
			WeekGr          float64 `json:"7d_gr"`
			Month           float64 `json:"30d"`
			MonthGr         float64 `json:"30d_gr"`
			Total           float64 `json:"total"`
			TotalDays       float64 `json:"total_days"`
			AllTimeHigh     float64 `json:"all_time_high"`
			AllTimeHighDate string  `json:"all_time_high_date"`
			Charts          struct {
				Labels   []string `json:"labels"`
				Datasets struct {
					Users []float64 `json:"Users"`
				} `json:"datasets"`
				Price struct {
					Name string    `json:"name"`
					Data []float64 `json:"data"`
				} `json:"price"`
			} `json:"charts"`
			ThreeMonth   float64 `json:"90d"`
			ThreeMonthGr float64 `json:"90d_gr"`
		} `json:"data"`
	} `json:"extract_stats"`
	// Ad struct {
	// 	URL         string      `json:"url"`
	// 	Image       interface{} `json:"image"`
	// 	MobileImage interface{} `json:"mobile_image"`
	// } `json:"ad"`
	// Stats struct {
	// 	Balance float64 `json:"balance"`
	// 	Token   struct {
	// 		ID         interface{} `json:"id"`
	// 		Name       string      `json:"name"`
	// 		Contract   string      `json:"contract"`
	// 		Icon       interface{} `json:"icon"`
	// 		Identifier string      `json:"identifier"`
	// 		Chain      struct {
	// 			Name         string      `json:"name"`
	// 			Icon         interface{} `json:"icon"`
	// 			IconNew      interface{} `json:"icon_new"`
	// 			ColorIcon    interface{} `json:"color_icon"`
	// 			SupportToken bool        `json:"support_token"`
	// 		} `json:"chain"`
	// 		FullName      string  `json:"full_name"`
	// 		Symbol        string  `json:"symbol"`
	// 		Price         float64 `json:"price"`
	// 		PriceGr       float64 `json:"price_gr"`
	// 		MktCap        float64 `json:"mkt_cap"`
	// 		MktCapGr      float64 `json:"mkt_cap_gr"`
	// 		OtherFiveData struct {
	// 		} `json:"other_five_data"`
	// 		Chart interface{} `json:"chart"`
	// 	} `json:"token"`
	// 	PriceAvailable bool `json:"price_available"`
	// 	Price          struct {
	// 		Name string `json:"name"`
	// 		Data []int  `json:"data"`
	// 	} `json:"price"`
	// 	RelatedTokens []interface{} `json:"related_tokens"`
	// 	Components    []struct {
	// 		Name          string `json:"name"`
	// 		Description   string `json:"description"`
	// 		HasDollarSign bool   `json:"has_dollar_sign"`
	// 		IsCustom      bool   `json:"is_custom"`
	// 		NoTotal       bool   `json:"no_total"`
	// 		Data          struct {
	// 			Two4H           float64 `json:"24h"`
	// 			Two4HGr         float64 `json:"24h_gr"`
	// 			AllTimeHigh     float64 `json:"all_time_high"`
	// 			AllTimeHighDate string  `json:"all_time_high_date"`
	// 			Charts          struct {
	// 				Labels   []string `json:"labels"`
	// 				Datasets struct {
	// 					SocialSignal []float64 `json:"Social Signal"`
	// 				} `json:"datasets"`
	// 			} `json:"charts"`
	// 		} `json:"data"`
	// 	} `json:"components"`
	// 	Community struct {
	// 		Facebook    float64 `json:"facebook"`
	// 		FacebookGr  float64 `json:"facebook_gr"`
	// 		Twitter     float64 `json:"twitter"`
	// 		TwitterGr   float64 `json:"twitter_gr"`
	// 		Telegram    float64 `json:"telegram"`
	// 		TelegramGr  float64 `json:"telegram_gr"`
	// 		Discord     float64 `json:"discord"`
	// 		DiscordGr   float64 `json:"discord_gr"`
	// 		Community   float64 `json:"community"`
	// 		CommunityGr float64 `json:"community_gr"`
	// 	} `json:"community"`
	// } `json:"stats"`
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

func ExtractDappStatisticByHtmlText(dom *goquery.Document) bool {

	fmt.Println(dom.Html())
	indexResult := strings.Index(dom.Text(), "results")
	fmt.Println("index results", indexResult)

	indexUpdateAt := strings.Index(dom.Text(), "update_at")
	fmt.Println("index update_at", indexUpdateAt)
	indexUpdateAt = indexUpdateAt - 7
	indexResult = indexResult + 10
	fmt.Println(indexResult)
	fmt.Println(indexUpdateAt)

	dataJson := dom.Text()[indexResult:indexUpdateAt]

	fmt.Println(dataJson)

	listDappDotCom := ListDappDotCom{}
	err := json.Unmarshal([]byte(dataJson), &listDappDotCom.List)
	if err != nil {
		fmt.Println("unmarshal", err)
	}
	fmt.Println(len(listDappDotCom.List))
	fmt.Println(listDappDotCom.List[0])

	return strings.Contains(dom.Text(), "TomoPool")
}

type ListDappDotCom struct {
	List []DappDotCom
}

type DappDotCom struct {
	Ranking     float64 `json:"ranking"`
	RankingGr   float64 `json:"ranking_gr"`
	Tvl         float64 `json:"tvl"`
	TvlGr       float64 `json:"tvl_gr"`
	Community   int     `json:"community"`
	CommunityGr float64 `json:"community_gr"`
	User24H     int     `json:"user_24h"`
	Amount24H   int     `json:"amount_24h"`
	Volume24H   float64 `json:"volume_24h"`
	Usd24H      float64 `json:"usd_24h"`
	User24HGr   float64 `json:"user_24h_gr"`
	Amount24HGr float64 `json:"amount_24h_gr"`
	Volume24HGr float64 `json:"volume_24h_gr"`
	Usd24HGr    float64 `json:"usd_24h_gr"`
	User7D      int     `json:"user_7d"`
	Amount7D    int     `json:"amount_7d"`
	Volume7D    float64 `json:"volume_7d"`
	Usd7D       float64 `json:"usd_7d"`
	User7DGr    float64 `json:"user_7d_gr"`
	Amount7DGr  float64 `json:"amount_7d_gr"`
	Volume7DGr  float64 `json:"volume_7d_gr"`
	Usd7DGr     float64 `json:"usd_7d_gr"`
	User30D     int     `json:"user_30d"`
	Amount30D   int     `json:"amount_30d"`
	Volume30D   float64 `json:"volume_30d"`
	Usd30D      float64 `json:"usd_30d"`
	User30DGr   float64 `json:"user_30d_gr"`
	Amount30DGr float64 `json:"amount_30d_gr"`
	Volume30DGr float64 `json:"volume_30d_gr"`
	Usd30DGr    float64 `json:"usd_30d_gr"`
	Usds24H     []struct {
		Token string  `json:"token"`
		Usd   float64 `json:"usd"`
		Ratio float64 `json:"ratio"`
	} `json:"usds_24h"`
	Usds7D []struct {
		Token string  `json:"token"`
		Usd   float64 `json:"usd"`
		Ratio float64 `json:"ratio"`
	} `json:"usds_7d"`
	Usds30D []struct {
		Token string  `json:"token"`
		Usd   float64 `json:"usd"`
		Ratio float64 `json:"ratio"`
	} `json:"usds_30d"`
	Dapp struct {
		Pk         int    `json:"pk"`
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
		Icon       string `json:"icon"`
		Category   struct {
			Pk              int         `json:"pk"`
			Name            string      `json:"name"`
			Icon            string      `json:"icon"`
			BackgroundColor interface{} `json:"background_color"`
		} `json:"category"`
		Chains []struct {
			Pk           int    `json:"pk"`
			Name         string `json:"name"`
			Icon         string `json:"icon"`
			IconNew      string `json:"icon_new"`
			ColorIcon    string `json:"color_icon"`
			SupportToken bool   `json:"support_token"`
		} `json:"chains"`
		AddTime        string  `json:"add_time"`
		Abstract       string  `json:"abstract"`
		Heat           int     `json:"heat"`
		HasToken       bool    `json:"has_token"`
		AlexaRanking   int     `json:"alexa_ranking"`
		SocialSignal   int     `json:"social_signal"`
		SocialSignalGr float64 `json:"social_signal_gr"`
		SocialTop100   bool    `json:"social_top100"`
	} `json:"dapp"`
	HasToken bool `json:"has_token"`
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

func ExtractDappCategoryByHtmlDom(dom *goquery.Document) (string, bool) {
	// fmt.Println("find website", detailDapp.SourceUrl)
	category := ""

	domKey := `span` + utils.ConvertClassesFormatFromBrowserToGoQuery(`dappe-info-name-text`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			fmt.Println(s.Text())
			category = strings.Trim(s.Text(), "\\s+\n\t ")
		}
	})

	return category, true
}
