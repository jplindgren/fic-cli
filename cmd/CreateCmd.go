package main

import (
	"errors"
	"fmt"
	"log/slog"
)

func createStoreCmd(app App) func([]string) error {
	return func(args []string) error {
		if len(args) < 1 {
			return errors.New("you should provide the spreasheet name")
		}

		name := args[0]

		spreadsheetId, err := app.service.CreateStore(name)
		if err != nil {
			slog.Error("Unable to create fic spreasheet ->", "error", err)
			return err
		}

		fmt.Printf("Spreasheet %s created\n", name)
		fmt.Println(spreadsheetId)

		return nil
	}
}

func grantAccessCmd(app App) func([]string) error {
	return func(args []string) error {
		if len(args) < 2 {
			return errors.New("you should provide the spreasheet id and user email")
		}

		spreadsheetId := args[0]
		email := args[1]

		err := app.service.GrantAccess(spreadsheetId, email)
		if err != nil {
			slog.Error("Unable to grant permission", "spreasheetid", spreadsheetId, "user", email, "error", err)
			return err
		}

		fmt.Println("Access granted successfully.")
		return nil
	}
}

func deleteStoreCmd(app App) func([]string) error {
	return func(args []string) error {
		if len(args) < 1 {
			return errors.New("you should provide the spreasheet id")
		}

		spreadsheetId := args[0]

		err := app.service.DeleteStore(spreadsheetId)
		if err != nil {
			slog.Error("Unable to delete the store", "spreasheetid", spreadsheetId, "error", err)
			return err
		}

		fmt.Println("Spreadsheet deleted successfully.")
		return nil
	}
}
