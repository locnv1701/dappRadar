package dapp_expert

import (
	"fmt"
	"io"
	"net/http"
	"review-service/pkg/log"
	"review-service/pkg/utils"
	"review-service/service/review/crawl/dappcom"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Crawler() {
	fmt.Println("start")

	fmt.Println(`====================`)

	// listUrl := crawlerListByIndexOfHTML() //LASTEST 15/01 13:00 BMP WORLD CUP

	dappAddressListUrl := dappcom.DappAddressListUrl{}

	err := dappAddressListUrl.GetListUrlDappExpert()
	if err != nil {
		log.Println(log.LogLevelError, "listUrl.GetListUrlDappExpert() ", err)

	}

	fmt.Println("len list url:", len(dappAddressListUrl.ListUrl))

	for i, url := range dappAddressListUrl.ListUrl {
		fmt.Println(i, url)
		dappAddress, err := CrawlerDetailBySplitHTML(url)
		if err != nil {
			log.Println(log.LogLevelError, "CrawlerDetailBySplitHTML(url) "+url, err)
			continue
		}
		err = dappAddress.Insert()
		if err != nil {
			log.Println(log.LogLevelError, "dappAddress.Insert() "+url, err)
			continue
		}
	}

}

func CrawlerDetailBySplitHTML(url string) (dappcom.DappAddress, error) {
	dappAddress := dappcom.DappAddress{
		Url: url,
	}

	dappCode := strings.Replace(url, "https://dapp.expert/dapp/", "", -1)

	dappCode = strings.Replace(dappCode, "en-", "", -1)

	chainCode := strings.Split(dappCode, "-")[0]

	dappCode = strings.Replace(dappCode, chainCode+"-", "", -1)

	dappAddress.DappCode = dappCode

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(log.LogLevelError, "http.NewRequest", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(log.LogLevelError, "client.Do", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(log.LogLevelError, "io.ReadAll", err)
	}
	// fmt.Println(string(body))

	stringBody := string(body)

	nameIndexStart := strings.Index(stringBody, `<div class="dappe-info__name">`) + len(`<div class="dappe-info__name">`)
	nameIndexEnd := strings.Index(stringBody, `<div class="dappe-info__name">`) + len(`<div class="dappe-info__name">`)

	fmt.Println("index", nameIndexStart, nameIndexEnd)
	for i := nameIndexStart; i < nameIndexStart+200; i++ {
		if string(stringBody[i]) == "<" {
			nameIndexEnd = i
			break
		}
	}

	fmt.Println("Name: ", stringBody[nameIndexStart:nameIndexEnd])
	dappAddress.Name = stringBody[nameIndexStart:nameIndexEnd]

	// categoryIndexStart := strings.Index(stringBody, `data-title="Category" data-values="`) + len(`data-title="Category" data-values="`)
	// categoryIndexEnd := strings.Index(stringBody, `data-title="Category" data-values="`) + len(`data-title="Category" data-values="`)

	// fmt.Println("index", categoryIndexStart, categoryIndexEnd)
	// for i := categoryIndexStart; i < categoryIndexStart+200; i++ {
	// 	if string(stringBody[i]) == "\"" {
	// 		categoryIndexEnd = i
	// 		break
	// 	}
	// }

	// fmt.Println("category: ", stringBody[categoryIndexStart:categoryIndexEnd])
	// dappAddress.Category = stringBody[categoryIndexStart:categoryIndexEnd]

	//todo: extract html dom get category

	dom := utils.GetHtmlDomJsRenderByUrl(url)
	if dom == nil {
		log.Println(log.LogLevelDebug, `utils.GetHtmlDomJsRenderByUrl(url)`, `dom get by js loading is nil`)
	}

	category, succes := ExtractDappCategoryByHtmlDom(dom)
	fmt.Println(succes)
	dappAddress.Category = category

	blockchainIndexStart := strings.Index(stringBody, `data-title="Blockchain" data-values="`) + len(`data-title="Blockchain" data-values="`)
	blockchainIndexEnd := strings.Index(stringBody, `data-title="Blockchain" data-values="`) + len(`data-title="Blockchain" data-values="`)

	fmt.Println("index", blockchainIndexStart, blockchainIndexEnd)
	for i := blockchainIndexStart; i < blockchainIndexStart+200; i++ {
		if string(stringBody[i]) == "\"" {
			blockchainIndexEnd = i
			break
		}
	}

	fmt.Println("blockchain: ", stringBody[blockchainIndexStart:blockchainIndexEnd])
	dappAddress.Chain = stringBody[blockchainIndexStart:blockchainIndexEnd]

	listContractAddress := []string{}

	arr := strings.Split(stringBody, `<span class="copyTo" data-addr="`)
	for i, ele := range arr {
		if i == 0 {
			continue
		}
		ele = ele[:70]
		address := strings.Split(ele, `">`)[0]
		listContractAddress = append(listContractAddress, address)
	}

	dappAddress.Address = strings.ToLower(strings.Join(listContractAddress, ","))

	fmt.Println(len(listContractAddress))

	fmt.Println(listContractAddress)

	return dappAddress, nil
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

func crawlerListByIndexOfHTML() (detailUrls []string) {
	detailUrls = make([]string, 0)
	limit := 10
	for idx := 0; idx < 1; idx++ { //recordsTotal:1078
		offset := (idx * 10)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(`https://dapp.expert/dapps/get-ranktable-all?start=%d&length=%ddraw=1&columns%5B0%5D%5Bdata%5D=0&columns%5B0%5D%5Bname%5D=&columns%5B0%5D%5Bsearchable%5D=true&columns%5B0%5D%5Borderable%5D=false&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B0%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B1%5D%5Bdata%5D=1&columns%5B1%5D%5Bname%5D=&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=true&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B2%5D%5Bdata%5D=2&columns%5B2%5D%5Bname%5D=&columns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=true&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B3%5D%5Bdata%5D=3&columns%5B3%5D%5Bname%5D=&columns%5B3%5D%5Bsearchable%5D=true&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B4%5D%5Bdata%5D=4&columns%5B4%5D%5Bname%5D=&columns%5B4%5D%5Bsearchable%5D=true&columns%5B4%5D%5Borderable%5D=true&columns%5B4%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B4%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B5%5D%5Bdata%5D=5&columns%5B5%5D%5Bname%5D=&columns%5B5%5D%5Bsearchable%5D=true&columns%5B5%5D%5Borderable%5D=true&columns%5B5%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B5%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B6%5D%5Bdata%5D=6&columns%5B6%5D%5Bname%5D=&columns%5B6%5D%5Bsearchable%5D=true&columns%5B6%5D%5Borderable%5D=true&columns%5B6%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B6%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B7%5D%5Bdata%5D=7&columns%5B7%5D%5Bname%5D=&columns%5B7%5D%5Bsearchable%5D=true&columns%5B7%5D%5Borderable%5D=true&columns%5B7%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B7%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B8%5D%5Bdata%5D=8&columns%5B8%5D%5Bname%5D=&columns%5B8%5D%5Bsearchable%5D=true&columns%5B8%5D%5Borderable%5D=true&columns%5B8%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B8%5D%5Bsearch%5D%5Bregex%5D=false&order%5B0%5D%5Bcolumn%5D=6&order%5B0%5D%5Bdir%5D=dsc&&search=&blockchain=&category=&period=24h&platform=&site_lang=1&activeTab=dapps&_=1673766431070`, offset, limit), nil)
		if err != nil {
			log.Println(log.LogLevelError, "http.NewRequest", err)
		}
		// req.Header.Add("Authority", `www.dapp.com`)
		// req.Header.Add("Method", `GET`)
		// req.Header.Add("Path", fmt.Sprintf(`/api/ranking/dapp/?page=%d&sort=usd_24h`, (idx+1)))
		// req.Header.Add("Scheme", `https`)
		// req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		// // req.Header.Add("Accept-Encoding", "gzip, deflate, br")
		// // req.Header.Add("Accept-Language:", `en-US,en;q=0.9,vi;q=0.8`)
		// req.Header.Add("Cache-Control", `max-age=0`)
		// req.Header.Add("Cookie", `__gads=ID=a52fbad86fff0c09-221afd34bed80088:T=1670140310:RT=1670140310:S=ALNI_Ma2ltJxDZOB5izfxurarkICTuVniA; _vwo_uuid_v2=D3A19E54EBB7D126370B41A0263CF4AC8|fdc026fee6f6f19cbebcc80a6e5fe291; _fbp=fb.1.1670140310401.1580027654; wingify_donot_track_actions=0; i18n_redirected_dapp_com=en; _ga=GA1.2.1644116417.1671277479; _wingify_pc_uuid=7a5923f034f44b53aa731b74af01fc36; wingify_push_do_not_show_notification_popup=true; csrftoken=ZSlWb9g9LFuOunEHwa8BKaqEescQ17dBRURiMxVoMatgU6BjakXxwulZSfzUZmRz; G_ENABLED_IDPS=google; __gpi=UID=00000b8980fd6df7:T=1670140310:RT=1673758034:S=ALNI_MZ8pBiM3Y6RQWouTf-Z3zeuEVtxgw; _gid=GA1.2.1955783314.1673732835; _vis_opt_s=4%7C; _vis_opt_test_cookie=1`)
		// req.Header.Add("sec-ch-ua", `Not?A_Brand";v="8", "Chromium";v="108", "Microsoft Edge";v="108`)
		// req.Header.Add("sec-ch-ua-mobile", `?0`)
		// req.Header.Add("sec-ch-ua-platform", `Windows`)
		// req.Header.Add("sec-ch-ua-mobile", `?0`)
		// // req.Header.Add("sec-fetch-dest:", `document`)
		// req.Header.Add("sec-fetch-mode", `navigate`)
		// req.Header.Add("sec-fetch-site", `none`)
		// req.Header.Add("sec-fetch-user", `?1`)
		// req.Header.Add("upgrade-insecure-requests", `1`)
		// req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Edg/108.0.1462.76`)
		// req.Header.Add("Connection", "keep-alive")
		// req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Println(log.LogLevelError, "client.Do", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(log.LogLevelError, "io.ReadAll", err)
		}
		// fmt.Println(string(body))

		arr := strings.Split(string(body), `<div class='td-in t-text'><a href='`)
		fmt.Println(`=========`, len(arr))

		fmt.Println(arr[0][:100])

		for i, ele := range arr {
			if i == 0 {
				continue
			}
			fmt.Println(i)
			// startTextDetailUrl := `<div class='td-in t-text'><a href='`
			endTextDetailUrl := `'><span class='dapp-img-mini-wrap'><img src='`
			startDetailUrl := 0 //strings.Index(ele, startTextDetailUrl) + len(startTextDetailUrl)
			endDetailURl := strings.Index(ele, endTextDetailUrl)
			detailUrl := ele[startDetailUrl:endDetailURl]
			detailUrl = strings.ReplaceAll(detailUrl, `\`, ``) // \/ to /
			fmt.Println("detailUrl: ", detailUrl)
			detailUrls = append(detailUrls, detailUrl)
		}
	}
	return detailUrls
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
