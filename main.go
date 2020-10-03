package main

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var requestFields = []string{"symbol", "marketState", "regularMarketPrice", "regularMarketChange", "regularMarketChangePercent", "preMarketPrice", "preMarketChange", "preMarketChangePercent", "postMarketPrice", "postMarketChange", "postMarketChangePercent"}

const apiEndpoint = "https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com"

type financeMessage struct {
	Symbol      string `json:"symbol"`
	MarketState string `json:"marketState"`

	PreMarketPrice         float32 `json:"preMarketPrice"`
	PreMarketChange        float32 `json:"preMarketChange"`
	PreMarketChangePercent float32 `json:"preMarketChangePercent"`

	PostMarketPrice         float32 `json:"postMarketPrice"`
	PostMarketChange        float32 `json:"postMarketChange"`
	PostMarketChangePercent float32 `json:"postMarketChangePercent"`

	RegularMarketPrice         float32 `json:"regularMarketPrice"`
	RegularMarketChange        float32 `json:"regularMarketChange"`
	RegularMarketChangePercent float32 `json:"regularMarketChangePercent"`
}

type quoteResponse struct {
	Result []financeMessage `json:"result"`
	Error  string           `json:"error"`
}

type yfMessage struct {
	Response quoteResponse `json:"quoteResponse"`
}

type color string

const (
	colorBold  color = "\033[1;37m"
	colorGreen       = "\033[32m"
	colorRed         = "\033[31m"
	colorReset       = "\033[00m"
)

type preparedMarketInfo struct {
	symbol               string
	nonRegularMarketSign string
	price                float32
	diff                 float32
	percent              float32
}

func retrieveYFInfo(tickers []string, httpClient *http.Client, logger *zap.Logger) ([]byte, error) {
	symbols := strings.Join(tickers, ",")
	fields := strings.Join(requestFields, ",")

	if logger != nil {
		logger.Debug("symbols", zap.String("symbols", symbols))
	}

	fullRequestURL := fmt.Sprintf("%s&fields=%s&symbols=%s", apiEndpoint, fields, symbols)

	var res *http.Response
	var err error

	if httpClient != nil {
		res, err = httpClient.Get(fullRequestURL)
	} else {
		res, err = http.Get(fullRequestURL)
	}

	if err != nil && logger != nil {
		logger.Fatal("failed to fetch URL", zap.Error(err))
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if logger != nil {
		logger.Debug("Response body:", zap.String("body", string(body[:])))
	}

	return body, err
}

func printYahooFinanceMessages(body []byte, logger *zap.Logger) {
	var qr yfMessage

	jerr := json.Unmarshal(body, &qr)

	if jerr != nil && logger != nil {
		logger.Fatal("failed to unmarshal response body", zap.Error(jerr))
	}

	if logger != nil {
		logger.Debug("Result: ", zap.Any("result", qr))
	}

	// for _, financeMessage := range qr.Result {
	for _, financeMessage := range qr.Response.Result {
		var symbol, nonRegularMarketSign string
		var price, diff, percent float32

		if financeMessage.MarketState == "" {
			continue
		}

		symbol = financeMessage.Symbol

		if financeMessage.MarketState == "PRE" &&
			financeMessage.PreMarketChange != 0 {
			nonRegularMarketSign = "*"
			price = financeMessage.PreMarketPrice
			diff = financeMessage.PreMarketChange
			percent = financeMessage.PreMarketChangePercent
		} else if financeMessage.MarketState != "REGULAR" &&
			financeMessage.PostMarketChange != 0 {
			nonRegularMarketSign = "*"
			price = financeMessage.PostMarketPrice
			diff = financeMessage.PostMarketChange
			percent = financeMessage.PostMarketChangePercent
		} else {
			nonRegularMarketSign = ""
			price = financeMessage.RegularMarketPrice
			diff = financeMessage.RegularMarketChange
			percent = financeMessage.RegularMarketChangePercent
		}

		var color string
		if diff == 0 || os.Getenv("NO_COLOR") != "" {
			color = ""
		} else if diff < 0 {
			color = colorRed
		} else {
			color = colorGreen
		}

		if price != 0 {
			fmt.Printf("%-10s%s%8.2f%s", symbol, colorBold, price, colorReset)
			fmt.Printf("%s%10.2f%10s(%.2f%%)%s", color, diff, "", percent, colorReset)
			fmt.Printf(" %s\n", nonRegularMarketSign)
		}
	}
}

var logger *zap.Logger

func main() {
	tickers := os.Args[1:]

	body, err := retrieveYFInfo(tickers, nil, logger)

	if err != nil {
		os.Exit(1)
	}

	printYahooFinanceMessages(body, logger)
}
