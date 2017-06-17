/*
Copyright 2017 Yuce Tekol

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived
from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH
DAMAGE.
*/

package picon

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/chzyer/readline"
	pilosa "github.com/pilosa/go-pilosa"
)

type Console struct {
	httpClient        *Client
	pilosaClient      *pilosa.Client
	index             *pilosa.Index
	prompt            *promptInfo
	lastResponse      string
	inst              *readline.Instance
	homeDirectory     string
	sessionsDirectory string
	session           []string
	sessionName       string
	schema            *pilosa.Schema
}

func NewConsole(homeDirectory string) (*Console, error) {
	sessionsDirectory := ""
	if homeDirectory != "" {
		sessionsDirectory = path.Join(homeDirectory, "sessions")
	}
	console := &Console{
		prompt:            &promptInfo{address: "(not connected)", index: "(no index)"},
		homeDirectory:     homeDirectory,
		sessionsDirectory: sessionsDirectory,
		session:           []string{},
		sessionName:       autoSessionName(),
	}
	completer := readline.NewPrefixCompleter(
		readline.PcItem(":exit"),
		readline.PcItem(":connect", readline.PcItemDynamic(console.listConnections())),
		readline.PcItem(":use", readline.PcItemDynamic(console.listIndexes())),
		readline.PcItem(":ensure",
			readline.PcItem("index"),
			readline.PcItem("frame")),
		readline.PcItem(":create",
			readline.PcItem("index"),
			readline.PcItem("frame")),
		readline.PcItem(":delete",
			readline.PcItem("index"),
			readline.PcItem("frame")),
		readline.PcItem(":schema"),
		readline.PcItem(":save"),
		readline.PcItem(":session"),
	)
	config := &readline.Config{
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         ":exit",
		HistorySearchFold: true,
	}
	if homeDirectory != "" {
		config.HistoryFile = path.Join(homeDirectory, "history")
	}
	inst, err := readline.NewEx(config)
	if err != nil {
		return nil, err
	}
	console.inst = inst
	return console, nil
}

func (c *Console) Close() {
	c.inst.Close()
}

func (c *Console) Main() {
	c.ensureHomeDirectoryExists()
	log.SetOutput(c.inst.Stderr())
	c.updatePrompt()
	lines := []string{}
	for {
		line, err := c.inst.Readline()
		if err == io.EOF {
			goto exit
		}
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				fmt.Println("Type :exit to exit")
			} else {
				continue
			}
		}
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "\\") {
			c.inst.SetPrompt(">>> ")
			lines = append(lines, strings.TrimRight(line, "\\"))
			continue
		}
		if len(lines) > 0 {
			lines = append(lines, line)
			line = strings.Join(lines, "\n")
			lines = []string{}
			c.updatePrompt()
		}
		switch {
		case line == "":
			continue
		case line == ":exit":
			goto exit
		case strings.HasPrefix(line, "#"):
			c.inst.Operation.SetBuffer("# ")
		case strings.HasPrefix(line, ":"):
			err = c.executeCommand(line)
		case line == "_":
			if c.lastResponse != "" {
				fmt.Println(c.lastResponse)
			}
		default:
			err = c.executeQuery(line)
		}

		if err != nil {
			printError(err)
		} else {
			// do not save session commands to the session
			for _, s := range []string{":save", ":load", ":session"} {
				if strings.HasPrefix(line, s) {
					goto session_exit
				}
			}
			c.session = append(c.session, line)
		}
	session_exit:
	}
exit:
}

func (c *Console) listIndexes() func(string) []string {
	return func(line string) []string {
		indexNames := []string{}
		if c.schema != nil {
			for _, index := range c.schema.Indexes {
				indexNames = append(indexNames, index.Name)
			}
		}
		return indexNames
	}
}

func (c *Console) listConnections() func(string) []string {
	return func(line string) []string {
		addrs := []string{}
		file, err := os.Open(c.inst.Config.HistoryFile)
		if err != nil {
			return addrs
		}
		defer file.Close()
		addrSet := map[string]bool{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			if strings.HasPrefix(text, ":connect ") {
				fields := strings.Fields(text)
				if len(fields) > 1 {
					addrSet[fields[1]] = true
				}
			}
		}
		for addr := range addrSet {
			addrs = append(addrs, addr)
		}
		sort.Strings(addrs)
		return addrs
	}
}

