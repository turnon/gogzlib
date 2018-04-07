package main

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/turnon/gogzlib/opac"
)

func main() {

	flag.Parse()
	keywords := flag.Args()

	if len(keywords) <= 0 {
		fmt.Println("no keyword given")
		return
	}

	for i, k := range keywords {
		keywords[i] = url.QueryEscape(k)
	}

	keyword := strings.Join(keywords, "%20")

	start := time.Now()
	opac.Search(keyword)
	fmt.Println(time.Now().Sub(start))
}
