package opac

import "sync"

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
