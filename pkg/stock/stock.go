package stock

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jplindgren/fic/pkg/stock/model"
	"github.com/jplindgren/fic/pkg/store"
	"github.com/shopspring/decimal"
)

type StockService struct {
	store store.Store
}

func (s *StockService) List() ([]model.Stock, error) {
	return s.store.List()
}

func (s *StockService) Add(ticker string, target string) ([]model.Stock, error) {
	parsedTarget, ok, err := isValid(ticker, target)
	if !ok {
		return nil, err
	}

	return s.store.Add(addDefaultStockMarket(ticker), parsedTarget)
}

func (s *StockService) CreateStore(name string) (string, error) {
	return s.store.Create(name)
}

func (s *StockService) GrantAccess(spreadsheetId string, email string) error {
	if !isValidEmail(email) {
		return errors.New("invalid email")
	}
	return s.store.GrantAccess(spreadsheetId, email)
}

func (s *StockService) DeleteStore(storeId string) error {
	return s.store.Delete(storeId)
}

func New(store store.Store) *StockService {
	return &StockService{
		store: store,
	}
}

func addDefaultStockMarket(ticker string) string {
	parts := strings.Split(ticker, ":")
	if len(parts) == 1 {
		return fmt.Sprintf("BVMF:%s", ticker)
	}
	return ticker
}

func isValid(ticker string, sTarget string) (decimal.Decimal, bool, error) {
	target, err := decimal.NewFromString(sTarget)
	if err != nil {
		slog.Debug("unable to convert target to decimal", "target", sTarget, "error", err)
		return target, false, err
	}

	aboveLimit := decimal.NewFromInt(10000)

	if target.IsNegative() || target.GreaterThan(aboveLimit) {
		slog.Debug("Target out of range: ", "target", sTarget)
		return target, false, errors.New("target must be between 0 and 100000")
	}

	if len(ticker) < 3 {
		slog.Debug("Ticker name is too short: ", "ticker", ticker)
		return target, false, errors.New("ticker must have at least 3 characters")
	}

	return target, true, nil
}
