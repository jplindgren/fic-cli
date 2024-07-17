package stock

import (
	"testing"

	"github.com/jplindgren/fic-cli/pkg/stock/model"
	"github.com/shopspring/decimal"
)

type mockStore struct {
	stocks []model.Stock
}

func (m mockStore) List() ([]model.Stock, error) {
	return m.stocks, nil
}

func (m mockStore) Add(ticker string, target decimal.Decimal) ([]model.Stock, error) {
	stock := model.Stock{
		Ticker: ticker,
		Price:  decimal.NewFromFloat(39.00),
		Target: target,
	}
	m.stocks = append(m.stocks, stock)
	return m.stocks, nil
}

func (m mockStore) Create(name string) (string, error) {
	return "spreadsheetId", nil
}

func (m mockStore) GrantAccess(spreadsheetId string, email string) error {
	return nil
}

func (m mockStore) Delete(spreadsheetId string) error {
	return nil
}

func TestListShouldSortResultsBasedOnRatio(t *testing.T) {
	var tests = []struct {
		stocks []model.Stock
		want   []model.Stock
	}{
		{
			[]model.Stock{
				{Ticker: "BVMF:VALE3", Price: decimal.NewFromFloat(62.33), Target: decimal.NewFromFloat(62.33)},
				{Ticker: "BVMF:MDIA3", Price: decimal.NewFromFloat(28.14), Target: decimal.NewFromFloat(25.00)},
				{Ticker: "NYSE:KO", Price: decimal.NewFromFloat(57.55), Target: decimal.NewFromFloat(58.00)},
			},
			[]model.Stock{
				{Ticker: "NYSE:KO", Price: decimal.NewFromFloat(57.55), Target: decimal.NewFromFloat(58.00)},
				{Ticker: "BVMF:VALE3", Price: decimal.NewFromFloat(57.00), Target: decimal.NewFromFloat(62.33)},
				{Ticker: "BVMF:MDIA3", Price: decimal.NewFromFloat(28.14), Target: decimal.NewFromFloat(25.00)},
			},
		},
	}

	for _, tt := range tests {
		t.Run("correct order", func(t *testing.T) {
			mockStore := mockStore{
				stocks: tt.stocks,
			}
			srv := New(mockStore)
			ret, _ := srv.List()
			for i := range ret {
				if ret[i].Ticker != tt.want[i].Ticker {
					t.Errorf("got %s, want %s", ret[i].Ticker, tt.want[i].Ticker)
				}
			}
		})
	}
}

func TestAddParameters(t *testing.T) {
	var tests = []struct {
		name   string
		ticker string
		target string
		want   bool
	}{
		{"empty ticker", "", "39.00", false},
		{"small ticker", "A", "39.00", false},
		{"empty target", "BVMF:EGIE3", "", false},
		{"invalid target", "BVMF:EGIE3", "AA", false},
		{"negative target", "BVMF:EGIE3", "-1", false},
		{"target too big", "BVMF:EGIE3", "11000", false},
		{"valid ticker and target", "BVMF:EGIE3", "39.00", true},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok, _ := isValid(tt.ticker, tt.target)
			if ok != tt.want {
				t.Errorf("got %t, want %t", ok, tt.want)
			}
		})
	}
}
