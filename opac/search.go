package opac

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const searchURL = "http://opac.gzlib.gov.cn/opac/search?" +
	"searchType=standard&isFacet=true&view=standard&searchWay=title&" +
	"rows=10&sortWay=score&sortOrder=desc&hasholding=1&" +
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

	getBookInfoWorkers, result := makeWorker(4)
	allBooks := make([]book, 0)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for b := range result {
			allBooks = append(allBooks, b)
		}
		wg.Done()
	}()

	doc.Find(".bookmeta").Each(func(i int, s *goquery.Selection) {
		bookrecno, exists := s.Attr("bookrecno")

		if !exists {
			panic("no bookrecno")
		}

		getBookInfoWorkers <- book{No: i, Bookrecno: bookrecno}
	})

	close(getBookInfoWorkers)

	wg.Wait()
	fmt.Println(allBooks)

}

func makeWorker(n int) (chan book, chan book) {
	in := make(chan book)
	out := make(chan book)
	var wg sync.WaitGroup
	wg.Add(n)

	for ; n > 0; n-- {
		go func() {
			for b := range in {
				b.getBookInfo()
				out <- b
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return in, out
}
