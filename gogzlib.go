package main

import "github.com/turnon/gogzlib/opac"
import "fmt"

func main() {
	holds := opac.Get("3002404536")
	fmt.Println(holds)
}
