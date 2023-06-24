package dappcom

import (
	"fmt"
	"math/rand"
	"review-service/pkg/log"
	"review-service/pkg/server"
	"review-service/pkg/utils"
	"review-service/service/review/model/dao"
	dto_dappradar "review-service/service/review/model/dto/dappradar"
	"strings"
	"time"
)

var UpdateCount int
var InsertCount int
var LastUpdated string

type Tracking struct {
	LastUpdateTime string
	Update         int
	Insert         int
}

func ScheduleUpdateDappRadarByDappCom() {

	for {
		var listDappradar dto_dappradar.ListDappradar

		LastUpdated = utils.TimeNowString()
		InsertCount = 0
		UpdateCount = 0

		err := listDappradar.CallGetCodeDappComAllDappradar()
		if err != nil {
			log.Println(log.LogLevelError, `listDappradar.GetCodeDappComAllDappradar`, err.Error())
		}
		log.Println(log.LogLevelInfo, "Start update dapp by dapp.com", len(listDappradar.DappList))

		for _, dappRadar := range listDappradar.DappList {
			// fmt.Println(i, "----------------", dappRadar.CodeDappCom)
			listCodeDappCom := strings.Split(dappRadar.CodeDappCom, ",")
			dappComUpdate := dao.DappComStat{
				DappId: dappRadar.DAppId,
			}

			user24h := 0
			user7d := 0
			totaluser := 0

			volume24h := 0.0
			volume7d := 0.0
			totalVolume := 0.0

			txs24h := 0
			txs7d := 0
			totalTxs := 0

			fail := false

			// fmt.Println("len listCode", len(listCodeDappCom))

			for _, codeDappCom := range listCodeDappCom {
				countFail := 5

				for {
					dom := utils.GetHtmlDomJsRenderByUrl(fmt.Sprintf("https://www.dapp.com/api/v3/apps/%s/related_info/", codeDappCom))
					if dom == nil {
						log.Println(log.LogLevelDebug, "https://www.dapp.com/api/v3/apps/%s/related_info/  dom nil", codeDappCom)
						countFail -= 1
						if countFail == 0 {
							fail = true
							break
						}
					} else {
						dappComStat, succes := ExtractDappComStatisticByHtmlText(dom, codeDappCom)
						// fmt.Println(succes)
						if !succes {
							countFail -= 1
							if countFail == 0 {
								log.Println(log.LogLevelDebug, "ExtractDappComStatisticByHtmlText(dom, codeDappCom)", codeDappCom)
								fail = true
								break
							}
							timesleep := rand.Intn(2) + 3
							time.Sleep(time.Duration(timesleep) * time.Second)
						} else {
							user24h += dappComStat.User24h
							user7d += dappComStat.User7d
							totaluser += dappComStat.Totaluser

							volume24h += dappComStat.Volume24h
							volume7d += dappComStat.Volume7d
							totalVolume += dappComStat.TotalVolume

							txs24h += dappComStat.Txs24h
							txs7d += dappComStat.Txs7d
							totalTxs += dappComStat.TotalTxs

							break
						}
					}
				}
			}

			dappComUpdate.User24h = user24h
			dappComUpdate.User7d = user7d
			dappComUpdate.Totaluser = totaluser

			dappComUpdate.Volume24h = volume24h
			dappComUpdate.Volume7d = volume7d
			dappComUpdate.TotalVolume = totalVolume

			dappComUpdate.Txs24h = txs24h
			dappComUpdate.Txs7d = txs7d
			dappComUpdate.TotalTxs = totalTxs
			if !fail {
				err := dappComUpdate.CallUpdateDappRadarByDappCom()
				if err != nil {
					log.Println(log.LogLevelInfo, "dappComStat.UpdateDappRadarByDappCom()", err)
				} else {
					UpdateCount += 1
				}
			}
		}

		log.Println(log.LogLevelInfo, fmt.Sprintf(`Update dapp by dapp.com done insert: %d update: %d time: %s`, InsertCount, UpdateCount, LastUpdated), len(listDappradar.DappList))

		//todo: sleep
		timeSleepSchedule := server.Config.GetInt("SLEEP_SCHEDULE")
		time.Sleep(time.Duration(timeSleepSchedule) * time.Hour)
	}
}
