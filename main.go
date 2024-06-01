// api-key: "cpdlu61r01qh24flcgqgcpdlu61r01qh24flcgr0"

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

const (
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

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func fetchStockQuote(symbol, apiKey string) (*StockQuote, error) {
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

func fetchStockSymbol(companyName, apiKey string) (string, error) {
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

func printStockQuote(symbol string, quote *StockQuote) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	bold := color.New(color.FgHiWhite, color.Bold).SprintFunc()

	// Print stock name in a larger font style
	fmt.Printf("\n%s\n\n", bold(strings.ToUpper(symbol)))

	// Print table header
	fmt.Printf("%-25s %s\n", cyan("DESCRIPTION"), cyan("VALUE"))

	// Print table rows
	fmt.Printf("%-25s %s\n", "Current Price", green(fmt.Sprintf("%.2f", quote.C)))
	fmt.Printf("%-25s %s\n", "High Price of the Day", green(fmt.Sprintf("%.2f", quote.H)))
	fmt.Printf("%-25s %s\n", "Low Price of the Day", red(fmt.Sprintf("%.2f", quote.L)))
	fmt.Printf("%-25s %s\n", "Open Price of the Day", yellow(fmt.Sprintf("%.2f", quote.O)))
	fmt.Printf("%-25s %s\n", "Previous Close Price", yellow(fmt.Sprintf("%.2f", quote.Pc)))
}

func main() {
	loadEnv()

	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatalf("FINNHUB_API_KEY is not set in the .env file")
	}

	if len(os.Args) < 3 {
		fmt.Println("Please provide a company name and currency.")
		return
	}
	companyName := os.Args[1]

	// Fetch the stock symbol from the company name
	symbol, err := fetchStockSymbol(companyName, apiKey)
	if err != nil {
		fmt.Println("Error fetching stock symbol:", err)
		return
	}

	// Fetch the stock quote using the symbol
	quote, err := fetchStockQuote(symbol, apiKey)
	if err != nil {
		fmt.Println("Error fetching stock quote:", err)
		return
	}

	// Fetch exchange rates
	_, err = fetchExchangeRates()
	if err != nil {
		fmt.Println("Error fetching exchange rates:", err)
		return
	}

	// Print the stock quote
	printStockQuote(symbol, quote)
}
