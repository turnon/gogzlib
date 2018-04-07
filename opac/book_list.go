package opac

import (
	"sync"
)

type bookList struct {
	books   []*book
	workers chan *book
	wg      sync.WaitGroup
}

func (list *bookList) init(n int) {
	list.books = make([]*book, 0)
	list.workers = make(chan *book)
	list.wg.Add(n)

	for ; n > 0; n-- {
		go func() {
			for b := range list.workers {
				b.getBookInfo()
			}
			list.wg.Done()
		}()
	}
}

func (list *bookList) process(b *book) {
	list.books = append(list.books, b)
	list.workers <- b
}

func (list *bookList) done() []book {
	close(list.workers)
	list.wg.Wait()

	books := make([]book, 0, len(list.books))

	for _, b := range list.books {
		books = append(books, *b)
	}

	return books
}
