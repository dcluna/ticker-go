package main

import (
	"testing"
	"net/http"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveSymbols(t *testing.T) {
	// Start our recorder
	r, err := recorder.New("fixtures/yahoo-finance-basic")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	// Create an HTTP client and inject our transport
	client := &http.Client{
		Transport: r, // Inject as transport!
	}

	body, _ := retrieveYFInfo(( []string {"AAPL", "GOOG"} ), client)

	assert.Equal(
		t,
		string(body[:]),
		`{"quoteResponse":{"result":[{"language":"en-US","region":"US","quoteType":"EQUITY","quoteSourceName":"Delayed Quote","triggerable":true,"postMarketChangePercent":-0.38930705,"postMarketTime":1601683200,"postMarketPrice":112.58,"postMarketChange":-0.4399948,"regularMarketChange":-3.7700043,"regularMarketChangePercent":-3.2280195,"regularMarketTime":1601668801,"regularMarketPrice":113.02,"regularMarketPreviousClose":116.79,"fullExchangeName":"NasdaqGS","sourceInterval":15,"exchangeDataDelayedBy":0,"tradeable":false,"marketState":"CLOSED","exchange":"NMS","exchangeTimezoneName":"America/New_York","exchangeTimezoneShortName":"EDT","gmtOffSetMilliseconds":-14400000,"market":"us_market","esgPopulated":false,"firstTradeDateMilliseconds":345479400000,"priceHint":2,"symbol":"AAPL"},{"language":"en-US","region":"US","quoteType":"EQUITY","quoteSourceName":"Delayed Quote","triggerable":true,"postMarketChangePercent":-0.30375704,"postMarketTime":1601678790,"postMarketPrice":1453.99,"postMarketChange":-4.4300537,"regularMarketChange":-31.669922,"regularMarketChangePercent":-2.1253698,"regularMarketTime":1601668801,"regularMarketPrice":1458.42,"regularMarketPreviousClose":1490.09,"fullExchangeName":"NasdaqGS","sourceInterval":15,"exchangeDataDelayedBy":0,"tradeable":false,"marketState":"CLOSED","exchange":"NMS","exchangeTimezoneName":"America/New_York","exchangeTimezoneShortName":"EDT","gmtOffSetMilliseconds":-14400000,"market":"us_market","esgPopulated":false,"firstTradeDateMilliseconds":1092922200000,"priceHint":2,"symbol":"GOOG"}],"error":null}}`,
	)
}
