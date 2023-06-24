package dto_dappradar

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"review-service/pkg/db"
	"review-service/pkg/server"
	"review-service/pkg/utils"

	"github.com/google/uuid"
)

type DetailDapp struct {
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

	Balance   float64
	User24h   float64
	Volume24h float64
	TotalUser int

	Website string

	CodeDappCom string

	CreatedDate string
	UpdatedDate string
}

type Chain struct {
	ChainName string
	ChainId   string
}

type ListDappradar struct {
	DappList []DetailDapp
}

func (listDappradar *ListDappradar) GetSourceurlAllDapp() error {
	query := `SELECT sourceurl from dapp where dappsrc = 'dappradar'`

	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		dapp := DetailDapp{}
		err := rows.Scan(&dapp.SourceUrl)
		if err != nil {
			return err
		}

		listDappradar.DappList = append(listDappradar.DappList, dapp)
	}

	return nil

}

func (detailDapp *DetailDapp) GetDappByDappCode(dappSrc string) error {
	query := `select id, chains from dapp where dappsrc = $1 and dappcode = $2`

	chainJSONB := []byte{}
	// err := rows.Scan(&dao.Id, &dao.BlockchainId, &dao.BlockchainName,
	// 	&blockchainInfoJSONB, &dao.CreatedDate, &dao.UpdatedDate)
	// if err != nil {
	// 	return true, err
	// }
	// json.Unmarshal(blockchainInfoJSONB, &dao.Info)
	// if err != nil {
	// 	return true, err
	// }
	// return true, nil

	err := db.PSQL.QueryRow(query, dappSrc, detailDapp.DAppCode).Scan(&detailDapp.Id, &chainJSONB)
	if err != nil {
		return err
	}

	json.Unmarshal(chainJSONB, &detailDapp.Chains)
	if err != nil {
		return err
	}
	return nil
}

// func (listDappradar *ListDappradar) GetCodeDappComAllDappradar() error {
// 	// query := `select dappId, codeDappCom from dapp_tmp where codeDappcom is not null`
// 	query := `select dappId, codeDappCom from dapp_tmp where codeDappcom is not null`

// 	rows, err := db.PSQL.Query(query)
// 	if err != nil {
// 		return err
// 	}
// 	for rows.Next() {
// 		var detailDapp DetailDapp
// 		err := rows.Scan(&detailDapp.DAppId, &detailDapp.CodeDappCom)
// 		if err != nil {
// 			return err
// 		}

// 		listDappradar.DappList = append(listDappradar.DappList, detailDapp)
// 	}

// 	return nil
// }

func (listDappradar *ListDappradar) CallGetCodeDappComAllDappradar() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	requestURL := urlConnectDatabase + "listDappCode"
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 || res.StatusCode < 200 {
		return errors.New(res.Status + " " + string(resBody))
	}

	responseListDappCode := ResponseListDappCode{}
	err = json.Unmarshal(resBody, &responseListDappCode)
	if err != nil {
		return err
	}

	err = utils.Mapping(responseListDappCode.Data, listDappradar)
	if err != nil {
		return err
	}

	return nil
}

type ResponseListDappCode struct {
	Status  bool          `json:"status"`
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Data    ListDappradar `json:"data"`
}

func (listDappradar *ListDappradar) GetAllDappradar(dappSrc string) error {

	query := `select dappCode, dappName, chains, sourceUrl from dapp where dappsrc = $1`

	rows, err := db.PSQL.Query(query, dappSrc)
	if err != nil {
		return err
	}
	for rows.Next() {
		chainJSONB := []byte{}
		var detailDapp DetailDapp
		err := rows.Scan(&detailDapp.DAppCode, &detailDapp.DAppName, &chainJSONB, &detailDapp.SourceUrl)
		if err != nil {
			return err
		}
		json.Unmarshal(chainJSONB, &detailDapp.Chains)
		if err != nil {
			return err
		}
		listDappradar.DappList = append(listDappradar.DappList, detailDapp)
	}

	return nil
}

func (detailDapp *DetailDapp) UpdateNewChain(newChainName string, newChainId string) error {

	for chainName, _ := range detailDapp.Chains {
		if chainName == newChainName {
			return nil
		}
	}

	detailDapp.Chains[newChainName] = newChainId

	chainsJSONB, err := json.Marshal(detailDapp.Chains)
	if err != nil {
		return err
	}

	query := `Update dapp set chains = $1, updateddate = $2 where dappCode = $3`

	_, err = db.PSQL.Exec(query, chainsJSONB, utils.Timestamp(), detailDapp.DAppCode)
	if err != nil {
		return err
	}

	return nil
}

func (detailDapp *DetailDapp) UpdatedWebsite() error {
	query := `Update dapp set website = $1, updateddate = $2 where sourceurl = $3 and dappsrc = 'dappradar'`

	_, err := db.PSQL.Exec(query, detailDapp.Website, utils.Timestamp(), detailDapp.SourceUrl)
	if err != nil {
		return err
	}

	return nil
}

func (detailDapp *DetailDapp) UpdatedStatistic() error {
	query := `Update dapp set volume24h = $1, user24h = $2, balance = $3, updateddate = $4 where dappCode = $5 and dappsrc = 'dappradar'`

	_, err := db.PSQL.Exec(query, detailDapp.Volume24h, detailDapp.User24h, detailDapp.Balance, utils.TimeNowStringVietNam(), detailDapp.DAppCode)
	if err != nil {
		return err
	}

	return nil
}
