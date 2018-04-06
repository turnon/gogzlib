package main

import (
	"fmt"
	"time"

	"github.com/turnon/gogzlib/opac"
)

func main() {
	start := time.Now()
	opac.Search("ruby")
	fmt.Println(time.Now().Sub(start))
}
