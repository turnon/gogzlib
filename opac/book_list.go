package opac

import (
	"sort"
	"sync"
)

type bookList struct {
	books               []*book
	workers, result     chan *book
	workersWg, resultWg sync.WaitGroup
}

type booksSortByNo []book

func (bs booksSortByNo) Len() int {
	return len(bs)
}

func (bs booksSortByNo) Less(i, j int) bool {
	return bs[i].No < bs[j].No
}

func (bs booksSortByNo) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
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

	books := make(booksSortByNo, 0, len(list.books))

	for _, b := range list.books {
		books = append(books, *b)
	}

	sort.Sort(books)

	return books
}
