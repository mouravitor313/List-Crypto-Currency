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

func LoadAPIKey() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar .env:", err)
		return
	}

	CoinGeckoAPIKey = os.Getenv("COINGECKO_API_KEY")
	if CoinGeckoAPIKey == "" {
		fmt.Println("API_KEY não encontrada no .env!")
		os.Exit(1)
	}

	CurrencyLayerAPIKey = os.Getenv("CURRENCY_LAYER_API")
	if CurrencyLayerAPIKey == "" {
		fmt.Println("API_KEY não encontrada no .env!")
		os.Exit(1)
	}
}