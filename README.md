# Stock Info CLI

Stock Info CLI is a command-line application written in Go that allows users to fetch and display stock market information for a given company. The application utilizes the Finnhub API to retrieve real-time stock data and displays the information in a user-friendly format, complete with ASCII art and color-coded text.

Project state -> In Development

## Features

- Fetch stock information using the company name.
- Display stock prices in various currencies.
- Color-coded ASCII art for the company name.
- Horizontally aligned stock information for better readability.

## Prerequisites

- Go (version 1.16 or higher)
- Finnhub API Key (get one from [Finnhub.io](https://finnhub.io/))
- Internet connection for API requests

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/stock-info-cli.git
    cd stock-info-cli
    ```

2. Install the required Go packages:

    ```sh
    go get github.com/fatih/color
    go get github.com/common-nighthawk/go-figure
    ```

3. Replace the placeholder API key in `main.go` with your actual Finnhub API key:

    ```go
    const apiKey = "YOUR_FINNHUB_API_KEY"
    ```

## Usage

Run the application with the company name and desired currency as arguments:

```sh
go run main.go "Apple" "EUR"
```

## Supported Currencies

- USD (United States Dollar)
- EUR (Euro)
- GDP (British Pound)
 JPY (Japanese Yen)
