package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"bitbucket.org/yuce/picon"
)

const Version = "0.1.0"

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
	fmt.Printf(`       _ 
 _ __ (_) ___ ___  _ __  
| '_ \| |/ __/ _ \| '_ \ 
| |_) | | (_| (_) | | | |
| .__/|_|\___\___/|_| |_|
|_|                %s

	`, Version)
	console.Main()
	defer console.Close()
}
