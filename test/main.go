package main

import (
	"fmt"

	"gitlab.com/ollybritton/proxyfinder"
)

func main() {
	proxyfinder.Load()
	for i := 0; i < 3000; i++ {
		fmt.Println(proxyfinder.Random().url)
	}
}
