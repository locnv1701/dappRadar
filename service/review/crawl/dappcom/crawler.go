package dappcom

import (
	"fmt"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	dto_dappradar "review-service/service/review/model/dto/dappradar"
	"sort"
	"strings"
)

func CrawlDappCom() {
	listBlockchainCodeDappCom := []string{"bnb", "bsc", "eth", "polygon", "matic",
		"tron", "neo", "steem", "tomochain", "icon",
		"hive", "thundercore", "zilliqa", "kardiachain", "fure"}

	mapCodeDappCom := make(map[string][]DappDotCom)

	var listDappradar dto_dappradar.ListDappradar

	err := listDappradar.GetAllDappradar("dappradar")
	if err != nil {
		log.Println(log.LogLevelError, `Schedule listDappradar.GetAllDappradar`, err.Error())
	}

	fmt.Println(len(listDappradar.DappList))

	count := 0

	dappX := []string{}

	for page := 1; page <= 258; page++ { //page end = 258
		fmt.Print("page: ", page)
	Recall:
		dom := utils.GetHtmlDomJsRenderByUrl(fmt.Sprintf("https://www.dapp.com/api/ranking/dapp/?chain=0&page=%d&sort=usd_24h", page))
		if dom == nil {
			fmt.Println(`dom nil`)
		} else {
			listDapp, succes := ExtractDappListByHtmlText(dom)
			fmt.Println(succes)
			if !succes {
				fmt.Println("recall", page)
				goto Recall
			}

			for _, dapp := range listDapp.List {

				dappIdentifier := DappIdentifier{}

				identifier := strings.ToLower(dapp.Dapp.Identifier)
				splitIdentifier := strings.Split(identifier, "-")

				chain := splitIdentifier[len(splitIdentifier)-1]
				dappIdentifier.Chain = chain

				dappIdentifier.Identifier = identifier
				dappIdentifier.DappCode = identifier

				haveChain := false

				for _, blockchainCode := range listBlockchainCodeDappCom {
					if chain == blockchainCode {
						dappIdentifier.DappCode = strings.Join(splitIdentifier[:len(splitIdentifier)-1], "-")
						haveChain = true
						break
					}
				}
				if !haveChain {
					dappIdentifier.Chain = ""
				}

				dappIdentifier.DappName = dapp.Dapp.Name

				blockchain := []string{}
				for _, chain := range dapp.Dapp.Chains {
					blockchain = append(blockchain, chain.Name)
				}

				dappIdentifier.Blockchain = strings.Join(blockchain, " ")

				listDappCom, exist := mapCodeDappCom[dappIdentifier.DappCode]
				if exist {
					listDappComExist := listDappCom
					listDappComExist = append(listDappComExist, dapp)
					mapCodeDappCom[dappIdentifier.DappCode] = listDappComExist
				} else {
					listDappComExist := []DappDotCom{}
					listDappComExist = append(listDappComExist, dapp)
					mapCodeDappCom[dappIdentifier.DappCode] = listDappComExist
				}

				dappX = append(dappX, dapp.Dapp.Identifier)
				// }
			}
		}
	}

	for dappCode, listDappCom := range mapCodeDappCom {
		dappComUpdate := DappComUpdate{
			DappCode: dappCode,
			DappName: listDappCom[0].Dapp.Name,
		}
		volume24h := 0.0
		user24h := 0.0
		blockchain := []string{}

		for _, dappCom := range listDappCom {
			volume24h += dappCom.Usd24H
			user24h += float64(dappCom.User24H)
			chains := []string{}
			for _, chain := range dappCom.Dapp.Chains {
				chains = append(chains, chain.Name)
			}
			blockchain = append(blockchain, chains...)
		}
		dappComUpdate.Volume24h = volume24h
		dappComUpdate.User24h = user24h
		dappComUpdate.Blockchain = strings.Join(blockchain, ",")

		for _, dappRadar := range listDappradar.DappList {

			if strings.EqualFold(dappRadar.DAppCode, dappComUpdate.DappCode) || strings.EqualFold(dappRadar.DAppName, dappComUpdate.DappName) {
				count += 1

				dappComUpdate.DappCodeDappRadar = dappRadar.DAppCode
				dappComUpdate.DappNameDappRadar = dappRadar.DAppName
				dappComUpdate.Slug = dappRadar.DAppCode
				break
			}
		}

		err := dappComUpdate.Insert()
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(dappX) > 1 {
		sort.SliceStable(dappX, func(i, j int) bool {
			return dappX[i] < dappX[j]
		})
	}

	fmt.Println(strings.Join(dappX, " "))
	fmt.Println("count", count)
}

func UpdateDappRadarByDappCom() {

	listDappComUpdate := ListDappComUpdate{}
	err := listDappComUpdate.GetListDappComUpdate()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("len update", len(listDappComUpdate.List))

	for i, dappComUpdate := range listDappComUpdate.List {
		fmt.Print(i)
		err := dappComUpdate.UpdateDappRadarByDappCom()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("done")

}
