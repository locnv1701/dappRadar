package dappcom

import (
	"fmt"
	"review-service/pkg/db"
	"review-service/pkg/utils"

	"github.com/google/uuid"
)

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

type DetailDappCom struct {
	Id       uuid.UUID
	DAppId   string
	DAppSrc  string
	DAppCode string
	DAppName string

	Category      string
	SubCategories []string //tags
	Image         string
	Description   string
	Social        map[string]any
	Chains        map[string]any

	SourceUrl  string
	SourceName string

	User24h   float64
	Volume24h float64
	Txs24h    float64

	Website string

	CreatedDate string
	UpdatedDate string
}

type DappAddress struct {
	Id          uuid.UUID
	Url         string
	DappCode    string
	Name        string
	Address     string
	Category    string
	Chain       string
	CreatedDate string
	UpdatedDate string
}

type DappAddressListUrl struct {
	ListUrl []string
}

type ListDappComUpdate struct {
	List []DappComUpdate
}

type DappComUpdate struct {
	Id         uuid.UUID
	DappCode   string
	DappName   string
	Blockchain string
	Volume24h  float64
	User24h    float64

	DappCodeDappRadar string
	DappNameDappRadar string

	Slug string
}

type DappIdentifier struct {
	Id         uuid.UUID
	Identifier string
	DappCode   string
	Chain      string
	DappName   string
	Blockchain string
}

func (dappAddressListUrl *DappAddressListUrl) GetListUrlDappExpert() error {
	query := `select url from dapp_address`

	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		err := rows.Scan(&url)
		if err != nil {
			continue
		}
		dappAddressListUrl.ListUrl = append(dappAddressListUrl.ListUrl, url)
	}
	return err
}

func (dappAddress *DappAddress) Insert() error {
	query := `INSERT INTO public.dapp_addresses
	(url, dappCode, "name", address, category, "chain", createddate, updateddate)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8);
	`

	_, err := db.PSQL.Exec(query, dappAddress.Url, dappAddress.DappCode, dappAddress.Name, dappAddress.Address, dappAddress.Category, dappAddress.Chain, utils.TimeNowStringVietNam(), utils.TimeNowStringVietNam())
	return err
}

func (d *DappComUpdate) UpdateDappRadarByDappCom() error {
	query := `update dapp set codeDappCom = $1, blockchainDappCom = $2, 
	volume24hDappCom = $3, user24hDappCom = $4 where dappCode = $5 and dappsrc = 'dappradar';`

	count, err := db.PSQL.Exec(query, d.DappCode, d.Blockchain, d.Volume24h, d.User24h,
		d.Slug)
	c, _ := count.RowsAffected()
	if c != 1 {
		fmt.Println(d.DappCode)
	}
	return err
}

func (list *ListDappComUpdate) GetListDappComUpdate() error {
	query := `select dappcode, dappname, blockchain, volume24h, user24h,
	 dappCodeDappRadar, dappNameDappRadar, slug from dapp_com_data where slug != ''`

	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		dappComUpdate := DappComUpdate{}
		err := rows.Scan(&dappComUpdate.DappCode, &dappComUpdate.DappName, &dappComUpdate.Blockchain,
			&dappComUpdate.Volume24h, &dappComUpdate.User24h, &dappComUpdate.DappCodeDappRadar, &dappComUpdate.DappNameDappRadar, &dappComUpdate.Slug)
		if err != nil {
			return err
		}
		list.List = append(list.List, dappComUpdate)
	}
	return err
}

func (d *ListDappComUpdate) GetListDappComSlug() error {
	query := `select dappcode, slug from dapp_com_data where slug != ''`

	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		dappComUpdate := DappComUpdate{}
		err := rows.Scan(dappComUpdate.DappCode, dappComUpdate.Slug)
		if err != nil {
			return err
		}

	}
	return err
}

func (d *DappComUpdate) Insert() error {
	query := `INSERT INTO public.dapp_com_data
			(dappcode, dappname, blockchain, volume24h, user24h, dappCodeDappRadar, dappNameDappRadar, slug)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := db.PSQL.Exec(query, d.DappCode, d.DappName, d.Blockchain, d.Volume24h, d.User24h,
		d.DappCodeDappRadar, d.DappNameDappRadar, d.Slug)
	return err
}

func (d *DappIdentifier) Insert() error {
	query := `INSERT INTO public.dapp_com_test
			(identifier, dappcode, "chain", dappname, blockchain)
			VALUES($1, $2, $3, $4, $5);`

	_, err := db.PSQL.Exec(query, d.Identifier, d.DappCode, d.Chain, d.DappName, d.Blockchain)
	return err
}
