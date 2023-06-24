package dao

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"review-service/pkg/server"
)

type DappComStat struct {
	DappId string

	CodeDappCom string

	User24h   int
	User7d    int
	Totaluser int

	Volume24h   float64
	Volume7d    float64
	TotalVolume float64

	Txs24h   int
	Txs7d    int
	TotalTxs int
}

// func (dappComStat *DappComStat) UpdateDappRadarByDappCom() error {

// 	query := `update dapp_tmp set user24h = $1, user7d = $2, totalUser = $3,
// 	 volume24h = $4, volume7d = $5, totalVolume = $6,
// 	 txs24h = $7, txs7d = $8, totalTxs = $9, updatedDate = $10
// 	 where dappId = $11`

// 	_, err := db.PSQL.Exec(query, dappComStat.User24h, dappComStat.User7d, dappComStat.Totaluser,
// 		dappComStat.Volume24h, dappComStat.Volume7d, dappComStat.TotalVolume,
// 		dappComStat.Txs24h, dappComStat.Txs7d, dappComStat.TotalTxs, utils.TimeNowStringVietNam(), dappComStat.DappId)

// 	return err

// }

func (dappComStat *DappComStat) CallUpdateDappRadarByDappCom() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	jsonBody, err := json.Marshal(dappComStat)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := urlConnectDatabase + "/dapp"
	req, err := http.NewRequest(http.MethodPatch, requestURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// fmt.Printf("client: got response!\n")
	// fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 || res.StatusCode < 200 {
		return errors.New(res.Status + " " + string(resBody))
	}
	return nil
}
