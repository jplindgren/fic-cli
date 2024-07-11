package main

import (
	"flag"
	"log"
	"os"

	"github.com/jplindgren/fic-cli/pkg/stock"
	"github.com/jplindgren/fic-cli/pkg/store"
)

type App struct {
	service *stock.StockService
}

type cliCommand struct {
	name     string
	callback func(args []string) error
}

func getCommands(argument string, appConfig App) cliCommand {
	commands := map[string]cliCommand{
		"add": {
			name:     "add",
			callback: addCmd(appConfig),
		},
		"list": {
			name:     "list",
			callback: listCmd(appConfig),
		},
		"create": {
			name:     "create",
			callback: createStoreCmd(appConfig),
		},
		"grant-access": {
			name:     "grant-access",
			callback: grantAccessCmd(appConfig),
		},
		"delete": {
			name:     "delete",
			callback: deleteStoreCmd(appConfig),
		},
	}
	return commands[argument]
}

func main() {

	spreadsheetId := flag.String("spreadsheetid", os.Getenv("SPREADSHEET_ID"), "Google Spreadsheet ID")
	credentialPath := flag.String("credential", os.Getenv("CREDENTIAL"), "Path to credential file")

	//if *spreadsheetId == "" || *credentialPath == "" {
	if *credentialPath == "" {
		log.Fatalf("Spreadsheet Id and Credential file should be provided. Set SPREADSHEET_ID/CREDENTIAL env variables or use the flags spreadsheet-id/credential")
	}

	store := store.NewSpreadsheetStore(*credentialPath, *spreadsheetId)
	app := App{
		service: stock.New(store),
	}

	command := getCommands(os.Args[1], app)
	err := command.callback(os.Args[2:])
	if err != nil {
		log.Fatalln(err)
	}
}
