package opac

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type book struct {
	No                                                     int
	Bookrecno, Title, ISBN, PrimaryAuthor, SecondaryAuthor string
	holdings                                               []holding
}

const bookURL = "http://opac.gzlib.gov.cn/opac/book/"

var /* const */ isbnRegexp = regexp.MustCompile(`(?s)\n\s+价格.*`)

func (b *book) getBookInfo() {
	resp, err := http.Get(bookURL + b.Bookrecno)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		panic(err)
	}

	doc.Find("#bookInfoTable tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			b.Title = strings.TrimSpace(s.Text())
		} else {
			infoOtherThanTitle(b, s)
		}
	})

	b.checkHolding()
}

func infoOtherThanTitle(b *book, s *goquery.Selection) {
	left := strings.TrimSpace(s.Find(".leftTD").Text())

	switch left {
	case "ISBN:":
		isbn := infoValue(s)
		isbn = isbnRegexp.ReplaceAllString(isbn, "")
		b.ISBN = joinWithSlash(b.ISBN, strings.TrimSpace(isbn))
	case "著者:", "主要著者:":
		b.PrimaryAuthor = joinWithSlash(b.PrimaryAuthor, authorValue(s))
	case "次要著者:":
		b.SecondaryAuthor = joinWithSlash(b.SecondaryAuthor, authorValue(s))
	}
}

func infoValue(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find(".rightTD").Text())
}

func authorValue(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find(".rightTD a").Text())
}

func joinWithSlash(existing string, appending string) string {
	if len(existing) > 0 {
		return existing + "/ " + appending
	}
	return appending
}
