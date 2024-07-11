package main

import (
	"log/slog"
	"os"

	"github.com/jplindgren/fic-cli/pkg/stock/model"
	"github.com/olekukonko/tablewriter"
)

func listCmd(app App) func([]string) error {
	return func(args []string) error {
		stocks, err := app.service.List()
		if err != nil {
			slog.Error("Unable to retrieve data from sheet.")
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Ticker", "Price", "condition"})

		showAll := len(args) > 0 && args[0] == "-a"
		for _, stock := range stocks {
			isRecommended, condition := stock.IsRecommended()
			if !showAll && !isRecommended {
				continue
			}

			rowCustomization := getRowCustomization(condition)
			rowValues := []string{stock.Ticker, stock.Price.String(), condition}

			table.Rich(rowValues, rowCustomization)
		}

		table.Render()
		return nil
	}
}

func getRowCustomization(condition string) []tablewriter.Colors {
	if condition == model.Excellent {
		return []tablewriter.Colors{
			{tablewriter.Normal, tablewriter.FgGreenColor},   // ticker column
			{tablewriter.Normal, tablewriter.FgGreenColor},   // price column
			{tablewriter.Normal, tablewriter.FgHiGreenColor}, // condition column
		}
	} else if condition == model.Good {
		return []tablewriter.Colors{
			{tablewriter.Normal, tablewriter.FgCyanColor},
			{tablewriter.Normal, tablewriter.FgCyanColor},
			{tablewriter.Normal, tablewriter.FgHiCyanColor},
		}
	} else if condition == model.InTarget {
		return []tablewriter.Colors{
			{tablewriter.Normal, tablewriter.FgYellowColor},
			{tablewriter.Normal, tablewriter.FgYellowColor},
			{tablewriter.Normal, tablewriter.FgHiYellowColor},
		}
	} else {
		return []tablewriter.Colors{
			{tablewriter.Normal, tablewriter.FgRedColor},
			{tablewriter.Normal, tablewriter.FgRedColor},
			{tablewriter.Normal, tablewriter.FgHiRedColor},
		}
	}
}
