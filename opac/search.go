package opac

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const searchURL = "http://opac.gzlib.gov.cn/opac/search?" +
	"searchType=standard&isFacet=true&view=standard&searchWay=title&" +
	"rows=50&sortWay=score&sortOrder=desc&hasholding=1&" +
	"searchWay0=marc&logical0=AND&page=1&q="

// Search from http://opac.gzlib.gov.cn/opac/search
func Search(keyword string) {
	resp, err := http.Get(searchURL + keyword)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		panic(err)
	}

	doc.Find(".bookmeta").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".bookmetaTitle").Text()
		title = strings.TrimSpace(title)
		fmt.Println(title)
	})

}
