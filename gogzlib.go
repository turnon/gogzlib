package main

import (
	"flag"
	"fmt"
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

	keyword := strings.Join(keywords, "%20")

	start := time.Now()
	opac.Search(keyword)
	fmt.Println(time.Now().Sub(start))
}
