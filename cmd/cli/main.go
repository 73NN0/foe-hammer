package main

import (
	"os"
)

func main() {

	if err := NewCLI().Run(os.Args[1:]); err != nil {
		panic(err)
	}
}
