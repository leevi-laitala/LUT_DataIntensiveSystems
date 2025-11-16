package main

import (
	"fmt"
	"strconv"
)

// Define all CLI commands and their callback functions
// Must use "CallbackFunc" wrapper
var cmds = map[string]func([]string) error {
	"exit": CallbackFunc(cmdExit),
	"help": CallbackFunc(cmdPrintHelp),

	"select": CallbackFunc(cmdSelectDb),
	"deselect": CallbackFunc(cmdDeselectDb),

	"show": CallbackFunc(cmdShow),
	"find": CallbackFunc(cmdFind),
	"insert": CallbackFunc(cmdInsertData),
	"delete": CallbackFunc(cmdDeleteData),
}

func cmdPrintHelp() {
	fmt.Println("Available commands:")
	fmt.Println("- exit : Exits CLI")
	fmt.Println("- help : Shows this help menu")
	fmt.Println("- select <server-name> : Selects server to become active")
	fmt.Println("- deselect : Deselects active server")
	fmt.Println("- show <collections | servers> : List collection or servers")
	fmt.Println("- find <collection> <?id> : Print collection contents, or specific record")
	fmt.Println("- insert <collection> : Insert new data to collection from stdin")
	fmt.Println("- delete <collection> <id> : Delete specific record from collection")
}

func cmdExit() {
	runCli = false 
}

// List all connected servers and highlight connected server
func listServers() {
	fmt.Println("Connected servers:")
	for k, v := range servers {
		fmt.Printf("- %s", k)

		if v == currentSelectedServer {
			fmt.Printf(" [selected]")
		}

		fmt.Println("")
	}

	if currentSelectedServer == nil {
		fmt.Println("No server is yet selected")
	}
}

// Select active server
func cmdSelectDb(name string) {
	s, ok := servers[name]
	if !ok {
		fmt.Printf("Server '%s' does not exist. See 'list' command\n", name)
		return
	}

	currentSelectedServer = s
	prompt = currentSelectedServer.name
}

// Deselect active server
func cmdDeselectDb() {
	currentSelectedServer = nil
	prompt = ""
}

// Show either collections or servers
func cmdShow(show string) {
	switch show {
	case "collections":
		listCollections()
	case "servers":
		listServers()
	default:
		fmt.Println("Usage: 'show <collections | servers>'")
	}
}

// View contents of collection, or specific record
func cmdFind(collection string, idStr string) {
	if collection == "" {
		fmt.Println("Usage: 'find <collection> <?id>'")
		return
	}

	id := -1

	if idStr != "" {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Must provide integer ID")
			return
		}

		if idInt < 0 {
			fmt.Println("Cannot use negative ID")
			return
		}

		id = idInt
	}

	listCollectionContents(collection, id)
}

