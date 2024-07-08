package store

import (
	"github.com/jplindgren/stock-price-monitor/pkg/stock/model"
	"github.com/shopspring/decimal"
)

type Store interface {
	List() ([]model.Stock, error)
	Add(ticker string, target decimal.Decimal) ([]model.Stock, error)
	Create(name string) (string, error)
	GrantAccess(spreadsheetId string, name string) error
	Delete(spreadsheetId string) error
}
