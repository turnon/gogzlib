package opac

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type book struct {
	Bookrecno, Title, ISBN, PrimaryAuthor, SecondaryAuthor string
}

const bookURL = "http://opac.gzlib.gov.cn/opac/book/"

func (b *book) getInfobook() {
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
}

func infoOtherThanTitle(b *book, s *goquery.Selection) {
	left := strings.TrimSpace(s.Find(".leftTD").Text())

	switch left {
	case "ISBN:":
		r := regexp.MustCompile(`(?s)\n\s+价格.*`)
		isbn := r.ReplaceAllString(infoValue(s), "")
		b.ISBN = b.ISBN + "/" + strings.TrimSpace(isbn)
	case "著者:", "主要著者:":
		b.PrimaryAuthor = b.PrimaryAuthor + "/" + authorValue(s)
	case "次要著者:":
		b.SecondaryAuthor = b.SecondaryAuthor + "/" + authorValue(s)
	}
}

func infoValue(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find(".rightTD").Text())
}

func authorValue(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find(".rightTD a").Text())
}
