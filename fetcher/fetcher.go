package fetcher

import (
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/elastic/beats/libbeat/logp"
)

type SymbolQuote struct {
	Symbol string
	CompanyName string
	PrimaryExchange string
	Sector string
	CalculationPrice string
	LatestPrice float64
	LatestUpdate int64
	LatestVolume int64
}

type SymbolInfo struct {
	Quote SymbolQuote
}

type SymbolsResult map[string]SymbolInfo

func quoteURL(symbols []string) (url string) {
	return fmt.Sprintf("https://api.iextrading.com/1.0/stock/market/batch?symbols=%s&types=quote",
		strings.Join(symbols, ","))
}

func RetrieveQuotes(symbols []string) (result *[]SymbolQuote, err error) {
	url := quoteURL(symbols)
	resp, err := http.Get(url)

	if err != nil {
		return result, err
	}

	logp.Info("Fetched quote from url[%s] status[%d]", url, resp.StatusCode)

	var decoded SymbolsResult;
	err = json.NewDecoder(resp.Body).Decode(&decoded)

	if err != nil {
		return nil, err
	}

	// Remove the unnecessary nesting inside Quote.
	// We don't retrieve other types from the result

	result = &[]SymbolQuote{}
	for _, symbolInfo := range decoded {
		*result = append(*result, symbolInfo.Quote)
	}

	return result, nil
}
