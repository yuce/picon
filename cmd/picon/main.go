package main

import (
	"log"
	"strings"

	"fmt"

	"github.com/chzyer/readline"
	pilosa "github.com/pilosa/go-client-pilosa"
)

type promptInfo struct {
	address  string
	database string
}

var client *pilosa.Client
var database *pilosa.Database
var prompt promptInfo

func listDatabases() func(string) []string {
	return func(line string) []string {
		return []string{"sample-db", "foo"}
	}
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem(":exit"),
	readline.PcItem(":connect"),
	readline.PcItem(":use", readline.PcItemDynamic(listDatabases())),
	readline.PcItem(":create",
		readline.PcItem("db"),
		readline.PcItem("frame")),
	readline.PcItem(":ensure",
		readline.PcItem("db"),
		readline.PcItem("frame")),
	readline.PcItem(":schema"),
)

var inst *readline.Instance

func main() {
	var err error
	inst, err = readline.NewEx(&readline.Config{
		HistoryFile:       "/tmp/picon.tmp",
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         ":exit",
		HistorySearchFold: true,
	})
	if err != nil {
		panic(err)
	}
	defer inst.Close()
	log.SetOutput(inst.Stderr())
	updatePrompt()
	for {
		line, err := inst.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		}
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, ":"):
			executeCommand(line)
		case line == ":exit":
			goto exit
		case line == "":
		default:
			executeQuery(line)
		}
	}
exit:
}

func executeCommand(line string) {
	var err error
	parts := strings.Split(line, " ")
	// TODO: trim
	switch parts[0] {
	case ":connect":
		uri, err := pilosa.NewURIFromAddress(parts[1])
		if err != nil {
			fmt.Println("Invalid address: ", parts[1])
			return
		}
		prompt.address = uri.GetNormalizedAddress()
		client = pilosa.NewClientWithAddress(uri)
		updatePrompt()
	case ":use":
		if client == nil {
			fmt.Println("You must first connect to a server")
			return
		}
		databaseName := parts[1]
		database, err = pilosa.NewDatabase(databaseName)
		if err != nil {
			fmt.Println("Invalid database name: ", databaseName)
			return
		}
		prompt.database = databaseName
		updatePrompt()
	case ":ensure":
		if len(parts) != 3 {
			fmt.Println("Usage: :ensure db/frame what")
			return
		}
		which := parts[1]
		what := parts[2]
		switch which {
		case "db":
			if client == nil {
				fmt.Println("You must first connect to a server")
				return
			}
			databaseName := what
			database, err = pilosa.NewDatabase(what)
			if err != nil {
				fmt.Println("Invalid database name: ", databaseName)
				return
			}
			err = client.EnsureDatabaseExists(database)
			if err != nil {
				fmt.Println("Error ensuring database: ", err)
				return
			}
			prompt.database = databaseName
			updatePrompt()
		case "frame":
			if client == nil || database == nil {
				fmt.Println("You must first connect to a server and use a database")
				return
			}
			frameName := what
			frame, err := database.Frame(frameName)
			if err != nil {
				fmt.Println("Invalid frame name: ", frameName)
				return
			}
			if err != nil {
				fmt.Println("Error ensuring frame: ", err)
				return
			}
			err = client.EnsureFrameExists(frame)
		default:
			fmt.Println("Don't know how to ensure ", which)
		}
	default:
		fmt.Println("Invalid command: ", parts[0])
	}
}
func executeQuery(line string) {
	if client == nil || database == nil {
		fmt.Println("You must first connect to a server and use a database")
		return
	}
	response, err := client.Query(database, line)
	if err != nil {
		fmt.Println("Error executing query:", err)
	}
	fmt.Println(response)
}

func updatePrompt() {
	inst.SetPrompt(fmt.Sprintf("\033[36m%s\033[0m/\033[32m%s\033[0m» ", prompt.address, prompt.database))
}
