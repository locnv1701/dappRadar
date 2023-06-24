package crawl_dappradar_dapp

import (
	"fmt"
	"review-service/pkg/utils"
	dto_dappradar "review-service/service/review/model/dto/dappradar"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// TODO: check data length response equal zero
func CrawlDappDetailByHtmlDom(dom *goquery.Document, detailDapp *dto_dappradar.DetailDapp) bool {

	// fmt.Println("=====>", detailDapp.DAppName)
	fmt.Println(dom.Html())

	//Descripttion
	domKey := `h2#dappradar-full-description`
	description := ``
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Html())
		s.Parent().Children().Each(func(i int, s *goquery.Selection) {
			unspaceStr := strings.TrimSpace(s.Text())
			if unspaceStr != `` && unspaceStr != `Back to top` {
				description += fmt.Sprintf("%s\n", s.Text())
			}
		})

	})
	detailDapp.Description = description
	fmt.Println(description)

	//Tags
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-kgDmvS iCijWs`)
	tags := make([]string, 0)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey := `span`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			tags = append(tags, s.Text())
		})

	})
	detailDapp.SubCategories = tags

	//Social Media, (default type: Social media)
	domKey = `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-bYdiNP jzvkRo`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `a`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			socialUrl, foundSocialUrl := s.Attr(`href`)
			if foundSocialUrl {
				if detailDapp.Social == nil {
					detailDapp.Social = make(map[string]any, 0)
				}
				detailDapp.Social[socialUrl] = socialUrl
			}
		})

	})

	return detailDapp.Description != ``
}

func CrawlDappWebsiteByHtmlDom(dom *goquery.Document, detailDapp *dto_dappradar.DetailDapp) bool {
	// fmt.Println("find website", detailDapp.SourceUrl)

	domKey := `div` + utils.ConvertClassesFormatFromBrowserToGoQuery(`sc-ekFWYn uucLc`)
	dom.Find(domKey).Each(func(i int, s *goquery.Selection) {

		domKey = `a`
		s.Find(domKey).Each(func(i int, s *goquery.Selection) {
			websiteUrl, foundWebsiteUrl := s.Attr(`href`)
			if foundWebsiteUrl {
				detailDapp.Website = websiteUrl
			}
		})
	})

	return detailDapp.Website != ``
}
