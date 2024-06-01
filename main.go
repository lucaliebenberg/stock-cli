// api-key: "cpdlu61r01qh24flcgqgcpdlu61r01qh24flcgr0"

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

const (
	apiKey             = "cpdlu61r01qh24flcgqgcpdlu61r01qh24flcgr0"
	quoteApiUrl        = "https://finnhub.io/api/v1/quote?symbol=%s&token=%s"
	searchApiUrl       = "https://finnhub.io/api/v1/search?q=%s&token=%s"
	exchangeRateApiUrl = "https://api.exchangerate-api.com/v4/latest/USD"
)

type StockQuote struct {
	C  float64 `json:"c"`
	H  float64 `json:"h"`
	L  float64 `json:"l"`
	O  float64 `json:"o"`
	Pc float64 `json:"pc"`
}

type SearchResult struct {
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
}

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

func fetchStockQuote(symbol string) (*StockQuote, error) {
	url := fmt.Sprintf(quoteApiUrl, symbol, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var quote StockQuote
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return nil, err
	}
	return &quote, nil
}

func fetchStockSymbol(companyName string) (string, error) {
	url := fmt.Sprintf(searchApiUrl, companyName, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Result []SearchResult `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Result) == 0 {
		return "", fmt.Errorf("no results found for company: %s", companyName)
	}

	return result.Result[0].Symbol, nil
}

func fetchExchangeRates() (*ExchangeRates, error) {
	resp, err := http.Get(exchangeRateApiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rates ExchangeRates
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return nil, err
	}
	return &rates, nil
}

func convertCurrency(value float64, rates *ExchangeRates, currency string) float64 {
	rate, exists := rates.Rates[currency]
	if !exists {
		fmt.Printf("Currency %s not found. Using USD instead.\n", currency)
		return value
	}
	return value * rate
}

func printStockQuote(symbol string, quote *StockQuote, currency string, rates *ExchangeRates) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	// bold := color.New(color.FgHiWhite, color.Bold).SprintFunc()

	convertedC := convertCurrency(quote.C, rates, currency)
	convertedH := convertCurrency(quote.H, rates, currency)
	convertedL := convertCurrency(quote.L, rates, currency)
	convertedO := convertCurrency(quote.O, rates, currency)
	convertedPc := convertCurrency(quote.Pc, rates, currency)

	// Create ASCII art for the company name
	companyArt := figure.NewFigure(strings.ToUpper(symbol), "", true)
	companyArt.Print()

	// Print table header
	fmt.Printf("%-25s %s\n", cyan("DESCRIPTION"), cyan("VALUE"))

	// Print table rows
	fmt.Printf("%-25s %s %.2f\n", "Current Price", green(currency), convertedC)
	fmt.Printf("%-25s %s %.2f\n", "High Price of the Day", green(currency), convertedH)
	fmt.Printf("%-25s %s %.2f\n", "Low Price of the Day", red(currency), convertedL)
	fmt.Printf("%-25s %s %.2f\n", "Open Price of the Day", yellow(currency), convertedO)
	fmt.Printf("%-25s %s %.2f\n", "Previous Close Price", yellow(currency), convertedPc)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide a company name and currency.")
		return
	}
	companyName := os.Args[1]
	currency := os.Args[2]

	// Fetch the stock symbol from the company name
	symbol, err := fetchStockSymbol(companyName)
	if err != nil {
		fmt.Println("Error fetching stock symbol:", err)
		return
	}

	// Fetch the stock quote using the symbol
	quote, err := fetchStockQuote(symbol)
	if err != nil {
		fmt.Println("Error fetching stock quote:", err)
		return
	}

	// Fetch exchange rates
	rates, err := fetchExchangeRates()
	if err != nil {
		fmt.Println("Error fetching exchange rates:", err)
		return
	}

	// Print the stock quote
	printStockQuote(companyName, quote, currency, rates)
}
