package store

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/shopspring/decimal"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/jplindgren/stock-price-monitor/pkg/stock/model"
)

var defaultSheet = "FICSheet"

// The A1 notation to select all stocks from default sheet
var a1Range = fmt.Sprintf("%s!A:C", defaultSheet)

type spreadsheetStore struct {
	srv           *sheets.Service
	driveSrv      *drive.Service
	spreadsheetid string
}

func NewSpreadsheetStore(credentialFilePath string, spreadsheetid string) spreadsheetStore {
	// Load the service account key JSON file
	data, err := os.ReadFile(credentialFilePath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Create a Google Sheets service client
	config, err := google.JWTConfigFromJSON(data, sheets.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := config.Client(context.Background())
	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	driveSrv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	return spreadsheetStore{
		srv:           srv,
		driveSrv:      driveSrv,
		spreadsheetid: spreadsheetid,
	}
}

func (s spreadsheetStore) List() ([]model.Stock, error) {
	resp, err := s.srv.Spreadsheets.Values.Get(s.spreadsheetid, a1Range).Do()
	if err != nil {
		return nil, err
	}

	stocks := []model.Stock{}
	if len(resp.Values) == 0 {
		slog.Debug("No data found.")
		return stocks, nil
	} else {
		for _, row := range resp.Values {
			slog.Debug(fmt.Sprint(row))
			stock, err := parseStock(row)
			if err != nil {
				fmt.Printf("Unable to retrieve data from sheet: %v\n", err)
				continue
			}
			stocks = append(stocks, stock)
		}
	}

	return stocks, nil
}

func (s spreadsheetStore) Add(ticker string, target decimal.Decimal) ([]model.Stock, error) {
	values := make([][]interface{}, 1)
	row := make([]interface{}, 3)

	row[0] = ticker
	row[1] = fmt.Sprintf("=GOOGLEFINANCE(\"%s\", \"%s\")", ticker, "price")
	row[2] = target
	values[0] = row
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	resp, err := s.srv.Spreadsheets.Values.Append(s.spreadsheetid, a1Range, valueRange).IncludeValuesInResponse(true).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return nil, err
	}

	stocks := []model.Stock{}
	for _, row := range resp.Updates.UpdatedData.Values {
		slog.Debug(fmt.Sprint(row))
		stock, err := parseStock(row)
		if err != nil {
			fmt.Printf("Unable to retrieve data from sheet: %v\n", err)
			continue
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

// TODO: Move to new
func (s spreadsheetStore) Create(name string) (string, error) {
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: name,
			//Locale: "pt_BR", //not supporting other locales until fix GOOGLEFINANCE parse error on other locales
		},
		Sheets: []*sheets.Sheet{
			{
				Properties: &sheets.SheetProperties{
					Title: defaultSheet,
				},
			},
		},
	}

	createResp, err := s.srv.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		slog.Debug("Unable to create spreasheet")
		return "", err
	}

	spreadsheetId := createResp.SpreadsheetId

	// Define the requests for formatting
	boldFormat := &sheets.TextFormat{
		Bold: true,
	}
	numberFormat := &sheets.CellFormat{
		NumberFormat: &sheets.NumberFormat{
			Type:    "NUMBER",
			Pattern: "0.00",
		},
	}
	boldRequest := &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				SheetId:          createResp.Sheets[0].Properties.SheetId,
				StartColumnIndex: 0,
				EndColumnIndex:   1,
			},
			Cell: &sheets.CellData{
				UserEnteredFormat: &sheets.CellFormat{
					TextFormat: boldFormat,
				},
			},
			Fields: "userEnteredFormat.textFormat",
		},
	}
	numberRequest := &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				SheetId:          createResp.Sheets[0].Properties.SheetId,
				StartColumnIndex: 2,
				EndColumnIndex:   3,
			},
			Cell: &sheets.CellData{
				UserEnteredFormat: numberFormat,
			},
			Fields: "userEnteredFormat.numberFormat",
		},
	}

	// Batch update request
	batchUpdate := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{boldRequest, numberRequest},
	}
	_, err = s.srv.Spreadsheets.BatchUpdate(spreadsheetId, batchUpdate).Do()
	if err != nil {
		slog.Debug("Unable to apply batch update")
		return "", err
	}

	slog.Debug("Spreadsheet created and formatted successfully.")

	return spreadsheetId, nil
}

func (s spreadsheetStore) GrantAccess(spreadsheetId string, email string) error {
	// Define the permission
	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer",
		EmailAddress: email,
	}

	// Grant access to the spreadsheet
	_, err := s.driveSrv.Permissions.Create(spreadsheetId, permission).Do()
	return err
}

func (s spreadsheetStore) Delete(spreadsheetId string) error {
	err := s.driveSrv.Files.Delete(spreadsheetId).Do()
	return err
}

func parseStock(entry []interface{}) (model.Stock, error) {
	ticker := entry[0].(string)
	sPrice := entry[1].(string)
	sTarget := entry[2].(string)

	price, err := decimal.NewFromString(sPrice)
	if err != nil {
		return model.Stock{}, err
	}

	target, err := decimal.NewFromString(sTarget)
	if err != nil {
		return model.Stock{}, err
	}

	return model.Stock{
		Ticker: ticker,
		Price:  price,
		Target: target,
	}, nil
}
