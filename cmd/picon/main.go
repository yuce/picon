package main

import (
	"fmt"
	"os"

	"bitbucket.com/yuce/picon"
)

func main() {
	var err error
	console, err := picon.NewConsole()
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}
	console.Main()
	defer console.Close()
}
