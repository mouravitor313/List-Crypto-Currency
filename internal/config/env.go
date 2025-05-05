package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

var (
	CoinGeckoAPIKey string
	CurrencyLayerAPIKey string
)

func LoadAPIKey() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("fail to to load .env: %v", err)
	}

	CoinGeckoAPIKey = os.Getenv("COINGECKO_API_KEY")
	if CoinGeckoAPIKey == "" {
		return fmt.Errorf("CoinGecko API_KEY not found")
	}

	CurrencyLayerAPIKey = os.Getenv("CURRENCY_LAYER_API")
	if CurrencyLayerAPIKey == "" {
		return fmt.Errorf("CurrencyLayer API_KEY not found")
	}

	return nil
}