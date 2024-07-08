package main

import (
	"errors"
	"fmt"
	"log/slog"
)

func addCmd(app App) func([]string) error {
	return func(args []string) error {
		if len(args) < 2 {
			return errors.New("you should provide ticker and target price")
		}

		ticker := args[0]
		target := args[1]

		stocks, err := app.service.Add(ticker, target)
		if err != nil {
			slog.Error("Unable to add target ->", "ticker", args[0], "target", args[1], "error", err)
			return err
		}

		fmt.Println("Udated Data:")
		fmt.Println(stocks)

		return nil
	}
}
