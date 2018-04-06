package opac

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const searchURL = "http://opac.gzlib.gov.cn/opac/search?" +
	"searchWay=title&searchSource=reader&booktype=&scWay=dim&" +
	"marcformat=&sortWay=score&sortOrder=desc&startPubdate=&" +
	"endPubdate=&rows=10&hasholding=1&curlibcode=GT&q="

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

	list := bookList{}
	list.init(2)

	doc.Find(".bookmeta").Each(func(i int, s *goquery.Selection) {
		bookrecno, exists := s.Attr("bookrecno")

		if !exists {
			panic("no bookrecno")
		}

		list.process(&book{No: i, Bookrecno: bookrecno})
	})

	books := list.done()

	fmt.Println(books)
}