func (c *Console) executeCommand(line string) (err error) {
	args := strings.Fields(line)
	cmd := args[0]
	switch cmd {
	case ":connect":
		err = c.executeConnectCommand(cmd, args[1:])
	case ":use":
		err = c.executeUseCommand(cmd, args[1:])
	case ":ensure":
		err = c.executeCreateOrEnsureCommand(cmd, args[1:])
	case ":create":
		err = c.executeCreateOrEnsureCommand(cmd, args[1:])
	case ":delete":
		err = c.executeDeleteCommand(cmd, args[1:])
	case ":save":
		err = c.executeSaveCommand(cmd, args[1:])
	case ":session":
		err = c.executeSessionCommand(cmd, args[1:])
	case ":schema":
		err = c.executeSchemaCommand(cmd, args[1:])
	default:
		err = fmt.Errorf("Invalid command: %s", cmd)
	}
	return err
}

func (c *Console) executeConnectCommand(cmd string, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: :connect pilosa-address")
	}
	uri, err := pilosa.NewURIFromAddress(args[0])
	if err != nil {
		return err
	}
	c.pilosaClient = pilosa.NewClientWithURI(uri)
	c.httpClient, _ = NewClient(uri.Normalize())
	err = c.updateSchema()
	if err != nil {
		c.pilosaClient = nil
		c.httpClient = nil
		return err
	}
	c.prompt.address = uri.Normalize()
	c.updatePrompt()
	return nil
}

func (c *Console) executeUseCommand(cmd string, args []string) (err error) {
	if len(args) != 1 {
		return errors.New("usage: :use index-name")
	}
	if c.pilosaClient == nil {
		return errNotConnected
	}
	indexName := args[0]
	c.index, err = pilosa.NewIndex(indexName, nil)
	if err != nil {
		return err
	}
	c.prompt.index = indexName
	c.updatePrompt()
	return nil
}

func (c *Console) executeCreateOrEnsureCommand(cmd string, args []string) (err error) {
	if c.pilosaClient == nil {
		return errNotConnected
	}
	if len(args) < 2 {
		return fmt.Errorf("Usage: %s {index | frame} name [option1=value1, ...]", cmd)
	}

	rawOptions, err := parseOptions(args[2:])
	if err != nil {
		return err
	}

	what := args[1]
	which := args[0]
	switch which {
	case "index":
		options, err := makeIndexOptions(rawOptions)
		if err != nil {
			return err
		}
		c.index, err = pilosa.NewIndex(what, options)
		if err != nil {
			return err
		}
		switch cmd {
		case ":create":
			err = c.pilosaClient.CreateIndex(c.index)
		case ":ensure":
			err = c.pilosaClient.EnsureIndex(c.index)
		default:
			return fmt.Errorf("Invalid command in this context: %s", cmd)
		}
		if err != nil {
			return err
		}
		c.prompt.index = what
		c.updatePrompt()
		return c.updateSchema()
	case "frame":
		if c.index == nil {
			return errNoIndex
		}
		options, err := makeFrameOptions(rawOptions)
		if err != nil {
			return err
		}
		frame, err := c.index.Frame(what, options)
		if err != nil {
			return err
		}
		switch cmd {
		case ":create":
			err = c.pilosaClient.CreateFrame(frame)
		case ":ensure":
			err = c.pilosaClient.EnsureFrame(frame)
		default:
			return fmt.Errorf("Invalid command in this context: %s", cmd)
		}
		return err
	default:
		return fmt.Errorf("Don't know how to ensure %s", which)
	}
}

func (c *Console) executeDeleteCommand(cmd string, args []string) (err error) {
	if c.pilosaClient == nil {
		return errNotConnected
	}
	if len(args) < 2 {
		return errors.New("Usage: :delete {index | frame} name1, ...")
	}

	which := args[0]
	switch which {
	case "index":
		for _, what := range args[1:] {
			c.index, err = pilosa.NewIndex(what, nil)
			if err != nil {
				printWarning(fmt.Sprintf("Skipping invalid index `%s`: %s", what, err))
				continue
			}
			err = c.pilosaClient.DeleteIndex(c.index)
			if err != nil {
				printError(fmt.Errorf("Error deleting index `%s`: %s", what, err))
				continue
			}
		}
		err = c.updateSchema()
	case "frame":
		if c.index == nil {
			return errNoIndex
		}
		for _, what := range args[1:] {
			frame, err := c.index.Frame(what, nil)
			if err != nil {
				printWarning(fmt.Sprintf("Skipping invalid index `%s`: %s", what, err))
				continue
			}
			err = c.pilosaClient.DeleteFrame(frame)
			if err != nil {
				printError(fmt.Errorf("Error deleting frame `%s`: %s", what, err))
				continue
			}
		}
		err = nil
	default:
		return fmt.Errorf("Don't know how to delete %s", which)
	}
	return err
}

