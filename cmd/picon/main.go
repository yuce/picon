package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"bitbucket.com/yuce/picon"
)

func main() {
	var err error
	defaultHomeDir := ""
	usr, err := user.Current()
	if err == nil {
		defaultHomeDir = path.Join(usr.HomeDir, ".picon")
	}
	console, err := picon.NewConsole(defaultHomeDir)
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}
	console.Main()
	defer console.Close()
}
