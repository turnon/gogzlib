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

type bookList struct {
	books               []*book
	workers, result     chan *book
	workersWg, resultWg sync.WaitGroup
}

func (list *bookList) init(n int) {
	list.books = make([]*book, 0)
	list.workers, list.result = make(chan *book), make(chan *book)
	list.workersWg.Add(n)
	list.resultWg.Add(1)

	for ; n > 0; n-- {
		go func() {
			for b := range list.workers {
				b.getBookInfo()
				list.result <- b
			}
			list.workersWg.Done()
		}()
	}

	go func() {
		for b := range list.result {
			list.books = append(list.books, b)
		}
		list.resultWg.Done()
	}()
}

func (list *bookList) process(b *book) {
	list.workers <- b
}

func (list *bookList) done() []book {
	close(list.workers)
	list.workersWg.Wait()
	close(list.result)
	list.resultWg.Wait()
	books := make([]book, 0, len(list.books))
	for _, b := range list.books {
		books = append(books, *b)
	}
	return books
}

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