func (c *Console) executeSaveCommand(cmd string, args []string) error {
	if len(args) > 0 {
		return errors.New("Usage: :save")
	}
	if c.sessionsDirectory == "" {
		return errors.New("session directory was not set")
	}
	path := path.Join(c.sessionsDirectory, c.sessionName)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, line := range c.session {
		_, err := f.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Console) executeSessionCommand(cmd string, args []string) error {
	switch len(args) {
	case 0:
		c.sessionName = autoSessionName()
	case 1:
		c.sessionName = args[1]
	default:
		return errors.New("Usage: :session [session name]")
	}
	// reset session
	c.session = []string{}
	return nil
}

func (c *Console) executeSchemaCommand(cmd string, args []string) error {
	// TODO: check number of args
	// 0: schema for all
	// 1: schema for the given index
	if len(args) > 1 {
		return errors.New("usage: :schema [index name | *]")
	}
	indexName := ""
	if len(args) == 1 {
		if args[0] != "*" {
			indexName = args[0]
		}
	} else {
		if c.index != nil {
			indexName = c.index.Name()
		}
	}
	err := c.updateSchema()
	if err != nil {
		return err
	}
	for _, index := range c.schema.Indexes {
		if indexName == "" || indexName == index.Name {
			frameList := make([]string, 0, len(index.Frames))
			for _, frame := range index.Frames {
				frameList = append(frameList, frame.Name)
			}
			if indexName != "" {
				fmt.Printf("[%s]\n", strings.Join(frameList, ", "))
			} else {
				fmt.Printf("%s [%s]\n", index.Name, strings.Join(frameList, ", "))
			}
		}
	}
	return nil
}

func (c *Console) executeQuery(line string) error {
	if c.httpClient == nil {
		return errNotConnected
	}
	if c.index == nil {
		return errNoIndex
	}
	response, err := c.httpClient.query(c.index.Name(), line)
	if err != nil {
		return err
	}
	c.lastResponse = string(response)
	fmt.Println(c.lastResponse)
	return nil
}

func (c *Console) updatePrompt() {
	c.inst.SetPrompt(fmt.Sprintf("\033[36m%s\033[0m/\033[1m\033[32m%s\033[37m>\033[0m ",
		c.prompt.address, c.prompt.index))
}

func (c *Console) ensureHomeDirectoryExists() {
	if c.homeDirectory != "" {
		err := os.MkdirAll(c.homeDirectory, 0700)
		if err != nil {
			c.homeDirectory = ""
			printWarning(fmt.Sprintf("Cannot create %s, unsetting it.", c.homeDirectory))
			return
		}
	}
	if c.sessionsDirectory != "" {
		err := os.MkdirAll(c.sessionsDirectory, 0700)
		if err != nil {
			c.sessionsDirectory = ""
			printWarning(fmt.Sprintf("Cannot create %s, unsetting it.", c.sessionsDirectory))
		}
	}
}

func (c *Console) updateSchema() error {
	if c.pilosaClient == nil {
		return errNotConnected
	}
	schema, err := c.pilosaClient.Schema()
	if err != nil {
		return err
	}
	c.schema = schema
	return nil
}

func parseOptions(strOptions []string) (options map[string]string, err error) {
	options = make(map[string]string, 0)
	for _, stropt := range strOptions {
		parts := strings.SplitN(stropt, "=", 2)
		options[parts[0]] = parts[1]
	}
	return
}

func makeIndexOptions(rawOptions map[string]string) (*pilosa.IndexOptions, error) {
	opts := &pilosa.IndexOptions{}
	for k, v := range rawOptions {
		switch k {
		case "column_label", "columnLabel", "col", "c":
			opts.ColumnLabel = v
		case "time_quantum", "timeQuantum", "time", "t":
			opts.TimeQuantum = pilosa.TimeQuantum(v)
		default:
			return nil, fmt.Errorf("Invalid index option: %s", k)
		}
	}
	return opts, nil
}

func makeFrameOptions(rawOptions map[string]string) (*pilosa.FrameOptions, error) {
	opts := &pilosa.FrameOptions{}
	for k, v := range rawOptions {
		switch k {
		case "row_label", "rowLabel", "row", "r":
			opts.RowLabel = v
		case "time_quantum", "timeQuantum", "time", "t":
			opts.TimeQuantum = pilosa.TimeQuantum(v)
		case "inverse_enabled", "inverseEnabled", "inverse", "i":
			b, err := parseBool(v)
			if err != nil {
				return nil, err
			}
			opts.InverseEnabled = b
		default:
			return nil, fmt.Errorf("Invalid index option: %s", k)
		}
	}
	return opts, nil
}

func parseBool(v string) (bool, error) {
	switch v {
	case "true", "t", "1":
		return true, nil
	case "false", "f", "0":
		return false, nil
	}
	return false, fmt.Errorf("Invalid boolean value: %s. Try one of true/t/1/ or false/f/0", v)

}
