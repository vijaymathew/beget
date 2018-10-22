package main

import (
	"os"
	"fmt"
	"github.com/vijaymathew/beget"
)

func main() {
	fmt.Printf("%v\n", beget.Get(os.Args[1]))
}
